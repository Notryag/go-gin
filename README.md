# golang gin框架学习


运行go的文件直接`go run main.go`
## 快速开始

go在go 1.13以上的时候有一个工作区的概念, 可以创建不同的模块, 有点类似mono的不同服务的概念

工作区
```shell
go work init
mkdir common && cd common
# 添加工作区
go work use ./common
```


```shell
go work init
mkdir common && cd common
# mod初始化
go mod init
go work use ./common
# 添加mod
# 安装gin -u表示拉去网络最新版本的包
go get -u github.com/gin-gonic/gin
```
