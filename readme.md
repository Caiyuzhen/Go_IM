## 初始化
- go mod init ginchat


## 在 Navicat 中建立 localhost 数据库
- utf8mb4


## 在 Navicat 中连接 mySQL 数据库
- ginChat
- 123456


## 整体安装缺失依赖
- go mod download


## 汇总所有 go 的依赖
- go mod tidy


## 安装库
`GORM 是一个流行的 Go 语言 ORM（Object-Relational Mapping，对象关系映射）库, 它提供了一种高效的方式来在 【Go 应用程序】与【数据库】之间进行数据交互`
`GIN 是一个用 Go (Golang) 编写的 HTTP web 框架。它具有高性能的路由器和中间件，这使您能够创建功能全面的 API 和 Web`
`Swagger 允许你使用 YAML 或 JSON 格式来描述你的 API。这种描述被称为 Swagger 规范`
- go get github.com/jinzhu/gorm(旧)
- go get -u gorm.io/gorm(新)
- go get gorm.io/driver/mysql 
- go get -u github.com/gin-gonic/gin
  - 安装后需要在 config 的 app.yml 内配置 GinChat 数据库连接信息
- go get github.com/spf13/viper
- go get -u github.com/swaggo/swag/cd/swag 
- go get -u github.com/swaggo/swag/cmd/swag
  - 📄 SWAG 库 文档: https://pkg.go.dev/github.com/swaggo/gin-swagger#section-readme
  - swag init (🔥安装 swag 后记得做这步!)
  - go get -u github.com/swaggo/gin-swagger(🔥安装 swag 后记得做这步!)
  - go get -u github.com/swaggo/files(🔥安装 swag 后记得做这步!)
- go get gorm.io/gorm/logger
- go get github.com/asaskevich/govalidator
  - 检验账号跟密码等格式的库
- go get github.com/redis/go-redis/v9
  - 👍 引入百万级消息并发的缓存库, 用于缓存用户信息, 安装后需要在 config 的 app.yml 内配置 redis 的连接信息
- go get github.com/gorilla/websocket
  - 使用 websocket 来实现聊天功能
- go get gopkg.in/fatih/set.v0
- <!-- - go get github.com/thedevsaddam/govalidator  -->


## 注入测试数据
`go run testGorm.go`


## 更新 Swap 文档（🌟 每次新增接口都需要 init 一下！）
`swag init`


## 启动项目
`go run main.go`
启动后访问`http://localhost:8081/swagger/index.html`可以看到 Swagger 生成的接口文档


## 启动 Ridis 数据库
🌟 `cd` 到 Ridis 文件夹并 `redis-server` => 启动 Ridis 数据库


## 在线 WebSocket 测试工具
- https://www.easyswoole.com/wstool.html


## References
- GORM
  - https://gorm.io/zh_CN/docs/index.html
- GIN
  - https://pkg.go.dev/github.com/gin-gonic/gin#section-readme