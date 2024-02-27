package service

import (
	"fmt"
	"ginchat/utils" // 引入 utils 内的方法
	"io"
	"math/rand"
	"os"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传图片的接口
// 比如访问 http://localhost:8081/asset/upload/1708961431868671549.jpg
func Upload(c *gin.Context) {
	writer := c.Writer                           // 获取响应对象
	req := c.Request                             // 获取请求对象
	srcFile, header, err := req.FormFile("file") // 返回【文件 | 头部信息 | 报错】
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}

	// 1.后台也要判断传过来的文件类型  2.文件存在 upload 文件夹下
	mimeType := header.Header.Get("Content-Type")
	var suffix string
	switch mimeType {
    case "audio/mpeg":
        suffix = ".mp3"
    case "image/png":
        suffix = ".png"
    case "image/jpeg":
        suffix = ".jpg"
    // 可以根据需要添加更多的MIME类型和对应的文件后缀
    default:
        suffix = ".jpg"
    }
	// suffix := ".png"             // 文件后缀

	// ofileName := header.Filename // 拿到文件名称
	// fmt.Println("👀👀 后台拿到了前端传来的文件名称: ", ofileName)

	// backName := strings.Split(ofileName, ".") // 通过 . 分割文件的后缀
	// if len(backName) > 1 {
	// 	suffix = "." + backName[len(backName)-1] // 表示文件的后缀
	// }

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix) // 格式化文件名称, %d 拿到时间戳, %04d 拿到随机数, %s 拿到文件后缀

	// 保存文件到服务器上
	dirFile, err := os.Create("./asset/upload/" + fileName) // 创建文件
	if err != nil {
		utils.RespFail(writer, err.Error()) // RespFail 为自己封装的方法
		return
	}

	// IO 流
	_, err = io.Copy(dirFile, srcFile) // 拷贝文件到服务器上
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}
	if err != nil {
		utils.RespFail(writer, err.Error())
		return
	}

	domain := "http://localhost:8081" // 服务器域名
	url := domain + "/asset/upload/" + fileName

	fmt.Println("📁 文件上传成功", url)

	utils.RespOK(writer, url, "✅ 上传成功") // Resp 为自己封装的方法)

}
