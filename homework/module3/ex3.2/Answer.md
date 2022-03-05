## 1.编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化
详见[程序包](https://github.com/zhihai-tu/cncamp/tree/main/homework/module3/ex3.2)

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
[root@iZuf6hgwe067pstqgg2aj7Z ex3.2]# docker run -p 8000:80 cncamp/httpserver:v1.0
```

## 4.通过 nsenter 进入容器查看 IP 配置 
```shell
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker ps
CONTAINER ID        IMAGE                    COMMAND                  CREATED             STATUS              PORTS                  NAMES
5c7788c8f23e        cncamp/httpserver:v1.0   "/bin/sh -c /https..."   8 minutes ago       Up 8 minutes        0.0.0.0:8000->80/tcp   frosty_edison
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker ps 
CONTAINER ID        IMAGE                    COMMAND                  CREATED             STATUS              PORTS                  NAMES
5c7788c8f23e        cncamp/httpserver:v1.0   "/bin/sh -c /https..."   9 minutes ago       Up 9 minutes        0.0.0.0:8000->80/tcp   frosty_edison
[root@iZuf6hgwe067pstqgg2aj7Z ~]# docker inspect 5c7788c8f23e | grep -i pid
            "Pid": 16548,
            "PidMode": "",
            "PidsLimit": 0,
[root@iZuf6hgwe067pstqgg2aj7Z ~]# nsenter -t 16548 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:acff:fe11:2/64 scope link 
       valid_lft forever preferred_lft forever
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
```