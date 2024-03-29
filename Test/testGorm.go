package main

import (
	"time"
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_"fmt"
)

// 【测试】用户数据的 Model  =>  Schema
type UserBasic struct {
	gorm.Model // 继承 gorm.Model, 继承后可以使用 gorm.Model 的属性
	Identity string // 唯一标识
	Name string // 用户名
	Password string // 密码
	Phone string // 手机号
	Email string // 邮箱
	ClientIp string // 客户端 IP => 设备
	ClientPort string // 客户端端口 => 设备
	Salt string
	LoginTime *time.Time // 登录时间(使用指针类型, 让默认值为空), uint64 是时间戳, 使用 time.Time 可以避免为空时默认时间为 0 的状态
	HeartBeatTime *time.Time // 心跳时间(使用指针类型, 让默认值为空),, uint64 是时间戳,  使用 time.Time 可以避免为空时默认时间为 0 的状态
	LogoutTime *time.Time  // 登出时间(使用指针类型, 让默认值为空),, uint64 是时间戳,  使用 time.Time 可以避免为空时默认时间为 0 的状态  || `` 为表达式, 自定义在数据库内的字段名 `gorm:"column:logOut_time" json:"logOut_time`
	IsLogOut bool // 是否登出
	DeviceInfo string // 设备信息
}


func main() {
	// 连接数据库 ————————————————————————————————————————————————————————————————————————————————————————————————
	// 👇 后续封装到 utils 内去连接数据库
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/ginChat?charset=utf8mb4&parseTime=True&loc=UTC"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}


	// 👇 创建【用户】、【消息】、【群组】、【用户关系】几张表
	db.AutoMigrate(&models.UserBasic{}) // 如果没有表, 则创建一张【用户】表  =>  schema
	db.AutoMigrate(&models.MessageBasic{}) // 如果没有表, 则创建一张【消息】表  =>  schema
	db.AutoMigrate(&models.Thread{}) // 如果没有表, 则创建一张【群组】表  =>  schema
	db.AutoMigrate(&models.GroupBasic{}) // 如果没有表, 则创建一张【群组】表  =>  schema
	db.AutoMigrate(&models.ContactBasic{}) // 如果没有表, 则创建一张【用户关系】表  =>  schema

	// 新增一个用户 user ————————————————————————————————————————————————————————————————————————————————————————————————
	currentTime := time.Now() // 使用 time.Now() 获取当前时间
	user := &models.UserBasic{
		Name: "Jimmy",
		LoginTime: &currentTime, 
		HeartBeatTime: nil, // 默认为空
		LogoutTime: nil, // 默认为空
	}
	// user.Name = "Zeno"
	db.Create(user)



	// 读取用户 user  ————————————————————————————————————————————————————————————————————————————————————————————————
	// db.First(user, 1) // 根据整型主键查找 => 设置 id 为 1 的 user
	// fmt.Println(db.First(user, 1))

	// db.First(user, "code = ?", "D42") // 查找 code 字段值为 D42 的记录



	// 修改用户 user (更新) ————————————————————————————————————————————————————————————————————————————————————————————————
	// db.Model(user).Update("Password", 1234)


	// // Update - 更新多个字段
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})



	// Delete - 删除用户 user  ————————————————————————————————————————————————————————————————————————————————————————————————
	// db.Delete(user, 1) // 🔥 别一新建就删除了 !
}