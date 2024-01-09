package service

import (
	"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
)

// 处理路由的数据 => 获取用户列表
func UserListService(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserListModel()
	if len(data) == 0 {
        fmt.Println("❓未查询到数据")
    }

	c.JSON(200, gin.H {
		"message": data,
	})
}