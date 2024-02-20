package service

import (
	_"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) { // 处理路由的数据服务 => 初始化, 并且配置 Swagger 文档 !! 👉 配置完后在 cli 中输入 swag init 生成 docs 文件夹

	
	// 👇 使用 Gin 提供的的方法来渲染 html
	c.HTML(http.StatusOK, "index.html", nil) // 接收状态码、模板文件名、传给模板的数据

	

	// 👇 使用模板方法渲染 html
	// index, err := template.ParseFiles("index.html")
	// if err != nil {
	// 	panic(err) // panic 表示程序发生了不可恢复的错误, 会导致程序中断
	// }
	// index.Execute(c.Writer, "index") // 执行模板的渲染, 并将渲染后的结果返回给客户端


	// 👇 测试用的假的首页, 还没引入 view html 文件, 先返回一些内容看能不能跑通
	// c.JSON(200, gin.H {
		// "message": "Welcome to Gin Server!",
	// })
}
