## 第一部分
现在你对 Kubernetes 的控制面板的工作机制是否有了深入的了解呢？  
是否对如何构建一个优雅的云上应用有了深刻的认识，那么接下来用最近学过的知识把你之前编写的 http 以优雅的方式部署起来吧，你可能需要审视之前代码是否能满足优雅上云的需求。  
作业要求：编写 Kubernetes 部署脚本将 httpserver 部署到 Kubernetes 集群，以下是你可以思考的维度。  

* 优雅启动
* 优雅终止
* 资源需求和 QoS 保证
* 探活
* 日常运维需求，日志等级
* 配置和代码分离

提交地址： https://jinshuju.net/f/rJC4DG  
截止日期：2022 年 4 月 10 日 23:59

## 作业解答

### 优雅启动、探活
设置了livenessProbe探针和readinessProbe探针，实现优雅启动，yaml文件片段如下：
```
livenessProbe:
  failureThreshold: 3
  httpGet:
    path: /healthz
    port: 8080
    scheme: HTTP
  initialDelaySeconds: 5
  periodSeconds: 10
  successThreshold: 1
  timeoutSeconds: 1
readinessProbe:
  failureThreshold: 3
  httpGet:
	path: /healthz
	port: 8080
	scheme: HTTP
  initialDelaySeconds: 5
  periodSeconds: 10
  successThreshold: 1
  timeoutSeconds: 1
```
### 优雅终止
利用 preStop先sleep 5秒钟，防止pod删除瞬间部分请求失败，yaml文件片段如下：
```
lifecycle:
  preStop:
    exec:
      command:
      - sleep
      - 5s
```
### 资源需求和Qos保证
设置资源需求（Requests）和限制（ Limits），Qos级别设置为Burstable，yaml文件片段如下：
```
resources:
  limits:
    cpu: 200m
    memory: 100Mi
  requests:
    cpu: 20m
    memory: 20Mi
```

### 配置代码分离（重点示例）
* 实现功能：将httpserver的端口定义在配置文件中
* 实现思路：httpport变量定义在app.properties中，根据此配置文件生成configmap，然后在deployment对象的yaml文件中，利用env读取configmap中的httpport变量的值，作为容器内部的系统变量。容器内的go程序会读取此系统变量作为访问端口。
* 实现步骤：  
修改go语言程序，读取环境变量HTTP_PORT
```go
	//根据环境变量获取访问端口，默认为80
	httpport := os.Getenv("HTTP_PORT")
	fmt.Printf("GET ENV: HTTP_PORT=%s", httpport)
	if httpport == "" {
		httpport = "80"
	}
	err := http.ListenAndServe(":"+httpport, mux)
	if err != nil {
		log.Fatal(err)
	}
```

重新构建httpserver镜像并上传至仓库
```sh
root@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8# docker build -t tuzhihai1986/httpserver:v3.0 .
Sending build context to Docker daemon  10.24kB
Step 1/9 : FROM golang:1.17 AS build
1.17: Pulling from library/golang
dbba69284b27: Pull complete 
9baf437a1bad: Pull complete 
6ade5c59e324: Pull complete 
b19a994f6d4c: Pull complete 
e3c59fd148be: Pull complete 
c8cbe1f86def: Pull complete 
74b19341d753: Pull complete 
Digest: sha256:f675106e44f205a7284e15cd75c41b241329f3c03ac30b0ba07b14a6ea7c99d9
Status: Downloaded newer image for golang:1.17
 ---> 5bd8c5733e7c
Step 2/9 : WORKDIR /app/
 ---> Running in 072d05a87a95
Removing intermediate container 072d05a87a95
 ---> 4231f23724d0
Step 3/9 : COPY main.go .
 ---> 0b28b087d943
Step 4/9 : ENV GO111MODULE=off     CGO_ENABLED=0     GOOS=linux     GOARCH=amd64
 ---> Running in 18f8bd95845d
Removing intermediate container 18f8bd95845d
 ---> 246f82b32be2
Step 5/9 : RUN go build -o httpserver .
 ---> Running in ef2b001a4910
Removing intermediate container ef2b001a4910
 ---> 4f0f5b974618
Step 6/9 : FROM scratch
 ---> 
Step 7/9 : COPY --from=build /app/httpserver /
 ---> a8719818ed36
Step 8/9 : EXPOSE 80
 ---> Running in c5a7132531d5
Removing intermediate container c5a7132531d5
 ---> 8302abcbd2e0
Step 9/9 : ENTRYPOINT ["/httpserver"]
 ---> Running in acb3ee1d8368
Removing intermediate container acb3ee1d8368
 ---> fef6ad0cdc77
Successfully built fef6ad0cdc77
Successfully tagged tuzhihai1986/httpserver:v3.0
```
```sh
root@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8# docker images
REPOSITORY                                                        TAG       IMAGE ID       CREATED          SIZE
tuzhihai1986/httpserver                                           v3.0      fef6ad0cdc77   16 seconds ago   7.04MB
<none>                                                            <none>    4f0f5b974618   17 seconds ago   963MB
golang                                                            1.17      5bd8c5733e7c   10 days ago      941MB
```
```sh
root@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8# docker push tuzhihai1986/httpserver:v3.0
The push refers to repository [docker.io/tuzhihai1986/httpserver]
695c2c1f8614: Pushed 
v3.0: digest: sha256:48d0a849f6ec8670db88d103bf405fdf1870766a7a818fc15dfd6eefc67cd477 size: 528
```

创建configmap，来源于app.properties文件，其中变量httpport=8080
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8$ k create cm myenv1 --from-env-file=app.properties
configmap/myenv1 created
```
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8$ k get cm myenv1 -oyaml
apiVersion: v1
data:
  httpport: "8080"
kind: ConfigMap
metadata:
  creationTimestamp: "2022-04-10T08:40:38Z"
  name: myenv1
  namespace: default
  resourceVersion: "2825765"
  uid: 03ce6c93-6c38-4294-943a-7ee614525120
```

创建deployment，观察两个副本创建成功
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8$ k create -f httpserver-deployment.yaml 
deployment.apps/httpserver created

ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8$ k get po -w
NAME                          READY   STATUS              RESTARTS   AGE
httpserver-86479f944c-bt76m   0/1     ContainerCreating   0          7s
httpserver-86479f944c-kbwkh   0/1     ContainerCreating   0          7s
httpserver-86479f944c-bt76m   1/1     Running             0          10s
httpserver-86479f944c-kbwkh   1/1     Running             0          12s
```

查看pod的ip，并测试访问8080（来自于配置文件）端口的健康检查是否生效
```sh
ubuntu@VM-4-4-ubuntu:~$ k get pod -owide
NAME                          READY   STATUS    RESTARTS   AGE   IP                NODE            NOMINATED NODE   READINESS GATES
httpserver-86479f944c-bt76m   1/1     Running   0          85s   192.168.182.208   vm-4-4-ubuntu   <none>           <none>
httpserver-86479f944c-kbwkh   1/1     Running   0          85s   192.168.182.207   vm-4-4-ubuntu   <none>           <none>

ubuntu@VM-4-4-ubuntu:~$ curl 192.168.182.207/healthz
curl: (7) Failed to connect to 192.168.182.207 port 80: Connection refused
ubuntu@VM-4-4-ubuntu:~$ curl 192.168.182.207:8080/healthz
system is working... httpcode: 200 
```
### 日常运维需求、日志等级（重点示例）
* 实现功能：app.properties文件中定了glog日志的保存路径和日志等级
* 实现思路：将app.properties（整个文件）创建为configmap，ymal文件中将其mount到容器内，容器内的go程序会去读取配置文件，根据配置文件中的保存路径和日志等级来记录日志
* 实现步骤  
根据app.properties（loglevel=5,logpath=gologs）来创建configmap
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.1$ k create cm myenv --from-file=app.properties
configmap/myenv created
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/zhihai-tu/cncamp/homework/module8/ex8.1$ k get cm myenv -oyaml
apiVersion: v1
data:
  app.properties: |-
    loglevel=5
    logpath=gologs
kind: ConfigMap
metadata:
  creationTimestamp: "2022-04-10T13:35:17Z"
  name: myenv
  namespace: default
  resourceVersion: "2864505"
  uid: de198c7f-d3d0-411f-a419-7d36de04441c
```
将配置文件挂载到容器内的根目录下，yaml文件片段如下：
```
  volumeMounts:
  - name: http-config
    mountPath: "/"
    readOnly: true
volumes:
- name: http-config
  configMap:
    name: myenv
```