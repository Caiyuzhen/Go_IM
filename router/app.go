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

	// 🌟 使用 【ginSwagger】 中间件来生成 API 文档 => API文档化：Swagger可以自动从你的代码生成API文档，并生成可视化界面来调用API，还可以为每个API设置测试用例，方便测试
	docs.SwaggerInfo.BasePath = "" // 🔥 访问 swagger 生成的 API 文档 => http://localhost:8081/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler((swaggerfiles.Handler))) // 表示任何路由都可以访问 swagger


	// 🌟 静态资源（前端）
	router.Static("/asset", "asset/") // 各种静态文件
	router.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	router.LoadHTMLGlob("views/**/*") // 加载 HTML 视图文件 


	// 🌟【路由 API】
	// router.GET("/index", service.GetIndex) // 🌟【http://localhost:8081/index】 首页
	router.GET("/", service.GetIndex) // 🌟【http://localhost:8081/ 首页
	router.GET("/index", service.GetIndex) // 🌟【http://localhost:8081/index 首页
	router.GET("/user/getUserList", service.UserListService) // 🌟 获取用户列表 【http://localhost:8081/user/getUserList】 router 内数据的处理方式放在 Server 层  =>  比如 GetUserList 方法
	router.GET("/user/createUser", service.CreateUser) // 新增用户的接口 => http://localhost:8081/user/createUser?name=Annie&password=123456&rePassword=123456'
	router.GET("/user/deleteUser", service.DeleteUser) // 删除用户的接口 => http://localhost:8081/user/deleteUser?id=1
	router.POST("/user/updateUser", service.UpdateUser) // 更新用户的接口 => http://localhost:8081/user/updateUser
	router.POST("/user/login", service.FindUserByNameAndPassword) // 用户登录的接口 => http://localhost:8081/user/login?name=海绵宝宝&password=123456


	// 🌟 发送 websocket 消息 (Redis)
	router.GET("/user/sendMsg", service.SendMsgServer) // 访问在线测试工具: https://www.easyswoole.com/wstool.html  => 【ws://127.0.0.1:8081/user/sendMsg】
	router.GET("/user/sendUserMsg", service.SendUserMsg) // 发送消息给指定用户 => 【ws://127.0.0.1:8081/user/sendUserMsg】
	return router
}