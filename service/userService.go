package service

import (
	_ "fmt"
	"ginchat/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Summary 获取所有用户列表
// @Tags 用户模块
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
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param rePassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) { // 处理路由的数据 => 获取用户列表
	user := models.UserBasic{}
	user.Name = c.Query("name") // 【因为 user 在上面 user := models.UserBasic{} 实例化了, 因此直接 user.Name 】 => 获取路由中的 name 参数 => Query 是 gin 框架的方法

	password := c.Query("password") // 获取路由中的 password 参数 => Query 是 gin 框架的方法
	rePassword := c.Query("rePassword") // 获取路由中的 rePassword 参数 => Query 是 gin 框架的方法

	if password != rePassword {
		c.JSON(-1, gin.H {
			"message": "❌ 两次输入的密码不一致!",
		})
		return
	}

	// 如果不是密码不一致, 则将密码赋值给 user.Password
	user.Password = password
	// 创建用户成功后的返回值
	models.CreateUser(user) // 调用 model 内的方法
	c.JSON(200, gin.H {
		"message": "新增用户成功",
	})
}



// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "用户 id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) { // 处理路由的数据 => 获取用户列表
	user := models.UserBasic{}
	id, err := strconv.Atoi(c.Query("id")) // 👈👈 将路由中的 id 参数转换为 int 类型 => Atoi 是 strconv 包的方法
	if err != nil {
		c.JSON(-1, gin.H {
			"message": "❌ id 参数错误",
		})
		return
	}
	user.ID = uint(id) // 将转换后的 id 赋值给 user.ID  | 🔥 ID 在继承的 gorm 的 class 中有, 为 大写 | ubit 为无符号整型, 表示非负整数的数据类型

	models.DeleteUser(user) // 调用 model 内的方法

	// 删除用户成功后的返回值
	c.JSON(200, gin.H {
		"message": "删除用户成功",
	})
}


// UpdateUser
// @Summary 更新用户数据
// @Tags 用户模块
// @param id formData string false "用户 id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) { // 处理路由的数据 => 获取用户列表
	user := models.UserBasic{}
	id, err := strconv.Atoi(c.PostForm("id")) // 👈👈 将路由中的 id 参数转换为 int 类型 => Atoi 是 strconv 包的方法, 通过  PostForm 来获得数据!!
	if err != nil {
		c.JSON(-1, gin.H {
			"message": "❌ id 参数错误",
		})
		return
	}
	// 🔥拿到 id, 传给下一层的 model 去修改数据库
	user.ID = uint(id) // 将转换后的 id 赋值给 user.ID  | 🔥 ID 在继承的 gorm 的 class 中有, 为 大写 | ubit 为无符号整型, 表示非负整数的数据类型

	// 修改 user 的 name 和 password
	user.Name = c.PostForm("name") // 获取路由中的 name 参数 => PostForm 是 gin 框架的方法
	user.Password = c.PostForm("password") // 获取路由中的 password 参数 => PostForm 是 gin 框架的方法

	models.UpdateUser(user) // 调用 model 内的方法, 👉传入 id

	// 用户用户成功后的返回值
	c.JSON(200, gin.H {
		"message": "修改用户成功",
	})
}