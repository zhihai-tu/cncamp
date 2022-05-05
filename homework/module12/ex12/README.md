## 模块十二作业（必交）
把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：  
1. 如何实现安全保证；
2. 七层路由规则；
3. 考虑 open tracing 的接入。

## 作业解答
### 通过Istio Ingress Gateway的形式来发布httpserver
1. 安装istio
```sh
root@VM-4-4-ubuntu:~/app$ export ISTIO_VERSION=1.12.0
root@VM-4-4-ubuntu:~/app# curl -L https://istio.io/downloadIstio | sh -
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   101  100   101    0     0    289      0 --:--:-- --:--:-- --:--:--   289
100  4926  100  4926    0     0   8796      0 --:--:-- --:--:-- --:--:--  8796

Downloading istio-1.12.0 from https://github.com/istio/istio/releases/download/1.12.0/istio-1.12.0-linux-amd64.tar.gz ...

Istio 1.12.0 Download Complete!

Istio has been successfully downloaded into the istio-1.12.0 folder on your system.

Next Steps:
See https://istio.io/latest/docs/setup/install/ to add Istio to your Kubernetes cluster.

To configure the istioctl client tool for your workstation,
add the /root/app/istio-1.12.0/bin directory to your environment path variable with:
         export PATH="$PATH:/root/app/istio-1.12.0/bin"

Begin the Istio pre-installation check by running:
         istioctl x precheck 

Need more information? Visit https://istio.io/latest/docs/setup/install/ 

root@VM-4-4-ubuntu:~/app# cd istio-1.12.0
root@VM-4-4-ubuntu:~/app/istio-1.12.0# cp bin/istioctl /usr/local/bin
```
切换到ubuntu用户，安装istio
```sh
ubuntu@VM-4-4-ubuntu:~$ istioctl install --set profile=demo -y
✔ Istio core installed                                                                                                                     
✔ Istiod installed                                                                                                                         
✔ Egress gateways installed                                                                                                                
✔ Ingress gateways installed                                                                                                               
✔ Installation complete                                                                                                                    Making this installation the default for injection and validation.

Thank you for installing Istio 1.12.  Please take a few minutes to tell us about your install/upgrade experience!  https://forms.gle/FegQbc9UvePd4Z9z7
```
验证
```sh
ubuntu@VM-4-4-ubuntu:~$ k get ns
NAME               STATUS   AGE
calico-apiserver   Active   40d
calico-system      Active   40d
default            Active   40d
ingress-nginx      Active   18d
istio-demo         Active   107m
istio-system       Active   7m12s
kube-node-lease    Active   40d
kube-public        Active   40d
kube-system        Active   40d
tigera-operator    Active   40d

ubuntu@VM-4-4-ubuntu:~$ k get po -n istio-system
NAME                                    READY   STATUS    RESTARTS   AGE
istio-egressgateway-7689b99d44-pxtbt    1/1     Running   0          5m6s
istio-ingressgateway-579b8b4bf4-ppwd4   1/1     Running   0          5m6s
istiod-7f67886588-5m28r                 1/1     Running   0          7m28s
```

2. 在新的namespace中创建httpserver的pod和service
```sh
ubuntu@VM-4-4-ubuntu:~$ k create ns istio-demo
namespace/istio-demo created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create cm myenv -n istio-demo --from-file=app.properties
configmap/myenv created
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create cm myenv1 -n istio-demo --from-env-file=app.properties
configmap/myenv1 created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get cm -n istio-demo
NAME               DATA   AGE
kube-root-ca.crt   1      12m
myenv              1      46s
myenv1             3      20s

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -f httpserver-deployment.yaml -n istio-demo
deployment.apps/httpserver created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get po -n istio-demo -owide
NAME                          READY   STATUS    RESTARTS   AGE     IP                NODE            NOMINATED NODE   READINESS GATES
httpserver-7d4bbb44f5-grv98   1/1     Running   0          8m27s   192.168.182.255   vm-4-4-ubuntu   <none>           <none>
httpserver-7d4bbb44f5-tlstj   1/1     Running   0          8m27s   192.168.182.195   vm-4-4-ubuntu   <none>           <none>
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 192.168.182.255
curl: (7) Failed to connect to 192.168.182.255 port 80: Connection refused
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 192.168.182.255:8080
<html h1>Welcome to cncamp...</html>

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -f httpserver-service.yaml -n istio-demo
service/httpserver created
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get svc -n istio-demo
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
httpserver   ClusterIP   10.105.91.245   <none>        80/TCP    14s
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 10.105.91.245
<html h1>Welcome to cncamp...</html>
```