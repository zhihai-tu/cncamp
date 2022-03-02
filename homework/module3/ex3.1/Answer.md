## 将进程添加到cgroup.procs中去

[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# ps -ef|grep malloc |grep -v grep|awk '{print $2}'  
32413  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# cat cgroup.procs  
1  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# echo 32413 > cgroup.procs  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# cat cgroup.procs   
1  
32413  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]#  

---

## 设置cgroup内存限制
[root@iZuf6hgwe067pstqgg2aj7Z ~]# cd /sys/fs/cgroup/memory  
[root@iZuf6hgwe067pstqgg2aj7Z memory]#   
[root@iZuf6hgwe067pstqgg2aj7Z memory]#   
[root@iZuf6hgwe067pstqgg2aj7Z memory]# cd memorydemo  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# cat memory.limit_in_bytes  
9223372036854771712  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]#   
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]#   
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]#   
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# echo 104960000 > memory.limit_in_bytes  
[root@iZuf6hgwe067pstqgg2aj7Z memorydemo]# cat memory.limit_in_bytes  
104960000  

---
## 执行malloc程序，观察是否会自动killed
[root@iZuf6hgwe067pstqgg2aj7Z ~]# cd Applications/Go/  
[root@iZuf6hgwe067pstqgg2aj7Z Go]# ll  
total 8  
drwxr-xr-x 3 root root 4096 Feb 28 23:00 pkg  
drwxr-xr-x 4 root root 4096 Feb 28 23:18 src  
[root@iZuf6hgwe067pstqgg2aj7Z Go]# cd src/github.com/  
[root@iZuf6hgwe067pstqgg2aj7Z github.com]# ll  
total 4  
drwxr-xr-x 3 root root 4096 Feb 28 23:18 cncamp  
[root@iZuf6hgwe067pstqgg2aj7Z github.com]# cd cncamp/  
[root@iZuf6hgwe067pstqgg2aj7Z cncamp]# ll  
total 4  
drwxr-xr-x 6 root root 4096 Feb 28 23:08 golang  
[root@iZuf6hgwe067pstqgg2aj7Z cncamp]# cd golang/examples/module3  
[root@iZuf6hgwe067pstqgg2aj7Z module3]# ll  
total 8  
drwxr-xr-x 2 root root 4096 Feb 28 23:08 busyloop  
drwxr-xr-x 3 root root 4096 Feb 28 23:25 malloc  
[root@iZuf6hgwe067pstqgg2aj7Z module3]#   
[root@iZuf6hgwe067pstqgg2aj7Z module3]# cd malloc/  
[root@iZuf6hgwe067pstqgg2aj7Z malloc]# ll  
total 2784  
-rw-r--r-- 1 root root     438 Feb 28 23:08 main.go  
-rw-r--r-- 1 root root      63 Feb 28 23:08 Makefile  
-rwxr-xr-x 1 root root 2832328 Mar  2 22:17 malloc  
-rw-r--r-- 1 root root     211 Feb 28 23:08 malloc.c  
drwxr-xr-x 3 root root    4096 Feb 28 23:08 output  
[root@iZuf6hgwe067pstqgg2aj7Z malloc]#   
[root@iZuf6hgwe067pstqgg2aj7Z malloc]# ./malloc   
Allocating 100Mb memory, raw memory is 104960000  
Allocating 200Mb memory, raw memory is 209920000  
Allocating 300Mb memory, raw memory is 314880000  
**Killed**  
[root@iZuf6hgwe067pstqgg2aj7Z malloc]#   
