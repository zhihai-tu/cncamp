## Jenkins 练习
创建 Jenkins Master  
kubectl apply –f Jenkins.yaml  
kubectl apply –f sa.yaml  
  
等待 Jenkins-0 pod 运行，查看日志查找 root 密码  
kubectl logs –f Jenkins-0  
  
查看 Jenkins Service 的 NodePort，登录 Jenkins console  
 http://192.168.34.2:<nodePort>  
  
输入 root 密码并创建管理员用户，登录  


## 练习解答
