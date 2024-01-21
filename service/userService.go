package service

import (
	_"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
)


// GetIndex
// @Tags 首页
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func UserListService(c *gin.Context) { // 处理路由的数据 => 获取用户列表
	data := make([]*models.UserBasic, 10) // 创建一个切片来承接返回值
	data = models.GetUserListModel()

	c.JSON(200, gin.H {
		"message": data,
	})
}


// CreateUser
// @Tags 首页
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) { // 处理路由的数据 => 获取用户列表
	user := models.UserBasic{}
	user.Name: c.Query("name"), // 获取路由中的 name 参数 => Query 是 gin 框架的方法

	password: c.Query("password"), // 获取路由中的 password 参数 => Query 是 gin 框架的方法
	rePassword: c.Query("password"), // 获取路由中的 rePassword 参数 => Query 是 gin 框架的方法

	if password != rePassword {
		c.JSON(-1, gin.H {
			"message": "❌ 两次输入的密码不一致!",
		})
	}

	user := models.UserBasic{

	}

	c.JSON(200, gin.H {
		"message": data,
	})
}