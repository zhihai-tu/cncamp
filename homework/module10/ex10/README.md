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

制作镜像
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ex10]# docker build -o tuzhihai1986/httpserver:v3.0.1-metrics .
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
```