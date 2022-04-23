## 模块十作业（必交）

1. 为 HTTPServer 添加 0-2 秒的随机延时；
2. 为 HTTPServer 项目添加延时 Metric；
3. 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；
4. 从 Promethus 界面中查询延时指标数据；
5. （可选）创建一个 Grafana Dashboard 展现延时分配情况。

提交地址： https://jinshuju.net/f/awEgbi  
截止日期：2022 年 4 月 24 日 23:59  

## 作业解答
### STEP1.为 HTTPServer 添加 0-2 秒的随机延时
修改main.go程序，在rootHandler函数中增加如下代码
```go
//ex10-func-01:为 HTTPServer 添加 0-2 秒的随机延时
delay := randInt(0, 2000)
time.Sleep(time.Millisecond * time.Duration(delay))
glog.V(3).Infof("Respond in %d ms", delay)
```
### STEP2.为 HTTPServer 项目添加延时 Metric
新增metrics.go，定义prometheus采集指标  
修改main.go程序，增加merics相关内容
```go
func main() {
    ……
	//ex10-func-02:增加metrics的register，注册prometheus的指标采集器（直方图采集器，采集指标是execution_latency_seconds，详见metrics.go代码
	metrics.Register()
    ……
    //ex10-func-02:/metrics路径注册为prometheus的handler
	mux.Handle("/metrics", promhttp.Handler())
    ……
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	//ex10-func-02:记录整个函数的执行时间
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
    ……
}
```
### STEP3.将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；
#### 安装Prometheus
1. 参考如下安装教程，安装loki-stack：https://github.com/cncamp/101/tree/master/module10/loki-stack  
2. 如果image无法下载，手工docker pull后修改tag解决
3. 查看pod和service
```sh
cadmin@k8snode:~$ k get po
NAME                                            READY   STATUS    RESTARTS        AGE
loki-0                                          1/1     Running   1 (6m56s ago)   179m
loki-grafana-866d588467-nsvdh                   2/2     Running   4 (6m56s ago)   38m
loki-kube-state-metrics-5d666fbb55-vk9q9        1/1     Running   1 (6m56s ago)   143m
loki-prometheus-alertmanager-649cc4f455-gjp87   2/2     Running   5 (6m56s ago)   179m
loki-prometheus-node-exporter-c9n9k             1/1     Running   3 (6m56s ago)   179m
loki-prometheus-pushgateway-575b7f6bfd-7n7v7    1/1     Running   3 (6m56s ago)   179m
loki-prometheus-server-b4c6f96bd-xz7rv          2/2     Running   5 (6m56s ago)   179m
loki-promtail-26cwm                             1/1     Running   0               179m

cadmin@k8snode:~$ k get svc
NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)           AGE
envoy                           NodePort    10.97.219.161    <none>        10008:32565/TCP   36d
kubernetes                      ClusterIP   10.96.0.1        <none>        443/TCP           39d
loki                            ClusterIP   10.104.210.40    <none>        3100/TCP          179m
loki-grafana                    NodePort    10.106.221.30    <none>        80:32733/TCP      179m
loki-headless                   ClusterIP   None             <none>        3100/TCP          179m
loki-kube-state-metrics         ClusterIP   10.105.87.141    <none>        8080/TCP          179m
loki-prometheus-alertmanager    ClusterIP   10.111.236.91    <none>        80/TCP            179m
loki-prometheus-node-exporter   ClusterIP   None             <none>        9100/TCP          179m
loki-prometheus-pushgateway     ClusterIP   10.98.85.77      <none>        9091/TCP          179m
loki-prometheus-server          ClusterIP   10.103.103.103   <none>        80/TCP            179m
```
4. loki-grafana修改为NodePort，可正常打开grafana的网页

#### 部署应用
1. 制作镜像
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ex10]# docker build -t tuzhihai1986/httpserver:v3.0.1-metrics .
Sending build context to Docker daemon     64kB
Step 1/10 : FROM golang:1.17 AS build
 ---> 0659a535a734
Step 2/10 : WORKDIR /app/
 ---> Running in 235d20ed6db2
Removing intermediate container 235d20ed6db2
 ---> f4c0bedf9583
Step 3/10 : COPY main.go go.mod go.sum ./
 ---> 0455c06bac33
Step 4/10 : COPY metrics/metrics.go ./metrics/
 ---> e6890ea1bb4c
Step 5/10 : ENV GO111MODULE=on     CGO_ENABLED=0     GOOS=linux     GOARCH=amd64     GOPROXY=https://goproxy.cn,direct
 ---> Running in cb42dd8a3d1f
Removing intermediate container cb42dd8a3d1f
 ---> 7cd72b7c4f02
Step 6/10 : RUN go build -o httpserver .
 ---> Running in 619281b8560e
go: downloading github.com/golang/glog v1.0.0
go: downloading github.com/prometheus/client_golang v1.12.1
go: downloading github.com/prometheus/client_model v0.2.0
go: downloading github.com/prometheus/common v0.32.1
go: downloading github.com/beorn7/perks v1.0.1
go: downloading github.com/cespare/xxhash/v2 v2.1.2
go: downloading github.com/golang/protobuf v1.5.2
go: downloading github.com/prometheus/procfs v0.7.3
go: downloading google.golang.org/protobuf v1.26.0
go: downloading github.com/matttproud/golang_protobuf_extensions v1.0.1
go: downloading golang.org/x/sys v0.0.0-20220114195835-da31bd327af9
Removing intermediate container 619281b8560e
 ---> dd226a3b1f97
Step 7/10 : FROM busybox
 ---> 829374d342ae
Step 8/10 : COPY --from=build /app/httpserver .
 ---> 20a6566c5656
Step 9/10 : EXPOSE 80
 ---> Running in 4999830c5f5a
Removing intermediate container 4999830c5f5a
 ---> 7eb70e61061e
Step 10/10 : ENTRYPOINT ["/httpserver"]
 ---> Running in 95a91c03288c
Removing intermediate container 95a91c03288c
 ---> ab7e37dcbc51
Successfully built ab7e37dcbc51
Successfully tagged tuzhihai1986/httpserver:v3.0.1-metrics
```
2. 上传镜像至dockerhub
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ex10]# docker push tuzhihai1986/httpserver:v3.0.1-metrics
The push refers to repository [docker.io/tuzhihai1986/httpserver]
1b4aab581124: Pushed 
252fdf0c3b6a: Mounted from library/busybox 
v3.0.1-metrics: digest: sha256:72991c66708289e8344e0600baa07dd1a0e5a1df0d108a20b37ca08d4d71ffd6 size: 738
```
3. 修改httpserver-deployment.ymal文件,片段如下（详见注释处）
```
  template:
    metadata:
      ##ex10:定义这个deployment需要汇报指标，通过8080端口
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
      creationTimestamp: null
      labels:
        app: httpserver
    spec:
      containers:
      - image: tuzhihai1986/httpserver:v3.0.1-metrics
        imagePullPolicy: IfNotPresent
        name: httpserver
        ##ex10:暴露8080端口
        ports: 
          - containerPort: 8080
```
4. 创建configMap
```sh
cadmin@k8snode:~/tuzhihai/module10$ k create cm myenv --from-file=app.properties
configmap/myenv created

cadmin@k8snode:~/tuzhihai/module10$ k create cm myenv1 --from-env-file=app.properties
configmap/myenv1 created
```

5. 创建httpserver
```sh
cadmin@k8snode:~/tuzhihai/module10$ k create -f httpserver-deployment.yaml 
deployment.apps/httpserver created

cadmin@k8snode:~/tuzhihai/module10$ k get po
NAME                                            READY   STATUS    RESTARTS      AGE
httpserver-7b44cc8b47-7s8fs                     1/1     Running   0             10s
httpserver-7b44cc8b47-tzbvh                     1/1     Running   0             10s
```



