## 1.编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化
详见[程序包](https://github.com/zhihai-tu/cncamp/tree/main/homework/module3/ex3.2)
```shell
[root@iZuf6hgwe067pstqgg2aj7Z ex3.2]# docker build -t tuzhihai1986/httpserver:v2.0 .
Sending build context to Docker daemon  7.043MB
Step 1/9 : FROM golang:1.17 AS build
 ---> 0659a535a734
Step 2/9 : WORKDIR /app/
 ---> Using cache
 ---> 9fd3e453e178
Step 3/9 : COPY main.go .
 ---> Using cache
 ---> 9bf34853cb56
Step 4/9 : ENV GO111MODULE=off     CGO_ENABLED=0     GOOS=linux     GOARCH=amd64
 ---> Running in fe30bbba7736
Removing intermediate container fe30bbba7736
 ---> 88aedceb7a5d
Step 5/9 : RUN go build -o httpserver .
 ---> Running in b21d59035be2
Removing intermediate container b21d59035be2
 ---> e32f3a6d7b75
Step 6/9 : FROM scratch
 ---> 
Step 7/9 : COPY --from=build /app/httpserver /
 ---> 4f018a5bee07
Step 8/9 : EXPOSE 80
 ---> Running in c4f0727d2ed2
Removing intermediate container c4f0727d2ed2
 ---> 8bec7a88e4bd
Step 9/9 : ENTRYPOINT ["/httpserver"]
 ---> Running in e4d2896f5ef7
Removing intermediate container e4d2896f5ef7
 ---> b63ce3a5e4b0
Successfully built b63ce3a5e4b0
Successfully tagged tuzhihai1986/httpserver:v2.0
```

## 2.将镜像推送至 docker 官方镜像仓库
```shell
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
cncamp/httpserver   v1.0                1a5f3eb5c9d5        About an hour ago   12.6 MB
docker.io/ubuntu    latest              54c9d81cbb44        4 weeks ago         72.8 MB
docker.io/alpine    latest              c059bfaa849c        3 months ago        5.59 MB
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username (tuzhihai1986): ^C
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker tag cncamp/httpserver:v1.0 tuzhihai1986/httpserver:v1.0
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker images
REPOSITORY                TAG                 IMAGE ID            CREATED             SIZE
cncamp/httpserver         v1.0                1a5f3eb5c9d5        About an hour ago   12.6 MB
tuzhihai1986/httpserver   v1.0                1a5f3eb5c9d5        About an hour ago   12.6 MB
docker.io/ubuntu          latest              54c9d81cbb44        4 weeks ago         72.8 MB
docker.io/alpine          latest              c059bfaa849c        3 months ago        5.59 MB
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker push tuzhihai1986/httpserver:v1.0
The push refers to a repository [docker.io/tuzhihai1986/httpserver]
acd1579f406c: Pushed 
8d3ac3489996: Pushed 
v1.0: digest: sha256:2ce816b3f124998b1a8c16453bd79d2f3669f0c086e14db24602318c0659ff09 size: 739
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
```

## 3.通过 docker 命令本地启动 httpserver
```shell
[root@iZuf6hgwe067pstqgg2aj7Z ex3.2]# docker run -p 8000:80 -d tuzhihai1986/httpserver:v2.0
eafcea1af3b1d1adcde9329014e66b155d7d70aa72f5bb07415108564207bec0
```

```shell
[root@iZuf6hgwe067pstqgg2aj7Z ~]# curl localhost:8000/healthz
system is working... httpcode: 200 
```

## 4.通过 nsenter 进入容器查看 IP 配置 
```shell
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker ps
CONTAINER ID   IMAGE                          COMMAND         CREATED              STATUS              PORTS                                   NAMES
eafcea1af3b1   tuzhihai1986/httpserver:v2.0   "/httpserver"   About a minute ago   Up About a minute   0.0.0.0:8000->80/tcp, :::8000->80/tcp   inspiring_aryabhata
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker inspect eafcea1af3b1 | grep -i pid
            "Pid": 18922,
            "PidMode": "",
            "PidsLimit": null,
[root@iZuf6hgwe067pstqgg2aj7Z ~]# nsenter -t 18922 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
18: eth0@if19: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
```