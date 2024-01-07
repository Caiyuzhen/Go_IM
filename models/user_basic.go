package models

import (
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

// 类方法
func (table *UserBasic) TableName() string {
	return "user_basic"
}


// 普通方法 => 获取用户数据
func GetUserList() []*UserBasic { // UserBasic 类型指针的切片, 这里的每个元素都是指向 UserBasic 类型的指针, 这意味着可以直接修改这些指针指向的 UserBasic 对象
	return nil
}