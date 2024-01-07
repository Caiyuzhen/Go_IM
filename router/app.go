package router
import (
	"ginchat/service" // 因为 go mod 初始化的名字是 ginchat, 所以这里要用 ginchat/service!!
	"github.com/gin-gonic/gin"
)


func Router() *gin.Engine { // 返回值 *gin.Engin e是一个指向 Gin 框架的核心引擎的指针, 在Gin框架中, gin.Engine 是处理所有请求的主要结构体
	router := gin.Default() // 🚀 router 是 gin.Engine 的实例
	router.GET("/index", service.GetIndex) // 🌟 router 内数据的处理方式放在 Server 层  =>  比如 GetIndex 方法
	return router
}