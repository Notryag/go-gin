# golang gin框架学习


运行go的文件直接`go run main.go`
## 快速开始

go在go 1.13以上的时候有一个工作区的概念, 可以创建不同的模块, 有点类似mono的不同服务的概念

在代码中已经使用的情况下可以直接用`go mod tidy`来下载依赖

工作区
```shell
go work init
mkdir common && cd common
# 添加工作区
go work use ./common
```


```shell
# mod初始化
go mod init
# 添加mod
# 安装gin -u表示拉去网络最新版本的包
go get -u github.com/gin-gonic/gin
```

## 数据库
https://gorm.io/ gorm官网

启动mysql的docker容器
```shell
docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
```

需要安装gorm 和 mysql的驱动
```shell
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```



## 记录需要学习知识点
RandomString 方法, 查看其中的`make(byte[], n)`, `rand.Intn()`
安装jwt v4
```shell
go get -u github.com/dgrijalva/jwt-go/v4
```