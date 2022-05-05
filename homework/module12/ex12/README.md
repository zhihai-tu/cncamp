## 模块十二作业（必交）
把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：  
1. 如何实现安全保证；
2. 七层路由规则；
3. 考虑 open tracing 的接入。

## 作业解答
### 通过Istio Ingress Gateway的形式来发布httpserver
1. 在新的namespace中创建httpserver
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

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get po -n istio-demo
NAME                          READY   STATUS    RESTARTS   AGE
httpserver-7d4bbb44f5-grv98   1/1     Running   0          8m19s
httpserver-7d4bbb44f5-tlstj   1/1     Running   0          8m19s
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get po -n istio-demo -owide
NAME                          READY   STATUS    RESTARTS   AGE     IP                NODE            NOMINATED NODE   READINESS GATES
httpserver-7d4bbb44f5-grv98   1/1     Running   0          8m27s   192.168.182.255   vm-4-4-ubuntu   <none>           <none>
httpserver-7d4bbb44f5-tlstj   1/1     Running   0          8m27s   192.168.182.195   vm-4-4-ubuntu   <none>           <none>
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 192.168.182.255
curl: (7) Failed to connect to 192.168.182.255 port 80: Connection refused
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 192.168.182.255:8080
<html h1>Welcome to cncamp...</html>
```