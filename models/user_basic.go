package models

import (
	_"crypto/md5"
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

// 设计用户数据的 Model  =>  Schema
type UserBasic struct {
	gorm.Model // 继承 gorm.Model, 继承后可以使用 gorm.Model 的属性
	Identity string // 唯一标识
	Avatar string // 头像
	Name string // 用户名
	Password string // 密码
	Phone string `valid:"matches(^1[3-9]{1}\\d{9}$)"`// 手机号 => 使用 ValidateStruct 进行校验
	Email string  `valid:"email"` // 邮箱 => 使用 ValidateStruct 进行校验
	ClientIp string // 客户端 IP => 设备
	ClientPort string // 客户端端口 => 设备
	LoginTime *time.Time // 登录时间(使用指针类型, 让默认值为空), uint64 是时间戳, 使用 time.Time 可以避免为空时默认时间为 0 的状态
	HeartBeatTime *time.Time // 心跳时间(使用指针类型, 让默认值为空),, uint64 是时间戳,  使用 time.Time 可以避免为空时默认时间为 0 的状态
	LogoutTime *time.Time  // 登出时间(使用指针类型, 让默认值为空),, uint64 是时间戳,  使用 time.Time 可以避免为空时默认时间为 0 的状态  || `` 为表达式, 自定义在数据库内的字段名 `gorm:"column:logOut_time" json:"logOut_time`
	IsLogOut bool // 是否登出
	DeviceInfo string // 设备信息
	Salt string // 盐值 => 用于加密
}

// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *UserBasic) TableName() string { // 🔥 table 数据表, TableName 用于指定表名: 这是函数中用来引用传入的 UserBasic 指针的变量名, TableName() 是一个定义在 UserBasic 结构体上的方法，返回一个 string 类型的值
	return "user_basic"
}



// 🌟 普通方法 => 获取用户数据 (在 router 内定义一个 url, 然后通过 service 来调用这个 GetUserList 的 models 方法)
func GetUserListModel() []*UserBasic { // UserBasic 类型指针的切片, 这里的每个元素都是指向 UserBasic 类型的指针, 这意味着可以直接修改这些指针指向的 UserBasic 对象
	userData := []*UserBasic{} //【切片创建方法一来存放实例化的 UserBasic 】创建一个切片（能放一组用户数据）, 用于存放要查询的 userData 数据
	// userData := make([]*UserBasic, 10) //【切片创建方法二】 创建一个切片（能放一组用户数据）, 用于存放要查询的 userData 数据
	// var data []*UserBasic // 【切片创建方法三】创建一个空切片 => Find 函数会自动填充切片, 因此不用我们事先声明切片的长度
	ErrResult := utils.DB.Find(&userData) // 使用 utils 内的 DB 去 Find 查询数据库 => 传入 userData, 在所有数据中进行查询, 🔥 userData 会存放 Find() 后的所有结果！而 ❌ result 则是返回报错！！
	

	if ErrResult.Error != nil {
        fmt.Println("❌ 数据库查询错误: ", ErrResult.Error)
        return nil
    }

	if len(userData) == 0 {
		fmt.Println("❓未查询到数据")
		return nil
	}

	for _, v := range userData {
		fmt.Println("✅ 查询到的单条数据为: ", v) // 单条数据
	}
	return userData // 返回所有数据
}






// 🌟 普通方法 => 通过 【ID】去查找到某个用户
func FindUserByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user) // 👈 从数据库里边根据 id 找到某个用户
	return user
}




// 🌟 普通方法 => 通过【用户名】去查找到某个用户, 后续可以增加比如查询异性、同城等功能
func FindUserByName(name string) UserBasic {  // 返回 DB 内找到的用户数据
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user) // 从 DB 内找到 name = name 的数据, 并返回, 如果返回为空则表示还没有注册这个用户
	return user // 🔥因为传入的是【引用类型的数据】, 所以返回后都是被修改过的！
}


// 🌟 普通方法 => 通过【手机号】去查找到某个用户, 后续可以增加比如查询异性、同城等功能
func FindUserByPhone(phone string) UserBasic {  // 返回 DB 内找到的用户数据
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user) // 从 DB 内找到 name = name 的数据, 并返回, 如果返回为空则表示还没有注册这个用户
	return user // 🔥因为传入的是【引用类型的数据】, 所以返回后都是被修改过的！
}


// 🌟 普通方法 => 通过【邮箱】去查找到某个用户, 后续可以增加比如查询异性、同城等功能
func FindUserByEmail(email string) UserBasic {  // 返回 DB 内找到的用户数据
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user) // 从 DB 内找到 name = name 的数据, 并返回, 如果返回为空则表示还没有注册这个用户
	return user // 🔥因为传入的是【引用类型的数据】, 所以返回后都是被修改过的！
}


// 🌟 普通方法 => 新增用户
func CreateUser(user UserBasic) *gorm.DB { // 返回 DB 内新增的用户数据, *gorm.DB 表示返回的是一个 gorm.DB 类型的指针
	return utils.DB.Create(&user)
}


// 🌟 普通方法 => 删除用户
func DeleteUser(user UserBasic) *gorm.DB { // 返回 DB 内删除的用户数据
	return utils.DB.Delete(&user)
}


// 🌟 普通方法 => 更新用户(修改)
func UpdateUser(user UserBasic) *gorm.DB { // 返回 DB 内的用户数据
	return utils.DB.Model(&user).Updates(UserBasic{
		Name: user.Name,
		Password: user.Password,
		Phone: user.Phone,
		Email: user.Email,
		Salt: user.Salt,
		Avatar: user.Avatar,
	})
}


// 🌟 普通方法 => 用于登录
func FindUserByNameAndPasswordInModel(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? AND password = ?", name, password).First(&user) // 从 DB 内找到用户名跟密码都相同的数据, 并返回, 如果返回为空则表示 ❌ 错误

	// 生成 JWT （JSON Web Tokens）鉴权 token
	str := fmt.Sprintf("%d", time.Now().Unix())// 拿到系统时间戳
	temp :=  utils.MD5Encode(str)// 生成加密的 token
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp) // 当用户 id = user.ID 时, 将 identity 更新为 加密的 token

	return user
}