package service

import (
	_"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
)

// 处理路由的数据 => 获取用户列表
func UserListService(c *gin.Context) {
	data := make([]*models.UserBasic, 10) // 创建一个切片来承接返回值
	data = models.GetUserListModel()

	c.JSON(200, gin.H {
		"message": data,
	})
}