## 集群安装步骤参考以下链接
https://github.com/cncamp/101/tree/master/k8s-install

## 安装过程中的问题
安装完成后，发现node的状态是NotReady

通过以下命令查看日志
```shell
journalctl -f -u kubelet
```
找到报错信息“network plugin is not ready: cni config uninitialized”，百度后发现是网络插件calico没有安装成功。[参考文档](https://cloud.tencent.com/developer/article/1697300)

发现namespance为tigera-operator的pod失败了，describe看了下，原来是docker pull超时出错了
```shell
ubuntu@VM-4-4-ubuntu:~$ k get ns
NAME               STATUS   AGE
default            Active   108m
kube-node-lease    Active   108m
kube-public        Active   108m
kube-system        Active   108m
tigera-operator    Active   106m
ubuntu@VM-4-4-ubuntu:~$ k -n tigera-operator get pod
NAME                              READY   STATUS    RESTARTS      AGE
tigera-operator-b876f5799-4l625   0/1     Failed   1 (50m ago)   106m
ubuntu@VM-4-4-ubuntu:~$ 
ubuntu@VM-4-4-ubuntu:~$ 
ubuntu@VM-4-4-ubuntu:~$ k -n tigera-operator describe pod tigera-operator-b876f5799-4l625
…………………………………………………………
Events:
  Type     Reason          Age                 From               Message
  ----     ------          ----                ----               -------
  Normal   Scheduled       106m                default-scheduler  Successfully assigned tigera-operator/tigera-operator-b876f5799-4l625 to vm-4-4-ubuntu
  Warning  Failed          96m                 kubelet            Failed to pull image "quay.io/tigera/operator:v1.25.3": rpc error: code = Unknown desc = context canceled
  Warning  Failed          96m                 kubelet            Error: ErrImagePull
  Normal   BackOff         96m                 kubelet            Back-off pulling image "quay.io/tigera/operator:v1.25.3"
  Warning  Failed          96m                 kubelet            Error: ImagePullBackOff
  Normal   Pulling         96m (x2 over 106m)  kubelet            Pulling image "quay.io/tigera/operator:v1.25.3"
  Warning  Failed          69m                 kubelet            Error: ErrImagePull
  Warning  Failed          69m                 kubelet            Failed to pull image "quay.io/tigera/operator:v1.25.3": rpc error: code = Unknown desc = context canceled
```

由于安装在腾讯云服务器上，因此找到了腾讯云的Docker加速配置
```
{
  "registry-mirrors": ["https://mirror.ccs.tencentyun.com"]
}
```
具体可以查看如下文档：https://cloud.tencent.com/developer/article/1151242?from=15425



