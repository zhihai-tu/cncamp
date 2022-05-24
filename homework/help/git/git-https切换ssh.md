## 问题：git push的时候报错，如下所示：
```sh
surface@DESKTOP-GL33P9Q MINGW64 /d/cib/Oracle培训/OCP试题/ocp (master)
$ git push
fatal: unable to access 'https://github.com/zhihai-tu/ocp.git/': Failed to connect to 127.0.0.1 port 1081: Connection refused
```
## 解决方法
1. 查看关联的所有的远程仓储名称及地址
```sh
surface@DESKTOP-GL33P9Q MINGW64 /d/cib/Oracle培训/OCP试题/ocp (master)
$ git remote -v
origin  https://github.com/zhihai-tu/ocp.git (fetch)
origin  https://github.com/zhihai-tu/ocp.git (push)
```
2. 将https修改为ssh
```sh
surface@DESKTOP-GL33P9Q MINGW64 /d/cib/Oracle培训/OCP试题/ocp (master)
$ git remote set-url origin git@github.com:zhihai-tu/ocp.git

surface@DESKTOP-GL33P9Q MINGW64 /d/cib/Oracle培训/OCP试题/ocp (master)
$ git remote -v
origin  git@github.com:zhihai-tu/ocp.git (fetch)
origin  git@github.com:zhihai-tu/ocp.git (push)
```
3. 再次git push成功
```sh
surface@DESKTOP-GL33P9Q MINGW64 /d/cib/Oracle培训/OCP试题/ocp (master)
$ git push
Enumerating objects: 11, done.
Counting objects: 100% (11/11), done.
Delta compression using up to 4 threads.
Compressing objects: 100% (8/8), done.
Writing objects: 100% (8/8), 276.79 KiB | 805.00 KiB/s, done.
Total 8 (delta 2), reused 0 (delta 0)
remote: Resolving deltas: 100% (2/2), completed with 2 local objects.
To github.com:zhihai-tu/ocp.git
   91f576b..65a8591  master -> master
```