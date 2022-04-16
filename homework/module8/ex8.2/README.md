## 2. 第二部分
除了将 httpServer 应用优雅的运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。  
来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。  
在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：  
  
Service  
Ingress  
可以考虑的细节  
  
如何确保整个应用的高可用。  
如何通过证书保证 httpServer 的通讯安全。  
[strong_begin] 提交地址： https://jinshuju.net/f/XQcv68  
截止日期：2022 年 4 月 17 日 23:59  

## 作业解答
### Service
新建httpserver-service.yaml
```sh
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ vi httpserver-service.yaml 
apiVersion: v1
kind: Service
metadata:
  name: httpserver
spec:
  selector:
    app: httpserver
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```
创建Service对象，默认为ClusterIP类型，支持集群内部访问
```sh
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ k create -f httpserver-service.yaml 
service/httpserver created
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ k get svc
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
httpserver   ClusterIP   10.101.217.132   <none>        80/TCP    4s
kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP   18d
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ k get ep
NAME         ENDPOINTS                                   AGE
httpserver   192.168.182.241:8080,192.168.182.242:8080   10s
kubernetes   10.0.4.4:6443                               18d
```
测试
```sh
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ curl 10.101.217.132
<html h1>Welcome to cncamp...</html>ubuntu@VM-4-4-ubuntu:~/tuzhihai$ curl 10.101.217.132/healthz
<html h1>system is working... httpcode: 200 </html>ubuntu@VM-4-4-ubuntu:~/tuzhihai$ 
```
将server的类型设置为NodePort类型，供外部访问
```sh
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ k edit svc httpserver

  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: httpserver
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
```
测试
```
ubuntu@VM-4-4-ubuntu:~/tuzhihai$ k get svc
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
httpserver   NodePort    10.101.217.132   <none>        80:31322/TCP   8m2s
kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP        18d

ubuntu@VM-4-4-ubuntu:~$ curl 10.0.4.4:31322
<html h1>Welcome to cncamp...</html>ubuntu@VM-4-4-ubuntu:~$ curl 10.0.4.4:31322/healthz
<html h1>system is working... httpcode: 200 </html>ubuntu@VM-4-4-ubuntu:~$ 
```

### ingress
安装ingress
```sh
k create -f nginx-ingress-deployment.yaml
```
注意点：  
一、查看文件nginx-ingress-deployment.yaml，将以下两个镜像中的@sha开头的内容去掉：
1. k8s.gcr.io/ingress-nginx/controller:v1.0.0@sha256:0851b34f69f69352bf168e6ccf30e1e20714a264ab1ecd1933e4d8c0fc3215c6
2. k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0@sha256:f3b6b39a6062328c095337b4cadcefd1612348fdd5190b1dcbcb9b9e90bd8068
二、由于无法下载k8s.gcr.io镜像，需要下载国内镜像源中下载，然后在修改tag，操作如下：
```sh
root@VM-4-4-ubuntu:~# crictl pull liangjw/kube-webhook-certgen:v1.0
Image is up to date for sha256:17e55ec30f203e6acb1e2d35bf8af5e171b3734539e1d2b560c8e80f6b1b259a
root@VM-4-4-ubuntu:~# ctr -n k8s.io i tag docker.io/liangjw/kube-webhook-certgen:v1.0 k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0
k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0
root@VM-4-4-ubuntu:~# crictl pull liangjw/ingress-nginx-controller:v1.0.0
Image is up to date for sha256:ef43679c2cae7c3812814f91faa4c76de95152daa9dc6f52836f6262946f5825
root@VM-4-4-ubuntu:~# ctr -n k8s.io i tag docker.io/liangjw/ingress-nginx-controller:v1.0.0 k8s.gcr.io/ingress-nginx/controller:v1.0.0
k8s.gcr.io/ingress-nginx/controller:v1.0.0
root@VM-4-4-ubuntu:~# crictl images --digests| grep k8s.gcr.io
k8s.gcr.io/ingress-nginx/controller                               v1.0.0              0851b34f69f69       ef43679c2cae7       103MB
k8s.gcr.io/ingress-nginx/kube-webhook-certgen                     v1.0                f3b6b39a60623       17e55ec30f203       18.6MB
```

