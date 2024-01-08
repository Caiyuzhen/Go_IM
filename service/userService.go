package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
)

// 处理路由的数据 => 获取用户列表
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H {
		"message": data,
	})
}