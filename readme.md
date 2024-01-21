## 初始化
- go mod init ginchat


## 在 Navicat 中建立 localhost 数据库
- utf8mb4



## 整体安装缺失依赖
- go mod download


## 汇总所有 go 的依赖
- go mod tidy


## 安装库
`GORM 是一个流行的Go语言ORM（Object-Relational Mapping，对象关系映射）库, 它提供了一种高效的方式来在Go应用程序与数据库之间进行数据交互`
`GIN 是一个用 Go (Golang) 编写的 HTTP web 框架。它具有高性能的路由器和中间件，这使您能够创建功能全面的 API 和 Web`
`Swagger 允许你使用 YAML 或 JSON 格式来描述你的 API。这种描述被称为 Swagger 规范`
- go get github.com/jinzhu/gorm(旧)
- go get -u gorm.io/gorm(新)
- go get gorm.io/driver/mysql 
- go get -u github.com/gin-gonic/gin
- go get github.com/spf13/viper
- go get -u github.com/swaggo/swag/cd/swag 
- go get -u github.com/swaggo/swag/cmd/swag
  - 📄 文档: https://pkg.go.dev/github.com/swaggo/gin-swagger#section-readme
  - swag init (🔥安装 swag 后记得做这步!)
  - go get -u github.com/swaggo/gin-swagger(🔥安装 swag 后记得做这步!)
  - go get -u github.com/swaggo/files(🔥安装 swag 后记得做这步!)
- go get gorm.io/gorm/logger


## 注入测试数据
`go run testGorm.go`


## 启动项目
`go run main.go`



## References
- GORM
  - https://gorm.io/zh_CN/docs/index.html
- GIN
  - https://pkg.go.dev/github.com/gin-gonic/gin#section-readme