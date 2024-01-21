package service
import (
	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) { // 处理路由的数据服务 => 初始化, 并且配置 Swagger 文档 !! 👉 配置完后在 cli 中输入 swag init 生成 docs 文件夹
	c.JSON(200, gin.H {
		"message": "Welcome to Gin Server!",
	})
}
