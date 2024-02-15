package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"

	// "github.com/thedevsaddam/govalidator"
	"github.com/asaskevich/govalidator"
)

// GetAllUserList
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
func CreateUser(c *gin.Context) { // 处理路由的数据 => 🌟 注册用户
	user := models.UserBasic{} // 实例化一个 UserBasic 类型的 user 对象
	user.Name = c.Query("name") // 【因为 user 在上面 user := models.UserBasic{} 实例化了, 因此直接 user.Name 】 => 获取路由中的 name 参数 => Query 是 gin 框架的方法
	password := c.Query("password") // 获取路由中的 password 参数 => Query 是 gin 框架的方法
	rePassword := c.Query("rePassword") // 获取路由中的 rePassword 参数 => Query 是 gin 框架的方法

	salt := fmt.Sprintf("%06d", rand.Int31()) // 🔥🔥 表示生成一个 6 位的随机数, 因为 Sprintf 返回的是一个格式化的字符串, 而 rand.Int31() 返回的是一个 int32 类型的随机数, 因此需要使用 %06d 来格式化

	data := models.FindUserByName(user.Name) // 调用 model 内的方法来查找同名用户, 如果 FindUserByName 返回为空则表示还没有注册这个用户

	// 判断是否已经有同名的注册用户
	if data.Name != "" { // model 内的 FindUserByName 会返回 userr
		c.JSON(-1, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ 用户名已存在",
			"data": "",
		})
		return
	}
	

	if password != rePassword {
		c.JSON(-1, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ 两次输入的密码不一致",
			"data": "",
		})
		return
	}

	// 如果不是密码不一致, 则将密码赋值给 user.Password
	// user.Password = password // 简单的暴力赋值, 不安全
	user.Password = utils.MakePassword(password, salt) //【🔥🔥🔥 设置到数据库内!】调用生成加密值的方法, 传入【密码】与【盐值】来生成更安全的密码
	user.Salt = salt //【🔥🔥🔥 设置到数据库内!】
	
	fmt.Println("🔐🔐🔐 加密后的密码为: ", user.Password)


	// 创建用户成功后的返回值
	models.CreateUser(user) // 调用 model 内的方法
	c.JSON(200, gin.H {
		"code": 0, // 更好的返回值格式, 0 表示成功, -1 表示失败
		"message": "✅ 新增用户成功",
		"data": user, // 返回新增了谁
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
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ id 参数错误",
			"data": "",
		})
		return
	}
	user.ID = uint(id) // 将转换后的 id 赋值给 user.ID  | 🔥 ID 在继承的 gorm 的 class 中有, 为 大写 | ubit 为无符号整型, 表示非负整数的数据类型

	models.DeleteUser(user) // 调用 model 内的方法

	// 删除用户成功后的返回值
	c.JSON(200, gin.H {
		"code": 0, // 更好的返回值格式, 0 表示成功, -1 表示失败
		"message": "✅ 删除用户成功",
		"data": user, // 返回删除了谁
	})
}


// UpdateUser
// @Summary 更新用户数据
// @Tags 用户模块
// @param id formData string false "用户 id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "手机号"
// @param email formData string false "邮箱"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) { // 处理路由的数据 => 获取用户列表 👆（每次更新参数都需要 swag ini 一下!!）
	user := models.UserBasic{}
	id, err := strconv.Atoi(c.PostForm("id")) // 👈👈 将路由中的 id 参数转换为 int 类型 => Atoi 是 strconv 包的方法, 通过 PostForm (🔥 是 Gin 库内置的方法) 来获得数据!!
	if err != nil {
		c.JSON(-1, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ id 参数错误",
			"data": "",
		})
		return
	}
	// 🔥拿到 id, 传给下一层的 model 去修改数据库
	user.ID = uint(id) // 将转换后的 id 赋值给 user.ID  | 🔥 ID 在继承的 gorm 的 class 中有, 为 大写 | ubit 为无符号整型, 表示非负整数的数据类型

	// 👇 修改 user 的 name 或 password 或 phone 或 email
	user.Name = c.PostForm("name") // 获取路由中的 name 参数 => PostForm 是 gin 框架的方法
	// user.Password = c.PostForm("password") // 获取路由中的 password 参数 => PostForm 是 gin 框架的方法
	user.Phone = c.PostForm("phone") // 获取路由中的 phone 参数 => PostForm 是 gin 框架的方法
	user.Email = c.PostForm("email") // 获取路由中的 email 参数 => PostForm 是 gin 框架的方法

	// 生成新的盐值和加密密码 ————————————————————
	plainPassword := c.PostForm("password") // 获取前端传来的原始密码
	salt := fmt.Sprintf("%06d", rand.Int31()) // 🔥🔥 表示生成一个 6 位的随机数, 因为 Sprintf 返回的是一个格式化的字符串, 而 rand.Int31() 返回的是一个 int32 类型的随机数, 因此需要使用 %06d 来格式化
	encryptedPassword := utils.MakePassword(plainPassword, salt) // 加密密码
	user.Password = encryptedPassword // 给 user 实例传入加密后的密码, 再在下面传入 Model 层去修改数据库
	user.Salt = salt // 给 user 实例传入盐值, 再在下面传入 Model 层去修改数据库


	_, err2 := govalidator.ValidateStruct(user) // 使用 govalidator 内的 ValidateStruct 方法来验证 user 的数据是否符合要求
	if err2 != nil {
		c.JSON(-1, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ 数据格式错误, 修改失败!",
			"data": "",
		})
		return
	} else {
		// 👉 调用 model 内的方法, 传入 user, 修改数据库
		models.UpdateUser(user) 
		c.JSON(200, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "✏️ 修改用户成功", // 用户用户成功后的返回值
			"data": user, // 返回修改了谁
		})
	}
}



// Login
// @Summary 登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/FindUserByNameAndPassword [post]
func FindUserByNameAndPassword(c *gin.Context) { // 处理用户登录的路由服务
	data := models.UserBasic{}

	userInputName := c.Query("name") // 拿到用户输入的用户名
	userInputPwd := c.Query("password")  // 拿到用户输入的密码

	// 先从数据库内找到用户
	user := models.FindUserByName(userInputName) 
	dataBaseUserPassword := user.Password // 拿到数据库内的加密密码
	if user.Name == "" { // 不能用 Identity 来校验用户是否存在, 因为 Identity 经常变
		c.JSON(200, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ 用户不存在!",
			"data": "",
		})
		return
	}
	// fmt.Println("😄 找到了用户: ", user)
	// fmt.Println("😄 用户输入的密码: ", userInputPwd)
	// fmt.Println("😄 找到了用户的盐值: ", user.Salt)
	// fmt.Println("😄 找到了用户的加密密码: ", dataBaseUserPassword) // user.Password 是加密后的密码

	// 👆上面通过 name 拿到用户后, 拿到用户的【盐值】跟【用户所输入的密码】并进行 md5 的解密
	flag := utils.ValidPassword(userInputPwd, user.Salt, dataBaseUserPassword)// user.Password 是加密后的密码, 因为在数据库内的密码是加密过的, 因此这里需要解密后才能查询
	if !flag { // 如果密码不正确, !flag 表示 flag 为 false
		c.JSON(200, gin.H {
			"code": -1, // 更好的返回值格式, 0 表示成功, -1 表示失败
			"message": "❌ 密码错误!",
			"data": "",
		})
		return
	}

	// 解密密码 -> 因为数据库内储存的是 🔐 加密后的密码, 所以要重新加密再去数据库进行比对
	pwd := utils.MakePassword(userInputPwd, user.Salt)
	data = models.FindUserByNameAndPasswordInModel(userInputName, pwd) // 🔥 需要传入解密后的密码！！

	c.JSON(200, gin.H { // 密码正确的返回值
		"code": 0, // 更好的返回值格式, 0 表示成功, -1 表示失败
		"message": "✅ 登录成功",
		"data": data,
	})
}
