### Step1:按照课上讲解的方法在本地构建一个单节点的基于 HTTPS 的 etcd 集群
下载并安装etcd
```sh
ETCD_VER=v3.5.2
DOWNLOAD_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_DIR=/root/software
ETCD_HOME=/root/etcd
rm -f ${DOWNLOAD_DIR}/etcd-${ETCD_VER}-linux-amd64.tar.gz
curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o ${DOWNLOAD_DIR}/etcd-${ETCD_VER}-linux-amd64.tar.gz
mkdir -p ${ETCD_HOME}
# --strip-components=NUMBER 解压后去除前NUMBER个层级的目录
tar xzvf ${DOWNLOAD_DIR}/etcd-${ETCD_VER}-linux-amd64.tar.gz -C ${ETCD_HOME} --strip-components=1
```

添加环境变量
```sh
[root@iZuf6hgwe067pstqgg2aj7Z etcd]# cd ~
[root@iZuf6hgwe067pstqgg2aj7Z ~]# cat .bashrc
# .bashrc
……………………………………………………………………………………
# ETCD
export PATH=$PATH:/root/etcd
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
```

验证etcd版本
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcd --version
etcd Version: 3.5.2
Git SHA: 99018a77b
Go Version: go1.16.3
Go OS/Arch: linux/amd64
```

启动etcd（使用默认端口）
```
[root@iZuf6hgwe067pstqgg2aj7Z software]# etcd
{"level":"info","ts":"2022-03-27T16:40:21.145+0800","caller":"etcdmain/etcd.go:72","msg":"Running: ","args":["etcd"]}
{"level":"warn","ts":"2022-03-27T16:40:21.145+0800","caller":"etcdmain/etcd.go:104","msg":"'data-dir' was empty; using default","data-dir":"default.etcd"}
{"level":"info","ts":"2022-03-27T16:40:21.145+0800","caller":"embed/etcd.go:131","msg":"configuring peer listeners","listen-peer-urls":["http://localhost:2380"]}
{"level":"info","ts":"2022-03-27T16:40:21.146+0800","caller":"embed/etcd.go:139","msg":"configuring client listeners","listen-client-urls":["http://localhost:2379"]}
……………………
```

启动etcd(非默认端口，自行指定，此处在默认端口前加1)
```sh
etcd --listen-client-urls 'http://localhost:12379' \
 --advertise-client-urls 'http://localhost:12379' \
 --listen-peer-urls 'http://localhost:12380' \
 --initial-advertise-peer-urls 'http://localhost:12380' \
 --initial-cluster 'default=http://localhost:12380'
```

查看etcd的成员列表
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl member list --write-out=table --endpoints=localhost:2379
+------------------+---------+---------+-----------------------+-----------------------+------------+
|        ID        | STATUS  |  NAME   |      PEER ADDRS       |     CLIENT ADDRS      | IS LEARNER |
+------------------+---------+---------+-----------------------+-----------------------+------------+
| 8e9e05c52164694d | started | default | http://localhost:2380 | http://localhost:2379 |      false |
+------------------+---------+---------+-----------------------+-----------------------+------------+
```

### step2:写一条数据 
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 put name tuzhihai
OK
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 put age 37
OK
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 put nation China
OK
```

### step3:查看数据细节 
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get name
name
tuzhihai
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get age
age
37
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get age --keys-only
age
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get age --print-value-only
37
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get --prefix n 
name
tuzhihai
nation
China
# revision是3
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get age -wjson
{"header":{"cluster_id":14841639068965178418,"member_id":10276657743932975437,"revision":3,"raft_term":2},"kvs":[{"key":"YWdl","create_revision":3,"mod_revision":3,"version":1,"value":"Mzc="}],"count":1}
```

使用watch命令监控变量
```sh
## client-1
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age

## client-2
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 put age 40
OK
# revison是4
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get age -wjson
{"header":{"cluster_id":14841639068965178418,"member_id":10276657743932975437,"revision":4,"raft_term":2},"kvs":[{"key":"YWdl","create_revision":3,"mod_revision":4,"version":2,"value":"NDA="}],"count":1}

## client-1
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age
PUT
age
40
```

从某个revison开始监控变量的变化情况
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age --rev 1
PUT
age
37
PUT
age
40
^C
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age --rev 2
PUT
age
37
PUT
age
40
^C
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age --rev 3
PUT
age
37
PUT
age
40
^C
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age --rev 4
PUT
age
40
^C
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 watch age --rev 5
```

### step4:删除数据
```sh
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 del nation
1
[root@iZuf6hgwe067pstqgg2aj7Z ~]# etcdctl --endpoints=localhost:2379 get nation
[root@iZuf6hgwe067pstqgg2aj7Z ~]# 
```