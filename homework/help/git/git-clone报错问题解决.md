## 解决git clone各类报错问题

### 问题一：鉴权失败
* 报错信息：
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ git clone git@github.com:cncamp/101.git
Cloning into '101'...
The authenticity of host 'github.com (20.205.243.166)' can't be established.
ECDSA key fingerprint is SHA256:p2QAMXNIC1TJYWeIOttrVc98/R1BUFWu3/LiyKgUfQM.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added 'github.com,20.205.243.166' (ECDSA) to the list of known hosts.
git@github.com: Permission denied (publickey).
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
```
或者
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ git clone git@github.com:cncamp/101.git
Cloning into '101'...
git@github.com: Permission denied (publickey).
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
and the repository exists.
```

* 解决方法  
可以将本地的密钥上传至服务器解决  
进入密钥文件所在目录进行查看，发现没有公钥和私钥
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ cd ~/.ssh
ubuntu@VM-4-4-ubuntu:~/.ssh$ ll
total 16
drwx------ 2 ubuntu ubuntu 4096 Apr  9 11:03 ./
drwx------ 9 ubuntu ubuntu 4096 Apr  9 11:02 ../
-rw------- 1 ubuntu ubuntu  790 Mar 18 19:56 authorized_keys
-rw-r--r-- 1 ubuntu ubuntu  444 Apr  9 11:03 known_hosts
```

将公钥和私钥从本地上传至上述目录
```sh
ubuntu@VM-4-4-ubuntu:~/.ssh$ rz
rz waiting to receive.
Starting zmodem transfer.  Press Ctrl+C to cancel.
Transferring id_rsa...
  100%       1 KB       1 KB/sec    00:00:01       0 Errors  
Transferring id_rsa.pub...
  100%     395 bytes  395 bytes/sec 00:00:01       0 Errors  
```
![密钥文件上传](/help/git/git-keys-upload.jpg)

查看密钥文件
```sh
ubuntu@VM-4-4-ubuntu:~/.ssh$ ll
total 24
drwx------ 2 ubuntu ubuntu 4096 Apr  9 11:07 ./
drwx------ 9 ubuntu ubuntu 4096 Apr  9 11:02 ../
-rw------- 1 ubuntu ubuntu  790 Mar 18 19:56 authorized_keys
-rw-r--r-- 1 ubuntu ubuntu 1675 Feb 13 21:39 id_rsa
-rw-r--r-- 1 ubuntu ubuntu  395 Feb 13 21:39 id_rsa.pub
-rw-r--r-- 1 ubuntu ubuntu  444 Apr  9 11:03 known_hosts
```

* 验证  
再次执行git clone,发现仍旧报错（问题分析及解决，详见问题二）
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ git clone git@github.com:cncamp/101.git
Cloning into '101'...
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@         WARNING: UNPROTECTED PRIVATE KEY FILE!          @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
Permissions 0644 for '/home/ubuntu/.ssh/id_rsa' are too open.
It is required that your private key files are NOT accessible by others.
This private key will be ignored.
Load key "/home/ubuntu/.ssh/id_rsa": bad permissions
git@github.com: Permission denied (publickey).
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
and the repository exists.
```

### 问题二：权限太开放问题
* 报错信息
```
Permissions 0644 for '/home/ubuntu/.ssh/id_rsa' are too open.
It is required that your private key files are NOT accessible by others.
```

* 解决方法
```sh
ubuntu@VM-4-4-ubuntu:~/go/src/github.com/cncamp$ chmod 600 ~/.ssh/id_rsa ~/.ssh/id_rsa.pub
```

* 验证
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
