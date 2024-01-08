package service
import (
	"github.com/gin-gonic/gin"
)


// 处理路由的数据服务 => 初始化
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H {
		"message": "Welcome to Gin Server!",
	})
}