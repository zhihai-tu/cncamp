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
[前序作业](https://github.com/zhihai-tu/cncamp/tree/main/homework/module8/ex8.1)
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
### 应用高可用
httpserver的pod有两个副本，保证了高可用
```sh
ubuntu@VM-4-4-ubuntu:~$ k get po
NAME                          READY   STATUS    RESTARTS   AGE
httpserver-8588c7fd96-8p67z   1/1     Running   0          5d22h
httpserver-8588c7fd96-hr292   1/1     Running   0          5d22h
```

### ingress
安装ingress
```sh
k create -f nginx-ingress-deployment.yaml
```
注意点：  
1. 查看文件nginx-ingress-deployment.yaml，将以下两个镜像中的@sha开头的内容去掉：
   + k8s.gcr.io/ingress-nginx/controller:v1.0.0@sha256:0851b34f69f69352bf168e6ccf30e1e20714a264ab1ecd1933e4d8c0fc3215c6
   + k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0@sha256:f3b6b39a6062328c095337b4cadcefd1612348fdd5190b1dcbcb9b9e90bd8068     
2. 由于无法下载k8s.gcr.io镜像，需要下载国内镜像源中下载，然后在修改tag，操作如下：  
查询国内镜像源，可以使用docker search命令
```sh
root@VM-4-4-ubuntu:~# docker search kube-webhook-certgen
NAME                                 DESCRIPTION                                     STARS     OFFICIAL   AUTOMATED
jettech/kube-webhook-certgen         Kubernetes webhook certificate generator and…   13                   
liangjw/kube-webhook-certgen         k8s.gcr.io/ingress-nginx/kube-webhook-certgen   8                    
dyrnq/kube-webhook-certgen           k8s.gcr.io/ingress-nginx/kube-webhook-certgen   2                    
wangshun1024/kube-webhook-certgen                                                    2                    
wonderflow/kube-webhook-certgen                                                      1                    
lianyuxue1020/kube-webhook-certgen   new pull lianyuxue1020/kube-webhook-certgen:…   1                    
hzde0128/kube-webhook-certgen                                                        0                    
newrelic/kube-webhook-certgen                                                        0                    
spdplx2021/kube-webhook-certgen                                                      0                    
pangser/kube-webhook-certgen                                                         0                    
catalystcloud/kube-webhook-certgen                                                   0                    
```
或者可以登录dockerhub页面后，进行搜索，如下图
![dockerhub](https://github.com/zhihai-tu/cncamp/raw/main/homework/module8/ex8.2/dockerhub-search.jpg)

选定国内镜像后，使用crictl pull命令进行下载
```sh
root@VM-4-4-ubuntu:~# crictl pull liangjw/kube-webhook-certgen:v1.0
Image is up to date for sha256:17e55ec30f203e6acb1e2d35bf8af5e171b3734539e1d2b560c8e80f6b1b259a

root@VM-4-4-ubuntu:~# crictl pull liangjw/ingress-nginx-controller:v1.0.0
Image is up to date for sha256:ef43679c2cae7c3812814f91faa4c76de95152daa9dc6f52836f6262946f5825
```
使用ctr命令修改tag（注意containd有namespace的概念，需要加上-n k8s.io)
```sh
root@VM-4-4-ubuntu:~# ctr -n k8s.io i tag docker.io/liangjw/kube-webhook-certgen:v1.0 k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0
k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0

root@VM-4-4-ubuntu:~# ctr -n k8s.io i tag docker.io/liangjw/ingress-nginx-controller:v1.0.0 k8s.gcr.io/ingress-nginx/controller:v1.0.0
k8s.gcr.io/ingress-nginx/controller:v1.0.0
```
查看镜像
```sh
root@VM-4-4-ubuntu:~# crictl images --digests| grep k8s.gcr.io
k8s.gcr.io/ingress-nginx/controller                               v1.0.0              0851b34f69f69       ef43679c2cae7       103MB
k8s.gcr.io/ingress-nginx/kube-webhook-certgen                     v1.0                f3b6b39a60623       17e55ec30f203       18.6MB
```
ingress安装完毕后，查看pod
```sh
ubuntu@VM-4-4-ubuntu:~$ k get po -n ingress-nginx
NAME                                       READY   STATUS      RESTARTS   AGE
ingress-nginx-admission-create-wzx6p       0/1     Completed   0          11h
ingress-nginx-admission-patch-5kszv        0/1     Completed   1          11h
ingress-nginx-controller-fd8b8b55b-7xklw   1/1     Running     0          11h
```
新建一个ingress网关
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ kubectl create -f ingress-http.yaml
ingress.networking.k8s.io/gateway created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ k get ing
NAME      CLASS   HOSTS   ADDRESS    PORTS   AGE
gateway   nginx   *       10.0.4.4   80      74s
```
查看ingress，已经按照规则绑定了backends后端服务
```sh 
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ k describe ing gateway
Name:             gateway
Labels:           <none>
Namespace:        default
Address:          10.0.4.4
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host        Path  Backends
  ----        ----  --------
  *           
              /   httpserver:80 (192.168.182.241:8080,192.168.182.242:8080)
Annotations:  <none>
Events:
  Type    Reason  Age                 From                      Message
  ----    ------  ----                ----                      -------
  Normal  Sync    48s (x2 over 103s)  nginx-ingress-controller  Scheduled for sync
```
通过ingress的service进行访问
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ k get svc  -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.103.35.254   <none>        80:30576/TCP,443:32200/TCP   112m
ingress-nginx-controller-admission   ClusterIP   10.108.78.3     <none>        443/TCP                      112m
```
通过ingress的入口进行访问，先把请求转到ingress的pod中去，然后通过ingress pod的转发规则（上述getway)再转到httpservice的pod中去。
```sh
ubuntu@VM-4-4-ubuntu:~$ curl 10.103.35.254
<html h1>Welcome to cncamp...</html>
ubuntu@VM-4-4-ubuntu:~$ curl localhost:30576
<html h1>Welcome to cncamp...</html>
```
问题：为什么curl 10.103.35.254/healthz没有返回结果？
```sh
ubuntu@VM-4-4-ubuntu:~$ curl 10.103.35.254/healthz
ubuntu@VM-4-4-ubuntu:~$ curl 10.103.35.254/healthz
ubuntu@VM-4-4-ubuntu:~$ curl 10.103.35.254/healthz
```

### 通过证书保证 httpServer 的通讯安全
签发域名为cncamp.com的证书
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
Generating a RSA private key
...........................................................+++++
.............................+++++
writing new private key to 'tls.key'
-----
```
创建secret，并进行查看(后续可提供https服务)
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ kubectl create secret tls cncamp-tls --cert=./tls.crt --key=./tls.key
secret/cncamp-tls created
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ k get secret
NAME                  TYPE                                  DATA   AGE
cncamp-tls            kubernetes.io/tls                     2      16s
default-token-zmnx9   kubernetes.io/service-account-token   3      21d
```
更新ingress，加入域名及https认证配置
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.2$ k replace -f ingress-https.yaml 
ingress.networking.k8s.io/gateway replaced
```
通过ingress的入口，先把请求转到ingress的pod中去，然后通过ingress pod的转发规则（上述getway)再转到httpservice的pod中去。
```sh
ubuntu@VM-4-4-ubuntu:~$ curl -H "Host: cncamp.com" https://10.103.35.254 -k
<html h1>Welcome to cncamp...</html>
ubuntu@VM-4-4-ubuntu:~$ curl -H "Host: cncamp.com" https://10.103.35.254/healthz -k
<html h1>system is working... httpcode: 200 </html>
```
问题：为什么不能通过nodeport的https进行访问？
```sh
ubuntu@VM-4-4-ubuntu:~$ curl -H "Host: cncamp.com" https://localhost:30576 -k
curl: (35) error:1408F10B:SSL routines:ssl3_get_record:wrong version number
```
