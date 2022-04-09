## 解决国内git clone慢的问题

### 参考材料：  
* https://www.funyan.cn/p/5321.html
* https://cloud.tencent.com/developer/article/1835785

### 推荐方法：
1、设置代理（只对 github 进行代理，对国内的仓库不影响）
```sh
git config --global http.https://github.com.proxy 127.0.0.1:1081
```

2、查看代理
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ git config --global http.https://github.com.proxy
127.0.0.1:1081
```

3、取消代理
```sh
git config --global --unset http.https://github.com.proxy
```

### 效果
设置代理加速后，下载的速度达到了4.86M/s，非常满意^^
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ git clone git@github.com:cncamp/101.git
Cloning into '101'...
remote: Enumerating objects: 1493, done.
remote: Counting objects: 100% (327/327), done.
remote: Compressing objects: 100% (112/112), done.
remote: Total 1493 (delta 239), reused 220 (delta 215), pack-reused 1166
Receiving objects: 100% (1493/1493), 24.81 MiB | 4.86 MiB/s, done.
Resolving deltas: 100% (718/718), done.
```