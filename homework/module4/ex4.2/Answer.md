## Step1 启动一个 Envoy Deployment
## step2 要求 Envoy 的启动配置从外部的配置文件 Mount 进 Pod。

获取envoy-deploy.yaml及envoy.yaml文件，获取地址：https://github.com/cncamp/101/tree/master/module4

创建configMap
```shell
cadmin@k8snode:~/tuzhihai/module4$ k create configmap envoy-config --from-file=envoy.yaml
configmap/envoy-config created
cadmin@k8snode:~/tuzhihai/module4$ k get cm
NAME               DATA   AGE
envoy-config       1      13s
kube-root-ca.crt   1      2d6h
```

查看configMap内容
```shell
cadmin@k8snode:~/tuzhihai/module4$ k describe cm envoy-config
Name:         envoy-config
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
envoy.yaml:
----
admin:
  address:
    socket_address: { address: 127.0.0.1, port_value: 9901 }

static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 10000 }
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                stat_prefix: ingress_http
                codec_type: AUTO
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: local_service
                      domains: ["*"]
                      routes:
                        - match: { prefix: "/" }
                          route: { cluster: some_service }
                http_filters:
                  - name: envoy.filters.http.router
  clusters:
    - name: some_service
      connect_timeout: 0.25s
      type: LOGICAL_DNS
      lb_policy: ROUND_ROBIN
      load_assignment:
        cluster_name: some_service
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: nginx
                      port_value: 80

BinaryData
====

Events:  <none>
cadmin@k8snode:~/tuzhihai/module4$
```

创建pod
```shell
cadmin@k8snode:~/tuzhihai/module4$ k create -f envoy-deploy.yaml
deployment.apps/envoy created
cadmin@k8snode:~/tuzhihai/module4$ k get po
NAME                     READY   STATUS              RESTARTS   AGE
envoy-6958c489d9-b4m8h   0/1     ContainerCreating   0          7s
```

查看pod的events
```shell
cadmin@k8snode:~/tuzhihai/module4$ k describe po envoy-6958c489d9-b4m8h
Name:         envoy-6958c489d9-b4m8h
Namespace:    default
Priority:     0
Node:         k8snode/192.168.231.128
Start Time:   Thu, 17 Mar 2022 06:10:51 +0000
Labels:       pod-template-hash=6958c489d9
              run=envoy
Annotations:  cni.projectcalico.org/containerID: fd74292ffd719e8b449a97ab8dda31bc8dbcb7ba6bc0ef7ccc6e45d8e52bda22
              cni.projectcalico.org/podIP: 192.168.145.229/32
              cni.projectcalico.org/podIPs: 192.168.145.229/32
Status:       Running
IP:           192.168.145.229
IPs:
  IP:           192.168.145.229
Controlled By:  ReplicaSet/envoy-6958c489d9
Containers:
  envoy:
    Container ID:   docker://6711d65765cf34e56dfc383eeb9c9dc3c664017e3332739dbfc174661d05df2b
    Image:          envoyproxy/envoy-dev
    Image ID:       docker-pullable://envoyproxy/envoy-dev@sha256:1edfdeb9164c77451030296a2f1503f89ec3d1c5a3f1dba2927f5b868b37dcdc
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Thu, 17 Mar 2022 06:11:28 +0000
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /etc/envoy from envoy-config (ro)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-6sbd2 (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             True 
  ContainersReady   True 
  PodScheduled      True 
Volumes:
  envoy-config:
    Type:      ConfigMap (a volume populated by a ConfigMap)
    Name:      envoy-config
    Optional:  false
  kube-api-access-6sbd2:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  47s   default-scheduler  Successfully assigned default/envoy-6958c489d9-b4m8h to k8snode
  Normal  Pulling    46s   kubelet            Pulling image "envoyproxy/envoy-dev"
  Normal  Pulled     12s   kubelet            Successfully pulled image "envoyproxy/envoy-dev" in 34.860013619s
  Normal  Created    12s   kubelet            Created container envoy
  Normal  Started    11s   kubelet            Started container envoy
```

再次查看pod状态，此时已经变成running
```shell
cadmin@k8snode:~/tuzhihai/module4$ k get po
NAME                     READY   STATUS    RESTARTS   AGE
envoy-6958c489d9-b4m8h   1/1     Running   0          54s
```

## step3 进入 Pod 查看 Envoy 进程和配置。

exec进入pod
```
cadmin@k8snode:~/tuzhihai/module4$ k exec envoy-6958c489d9-b4m8h -it bash
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
root@envoy-6958c489d9-b4m8h:/# 
```

查看进程
```shell
root@envoy-6958c489d9-b4m8h:/etc/envoy# ps -ef | grep envoy
envoy          1       0  0 06:11 ?        00:00:04 envoy -c /etc/envoy/envoy.yaml
root          35      20  0 06:35 pts/0    00:00:00 grep --color=auto envoy
```

查看配置，跟configMap内容相同
```shell
root@envoy-6958c489d9-b4m8h:/# cd /etc/envoy
root@envoy-6958c489d9-b4m8h:/etc/envoy# ls -lrt
total 0
lrwxrwxrwx 1 root root 17 Mar 17 06:10 envoy.yaml -> ..data/envoy.yaml
root@envoy-6958c489d9-b4m8h:/etc/envoy# cat envoy.yaml 
admin:
  address:
    socket_address: { address: 127.0.0.1, port_value: 9901 }

static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 10000 }
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
        ……………………………………………………………………
```
## step4 更改配置的监听端口并测试访问入口的变化。
修改configMap,讲port_value：10000修改为10008
```shell
cadmin@k8snode:~/tuzhihai/module4$ k edit cm envoy-config
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: v1
data:
  envoy.yaml: "admin:\r\n  address:\r\n    socket_address: { address: 127.0.0.1, port_value:
    9901 }\r\n\r\nstatic_resources:\r\n  listeners:\r\n    - name: listener_0\r\n
    \     address:\r\n        socket_address: { address: 0.0.0.0, port_value: 10000
    }\r\n      filter_chains:\r\n        - filters:\r\n            - name: envoy.filters.network.http_connection_manager\r\n
```
此时查看pod内的配置，仍旧为10000
```shell
root@envoy-6958c489d9-b4m8h:/etc/envoy# cat envoy.yaml    
admin:
  address:
    socket_address: { address: 127.0.0.1, port_value: 9901 }

static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 10000 }
      filter_chains:
```

大约30s到1min，再次查看，系统会将pod内的配置同步为10008
```shell
root@envoy-6958c489d9-b4m8h:/etc/envoy# cat envoy.yaml 
admin:
  address:
    socket_address: { address: 127.0.0.1, port_value: 9901 }

static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 10008 }
      filter_chains:
```
## step5 通过非级联删除的方法逐个删除对象。
通常情况下，删除某个对象为级联删除，例如删除deployment(deploy)的同时，会自动删除replicationset(rs)和pod(po)。如果需要非级联删除，可使用如下命令：
```shell
kubectl delete deployment nginx-deployment --cascade=orphan
```
[参考文档](https://kubernetes.io/zh/docs/tasks/administer-cluster/use-cascading-deletion/)

实操：非级联删除deploy后，发现rs和po仍然在运行
```shell
cadmin@k8snode:~/tuzhihai/module4$ k get deploy
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
envoy   1/1     1            1           49m
cadmin@k8snode:~/tuzhihai/module4$ k delete deploy envoy --cascade=orphan
deployment.apps "envoy" deleted
cadmin@k8snode:~/tuzhihai/module4$ k get rs
NAME               DESIRED   CURRENT   READY   AGE
envoy-6958c489d9   1         1         1       49m
cadmin@k8snode:~/tuzhihai/module4$ k get po
NAME                     READY   STATUS    RESTARTS   AGE
envoy-6958c489d9-b4m8h   1/1     Running   0          49m
cadmin@k8snode:~/tuzhihai/module4$ 
```