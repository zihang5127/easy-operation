# EasyOperation

[![Build Status](https://travis-ci.org/zihang5127/easy-operation.svg?branch=master)](https://travis-ci.org/zihang5127/easy-operation)
[![Build status](https://ci.appveyor.com/api/projects/status/tpm2k23umrqri2dd/branch/master?svg=true)](https://ci.appveyor.com/project/zihang5127/easy-operation/branch/master)

一个基于 Golang 开发的自动化部署系统。

界面和开发思路参考于 [SmartWebHook](https://github.com/lifei6671/go-git-webhook.git)。

# 如何使用？


## 源码安装使用

**1、拉取源码**

```
git clone https://github.com/zihang5127/easy-operation.git

``` 

**2、配置**

- 系统的配置文件位于 conf/app.conf 中：

```ini
sessionon = true
httpport = 8080
#数据库配置
db_host=127.0.0.1
db_port=3306
db_database=easyoperation
db_username=xxxxxx
db_password=xxxxxx
queue_size=50
```


**3、编译**

```
#更新依赖
go get -d ./...

#编译项目
go build -v -tags "pam" -ldflags "-w"
```


**4、运行**

```
#恢复数据库，请提前创建一个空的数据库
./easy-operation orm syncdb

#创建管理员账户
./easy-operation install -username=admin -password=123456 -email=512796048@qq.com

#启动
./easy-operation
```
# 访问

- localhost:8080

# 注意

## 添加 SSH Server

- 当添加的是一台 SSH 方式的服务器时，Server IP 为SSH外网IP地址，端口号为SSH端口号，账号为登录账号，SSH Private Key 可以是密码也可以是登录密钥，系统会自动识别类型。

## 添加 Project

- Project 的 Shell 脚本命令 **暂** 不支持换行，建议用服务器 Shell 脚本代替命令,也可使用 ‘&&’ 拼接命令。
- 当前的 SSH 连接不支持环境变量，需要 *export* 指定，或写绝对路径。