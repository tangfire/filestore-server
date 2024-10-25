# 在docker中配置mysql主从数据同步

## 主服务器
- step1：在docker中创建并启动MySQL主服务器：端口`3306`

```bash
docker run -d \
-p 3306:3306 \
-v /hitd/mysql/master/conf:/etc/mysql/conf.d \
-v /hitd/mysql/master/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=123456 \
--name master \
mysql:8.0.29
```
可以用我这个版本的mysql,其他版本用这个配置方法可能会出问题

- step2：创建MySQL主服务器配置文件：


默认情况下MySQL的binlog日志是自动开启的，可以通过如下配置定义一些可选配置

```bash
vim /hitd/mysql/master/conf/my.cnf
```


配置如下内容

```
[mysqld]
# 服务器唯一id，默认值1
server-id=1
# 设置日志格式，默认值ROW
binlog_format=STATEMENT
# 二进制日志名，默认binlog
# log-bin=binlog
# 设置需要复制的数据库，默认复制全部数据库
#binlog-do-db=mytestdb
# 设置不需要复制的数据库
#binlog-ignore-db=mysql
#binlog-ignore-db=infomation_schema

```

重启MySQL容器

```bash
docker restart master
```


`binlog格式说明：`

- binlog_format=STATEMENT：日志记录的是主机数据库的`写指令`，性能高，但是now()之类的函数以及获取系统参数的操作会出现主从数据不同步的问题。


- binlog_format=ROW（默认）：日志记录的是主机数据库的`写后的数据`，批量操作时性能较差，解决now()或者 user()或者 @@hostname 等操作在主从机器上不一致的问题。


- binlog_format=MIXED：是以上两种level的混合使用，有函数用ROW，没函数用STATEMENT，但是无法识别系统变量




- step3：使用命令行登录MySQL主服务器：

```bash
#进入容器：env LANG=C.UTF-8 避免容器中显示中文乱码
docker exec -it master env LANG=C.UTF-8 /bin/bash
#进入容器内的mysql命令行
mysql -uroot -p
#修改默认密码校验方式 
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';

```


- step4：主机中创建slave用户：

```sql
-- 创建slave用户
CREATE USER 'slave'@'%';
-- 设置密码
ALTER USER 'slave'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
-- 授予复制权限
GRANT REPLICATION SLAVE ON *.* TO 'slave'@'%';
-- 刷新权限
FLUSH PRIVILEGES;


```


- step5：主机中查询master状态：

执行完此步骤后`不要再操作主服务器MYSQL`，防止主服务器状态值变化

记下File值和Position值


## 准备从服务器

可以配置多台从机slave1、slave2…，这里以配置slave1为例

我们这个项目就只配置一个从服务器

- step1：在docker中创建并启动MySQL从服务器：`端口3307`

```bash
docker run -d \
-p 3307:3306 \
-v /hitd/mysql/slave1/conf:/etc/mysql/conf.d \
-v /hitd/mysql/slave1/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=123456 \
--name slave \
mysql:8.0.29

```


- step2：创建MySQL从服务器配置文件：

```bash
vim /hitd/mysql/slave1/conf/my.cnf

```

配置如下内容：

```
[mysqld]
# 服务器唯一id，每台服务器的id必须不同，如果配置其他从机，注意修改id
server-id=2
# 中继日志名，默认xxxxxxxxxxxx-relay-bin
#relay-log=relay-bin


```

重启MySQL容器


```bash
docker restart slave

```

- step3（可选）：使用命令行登录MySQL从服务器：

```bash
#进入容器：
docker exec -it slave env LANG=C.UTF-8 /bin/bash
#进入容器内的mysql命令行
mysql -uroot -p
#修改默认密码校验方式 
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';


```


- step4：在从机上配置主从关系：

在**从机**上执行以下SQL操作,此处所需的值对应第一步master数据库中所操作记录的值

```sql
CHANGE MASTER TO MASTER_HOST='192.168.188.100', 
MASTER_USER='hitd_slave',MASTER_PASSWORD='123456', MASTER_PORT=3306,
MASTER_LOG_FILE='binlog.000003',MASTER_LOG_POS=1348; 


```

MASTER_HOST对应master的IP地址(不知道IP的 可以用下面方法查看ip)

如果你想从宿主机上查看容器的 IP 地址，可以使用 Docker 命令：

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' master

```

将 master 替换为你的容器名称。这将返回容器的 IP 地址。

MASTER_USER对应master数据库所创建给从库使用的用户

MASTER_PASSWORD对应master数据库所创建给从库使用的用户密码

MASTER_LOG_FILE对应指令 master数据库执行SHOW MASTER STATUS; 命令的第一项值(File值)

MASTER_LOG_POS对应指令 master数据库执行SHOW MASTER STATUS; 命令的第二项值(Position值)




## 启动主从同步

启动从机的复制功能，执行SQL：

```sql
START SLAVE;
-- 查看状态（不需要分号）
SHOW SLAVE STATUS\G


```

执行后，主要看这两个参数：
Slave_IO_Running,Slave_SQL_Running 

这两个值都是Yes,则说明主从配置成功！

## 实现主从同步

在主机中执行以下SQL，在从机中查看数据库、表和数据是否已经被同步

```sql
CREATE DATABASE db_user;
USE db_user;
CREATE TABLE t_user (
 id BIGINT AUTO_INCREMENT,
 uname VARCHAR(30),
 PRIMARY KEY (id)
);
INSERT INTO t_user(uname) VALUES('zhang3');
INSERT INTO t_user(uname) VALUES(@@hostname);


```

## 停止和重置

需要的时候，可以使用如下SQL语句

```bash
-- 在从机上执行。功能说明：停止I/O 线程和SQL线程的操作。
stop slave; 

-- 在从机上执行。功能说明：用于删除SLAVE数据库的relaylog日志文件，并重新启用新的relaylog文件。
reset slave;

-- 在主机上执行。功能说明：删除所有的binglog日志文件，并将日志索引文件清空，重新开始所有新的日志文件。
-- 用于第一次进行搭建主从库时，进行主库binlog初始化工作；
reset master;

```



