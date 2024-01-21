package router
import (
	"ginchat/service" // 🌟 因为 go mod 初始化的名字是 ginchat, 所以这里要用 ginchat/service!!
	"github.com/gin-gonic/gin"
	"ginchat/docs"
	swaggerfiles "github.com/swaggo/files" // swaggerfiles 表示 swagger 的别名, 🌟 引入后还需要去 service 内去写 @ 注解!! 写完后还需要 swag init !!
	ginSwagger "github.com/swaggo/gin-swagger" // ginSwagger 表示 swagger 的别名, 🌟 引入后还需要去 service 内去写 @ 注解!! 写完后还需要 swag init !!
)


func Router() *gin.Engine { // 返回值 *gin.Engin e是一个指向 Gin 框架的核心引擎的指针, 在Gin框架中, gin.Engine 是处理所有请求的主要结构体
	router := gin.Default() // 🚀 router 是 gin.Engine 的实例

	// 使用 ginSwagger 中间件来生成 API 文档 => API文档化：Swagger可以自动从你的代码生成API文档，并生成可视化界面来调用API，还可以为每个API设置测试用例，方便测试
	docs.SwaggerInfo.BasePath = "" // 🔥 访问 swagger 生成的 API 文档 => http://localhost:8081/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler((swaggerfiles.Handler))) // 表示任何路由都可以访问 swagger

	router.GET("/index", service.GetIndex) // 🌟【http://localhost:8081/index】 router 内数据的处理方式放在 Server 层  =>  比如 GetIndex 方法
	router.GET("/user/getUserList", service.UserListService) // 🌟 【http://localhost:8081/user/getUserList】 router 内数据的处理方式放在 Server 层  =>  比如 GetUserList 方法

	return router
}