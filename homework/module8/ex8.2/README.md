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
安装helm
```sh
curl https://baltocdn.com/helm/signing.asc | sudo apt-key add -
sudo apt-get install apt-transport-https --yes
echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm
```
安装ingress
```sh

```