package service
import (
	"github.com/gin-gonic/gin"
)


// 路由的数据处理
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H {
		"message": "Welcome to Gin Server!",
	})
}