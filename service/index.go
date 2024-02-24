package service

import (
	"fmt"
	"ginchat/models"
	"html/template"
	_ "net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) { // 处理路由的数据服务 => 初始化, 并且配置 Swagger 文档 !! 👉 配置完后在 cli 中输入 swag init 生成 docs 文件夹
	// 👇 使用 Gin 提供的的方法来渲染 html
	// c.HTML(http.StatusOK, "index.html", nil) // 接收状态码、模板文件名、传给模板的数据


	// 👇 使用模板方法渲染 html
	index, err := template.ParseFiles("index.html", "views/chat/head.html") // 因为 index.html 里面引入了 head.html （拆分出去的模板）, 所以要一起解析
	if err != nil {
		panic(err) // panic 表示程序发生了不可恢复的错误, 会导致程序中断
	}
	index.Execute(c.Writer, "index") // 执行模板的渲染, 并将渲染后的结果返回给客户端


	// 👇 测试用的假的首页, 还没引入 view html 文件, 先返回一些内容看能不能跑通
	// c.JSON(200, gin.H {
		// "message": "Welcome to Gin Server!",
	// })
}



// 跳转到注册页面
func ToRegister(c *gin.Context) {
	fmt.Println("👍 跳转到注册页面")
	// c.HTML(http.StatusOK, "user/register.html", nil) // 接收状态码、模板文件名、传给模板的数据
	ind, err := template.ParseFiles("views/user/register.html") // 👈 使用 Gin 提供的的方法来渲染 html
	if err != nil {
		fmt.Println("❌ 解析模板文件失败: ", err)
	}
	ind.Execute(c.Writer, nil) // 渲染模板文件
}


// 登录后的路由跳转
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles(  // 👈 使用 Gin 提供的的方法来渲染 html (目标 html 有哪些 {{}} 就得写哪些 html!!)
		"views/chat/index.html", 
		"views/chat/head.html", 
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/foot.html",
	)
	if err != nil {
		fmt.Println("❌ 解析模板文件失败: ", err)
	}
	// 从登录页面（登录接口会返回这个信息）获取到的 userId 和 token, 可以拿到是因为【前端】登录页面的表单提交后会带上这两个参数
	// / 🔥 从登录成功的接口中拿到返回值 =>  🔥 把登录成功的返回值给到跳转 chat 的接口 （详见最外层的 index.htmll 根文件）
	userId, _ := strconv.Atoi(c.Query("userId")) // user Id 需要为 uint 类型, 因此使用 Atoi 转换为数字类型, 返回的第二个值 _ 是错误信息, 这里忽略了
	userToken := c.Query("token")
	fmt.Println("✅ 拿到了 userId: ", userId)

	// 把拿到的值存入 user
	user := models.UserBasic{}
	user.ID = uint(userId) // user Id 需要为 uint 类型, uint 为无符号类型, 只能存正整数和零
	user.Identity = userToken

	// 这里理论上还要校验一下 token 是否正确, 但是这里先不做了
	// http://localhost:8081/toChat?userId=0&token=3AA49DB7A7AFDEF91E2115A754405C5E
	ind.Execute(c.Writer, user) // 渲染模板文件
}


