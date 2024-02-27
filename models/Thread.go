package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

// 群 Model
type Thread struct {
	gorm.Model // 继承 Gorm
	Name string // 群名称
	OwnerId uint // 群主 ID
	Img string // 群头像
	Desc string // 群描述
}






// 创建群的普通房啊 -- 参数则是传入上面 类的实例
func CreateThread(thread Thread) (int, string) {
	if len(thread.Name) == 0 { // 群名称不能为空
		return -1, "❌ 群名称不能为空"
	}

	if thread.OwnerId == 0 { // 群主 ID 不能为空, 后续可以加一些安全性的校验, 拿到登录 token + 用户 ID 去校验这个人是否合法
		return -1, "❌ 请先登录"
	}

	// 把群新建到数据库内
	if err := utils.DB.Create(&thread).Error; err != nil {
		fmt.Println(err)
		return -1, "❌ 创建失败"
	}

	return 0, "✅ 创建成功"
}