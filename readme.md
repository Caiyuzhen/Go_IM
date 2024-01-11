## 初始化
- go mod init ginchat


## 在 Navicat 中建立 localhost 数据库
- utf8mb4


## 整体安装缺失依赖
- go mod tidy


## 安装库
`GORM 是一个流行的Go语言ORM（Object-Relational Mapping，对象关系映射）库, 它提供了一种高效的方式来在Go应用程序与数据库之间进行数据交互`
`GIN 是一个用 Go (Golang) 编写的 HTTP web 框架。它具有高性能的路由器和中间件，这使您能够创建功能全面的 API 和 Web`
- go get github.com/jinzhu/gorm(旧)
- go get -u gorm.io/gorm(新)
- go get gorm.io/driver/mysql 
- go get -u github.com/gin-gonic/gin
- go get github.com/spf13/viper
- go get -u github.com/swaggo/swag/cd/swag


## 注入测试数据
`go run testGorm.go`


## 启动项目
`go run main.go`



## References
- GORM
  - https://gorm.io/zh_CN/docs/index.html
- GIN
  - https://pkg.go.dev/github.com/gin-gonic/gin#section-readme