package models

import (
	"fmt"
	"ginchat/utils"
	"time"
	"gorm.io/gorm"
)

// 设计用户数据的 Model  =>  Schema
type UserBasic struct {
	gorm.Model // 继承 gorm.Model, 继承后可以使用 gorm.Model 的属性
	Identity string // 唯一标识
	Name string // 用户名
	Password string // 密码
	Phone string // 手机号
	Email string // 邮箱
	ClientIp string // 客户端 IP => 设备
	ClientPort string // 客户端端口 => 设备
	LoginTime *time.Time // 登录时间(使用指针类型, 让默认值为空), uint64 是时间戳, 使用 time.Time 可以避免为空时默认时间为 0 的状态
	HeartBeatTime *time.Time // 心跳时间(使用指针类型, 让默认值为空),, uint64 是时间戳,  使用 time.Time 可以避免为空时默认时间为 0 的状态
	LogoutTime *time.Time  // 登出时间(使用指针类型, 让默认值为空),, uint64 是时间戳,  使用 time.Time 可以避免为空时默认时间为 0 的状态  || `` 为表达式, 自定义在数据库内的字段名 `gorm:"column:logOut_time" json:"logOut_time`
	IsLogOut bool // 是否登出
	DeviceInfo string // 设备信息
}

// ⚠️ => 类方法
func (table *UserBasic) TableName() string {
	return "user_basic"
}


// 🌟 普通方法 => 获取用户数据 (在 router 内定义一个 url, 然后通过 service 来调用这个 GetUserList 的 models 方法)
func GetUserListModel() []*UserBasic { // UserBasic 类型指针的切片, 这里的每个元素都是指向 UserBasic 类型的指针, 这意味着可以直接修改这些指针指向的 UserBasic 对象
	data := make([]*UserBasic, 10) // 创建一个切片, 用于存放要查询的 userData 数据
	utils.DB.Find(&data) // 查询数据库 => 传入 date, 在所有数据中进行查询
	// var data []*UserBasic // 创建一个空切片 => Find 函数会自动填充切片, 因此不用我们事先声明切片的长度
	result := utils.DB.Find(&data) // 查询数据库, 传入 	&data,  &data 会存放 Find() 后的所有结果

	if result.Error != nil {
        fmt.Println("❌ 数据库查询错误: ", result.Error)
        return nil
    }

	for _, v := range data {
		fmt.Println("✅ 查询到的数据为: ", v) // 单条数据
	}
	return data
}