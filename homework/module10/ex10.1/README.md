## 课后练习 10.1
将 Nginx 容器镜像上传至 Harbor Demo server 并运行在测试环境中。

## 练习解答
1. 登录harbor demo server的web页面，地址如下：  
https://demo.goharbor.io/harbor/projects  
注册用户，登录，然后创建project

2. docker login
```sh
root@VM-4-4-ubuntu:~# docker login demo.goharbor.io
Username: tuzhihai
Password: 
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded
```
3. 上传镜像到harbor demo仓库
命令：docker push demo.goharbor.io/your-project/test-image
```sh
root@VM-4-4-ubuntu:~# docker push demo.goharbor.io/tuzhihai/httpserver:v3.1
The push refers to repository [demo.goharbor.io/tuzhihai/httpserver]
695c2c1f8614: Pushed 
797ac4999b67: Pushed 
v3.1: digest: sha256:fb28eebf8f9c7b58aa0cf006a9c9d3cc9edee561fd4b0d8a115aa803d485d989 size: 738
```
4. 在web页面中查看镜像已成功上传