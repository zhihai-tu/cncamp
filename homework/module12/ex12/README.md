## 模块十二作业（必交）
把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：  
1. 如何实现安全保证；
2. 七层路由规则；
3. 考虑 open tracing 的接入。

## 作业解答
### STEP1:安装istio
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

### STEP2:在新的namespace中创建httpserver的pod和service
```sh
ubuntu@VM-4-4-ubuntu:~$ k create ns istio-demo
namespace/istio-demo created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create cm myenv -n istio-demo --from-file=app.properties
configmap/myenv created
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create cm myenv1 -n istio-demo --from-env-file=app.properties
configmap/myenv1 created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get cm -n istio-demo
NAME                 DATA   AGE
istio-ca-root-cert   1      33s
kube-root-ca.crt     1      33s
myenv                1      17s
myenv1               3      8s

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -f httpserver-deployment.yaml -n istio-demo
deployment.apps/httpserver created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get po -n istio-demo -owide
NAME                          READY   STATUS    RESTARTS   AGE   IP                NODE            NOMINATED NODE   READINESS GATES
httpserver-7d4bbb44f5-fjxrf   1/1     Running   0          52s   192.168.182.204   vm-4-4-ubuntu   <none>           <none>
httpserver-7d4bbb44f5-wxt77   1/1     Running   0          52s   192.168.182.205   vm-4-4-ubuntu   <none>           <none>
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 192.168.182.204:8080
<html h1>Welcome to cncamp...</html>

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -f httpserver-service.yaml -n istio-demo
service/httpserver created
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k get svc -n istio-demo
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
httpserver   ClusterIP   10.108.170.227   <none>        80/TCP    6s
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ curl 10.108.170.227
<html h1>Welcome to cncamp...</html>
```

### STEP3:创建istio网关,即创建virtualservice和getway，实现七层路由
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -f istio-specs-http.yaml -n istio-demo
virtualservice.networking.istio.io/httpserver created
gateway.networking.istio.io/httpserver created
```
```sh
ubuntu@VM-4-4-ubuntu:~$ k get svc -nistio-system
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE
istio-egressgateway    ClusterIP      10.97.250.7     <none>        80/TCP,443/TCP                                                               4h1m
istio-ingressgateway   LoadBalancer   10.99.244.113   <pending>     15021:32429/TCP,80:30529/TCP,443:31031/TCP,31400:30496/TCP,15443:30553/TCP   4h1m
istiod                 ClusterIP      10.104.233.69   <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP                                        4h4m

ubuntu@VM-4-4-ubuntu:~$ curl -H "Host: httpserver.cncamp.io" http://10.99.244.113
<html h1>Welcome to cncamp...</html>
```
### STEP4:注入Sidecar
```sh
ubuntu@VM-4-4-ubuntu:~$ k label ns istio-demo istio-injection=enabled
namespace/istio-demo labeled

##删除pod，等待重建
ubuntu@VM-4-4-ubuntu:~$ k get po -n istio-demo
NAME                          READY   STATUS    RESTARTS   AGE
httpserver-7d4bbb44f5-fjxrf   1/1     Running   0          9d
httpserver-7d4bbb44f5-wxt77   1/1     Running   0          9d
ubuntu@VM-4-4-ubuntu:~$ k get rs -n istio-demo
NAME                    DESIRED   CURRENT   READY   AGE
httpserver-7d4bbb44f5   2         2         2       9d
ubuntu@VM-4-4-ubuntu:~$ k delete rs httpserver-7d4bbb44f5 -n istio-demo
replicaset.apps "httpserver-7d4bbb44f5" deleted
## 观察此时每个pod中会有两个容器，即已经注入了Sidecar
ubuntu@VM-4-4-ubuntu:~$ k get rs -n istio-demo
NAME                    DESIRED   CURRENT   READY   AGE
httpserver-7d4bbb44f5   2         2         2       56s
ubuntu@VM-4-4-ubuntu:~$ k get po -n istio-demo
NAME                          READY   STATUS    RESTARTS   AGE
httpserver-7d4bbb44f5-npvvj   2/2     Running   0          57s
httpserver-7d4bbb44f5-rjzp4   2/2     Running   0          57s
```

### STEP5:实现安全保证（签发证书，通过https访问）
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
Generating a RSA private key
......................................+++++
.............+++++
writing new private key to 'cncamp.io.key'
-----
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
secret/cncamp-credential created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module12/ex12$ k create -f istio-specs-https.yaml -n istio-demo
virtualservice.networking.istio.io/httpsserver created
gateway.networking.istio.io/httpsserver created
```
测试
```sh
ubuntu@VM-4-4-ubuntu:~$ export INGRESS_IP=10.99.244.113
ubuntu@VM-4-4-ubuntu:~$ curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP https://httpsserver.cncamp.io/healthz -v -k
* Added httpsserver.cncamp.io:443:10.99.244.113 to DNS cache
* Hostname httpsserver.cncamp.io was found in DNS cache
*   Trying 10.99.244.113:443...
* TCP_NODELAY set
* Connected to httpsserver.cncamp.io (10.99.244.113) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: O=cncamp Inc.; CN=*.cncamp.io
*  start date: May  5 15:05:50 2022 GMT
*  expire date: May  5 15:05:50 2023 GMT
*  issuer: O=cncamp Inc.; CN=*.cncamp.io
*  SSL certificate verify result: self signed certificate (18), continuing anyway.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* Using Stream ID: 1 (easy handle 0x55d637c3be30)
> GET /healthz HTTP/2
> Host: httpsserver.cncamp.io
> user-agent: curl/7.68.0
> accept: */*
> 
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* old SSL session ID is stale, removing
* Connection state changed (MAX_CONCURRENT_STREAMS == 2147483647)!
< HTTP/2 200 
< date: Thu, 05 May 2022 16:26:57 GMT
< content-length: 51
< content-type: text/html; charset=utf-8
< x-envoy-upstream-service-time: 0
< server: istio-envoy
< 
* Connection #0 to host httpsserver.cncamp.io left intact
<html h1>system is working... httpcode: 200 </html>
```

