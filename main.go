package main

import (
	// "github.com/gin-gonic/gin"
	"ginchat/router"
	"ginchat/utils"
)

func main() { // utils 初始化 => route => model => service => ...
	// 初始化配置文件 ________________________________________________
	utils.InitConfig() // 🔥 初始化配置文件 => 从 yml 内引入配置 !!
	utils.InitMySQL()  // 初始化数据库

	// 代码分层后的方式 ________________________________________________
	router := router.Router()
	router.Run(":8081") // listen and serve on localhost:8080 端口

	// 【代码没有分层的方式】建立一個 gin 的router 的示例 ________________________________________________
	// router := gin.Default()
	// router.GET("ping", func(c *gin.Context){ // 路由放在 router 层
	// 	c.JSON(200, gin.H { // 数据的处理放在 service 层
	// 		"message": "pong",
	// 	})
	// })
	// router.Run(":8081")
	// router.Run() // 默认端口
}
