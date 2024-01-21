package utils

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// _ "ginchat/models"
	// "ginchat/models"
)

// 🌟 定义一个全局的 db 变量, 用于接收初始化后的数据库连接
var DB *gorm.DB  // => 在 model 层会调用到 DB 这个全局变量！


// 应用的初始化配置
func InitConfig() { // 用 viper 读取配置文件内的流式数据, viper 为 GORM 内置的方法
	viper.SetConfigName("app") // 设置配置文件名的名称 (不带后缀)
	viper.AddConfigPath("config") // 设置配置文件的路径 => ginChat 是项目的根目录
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		fmt.Printf("❌ viper read config failed, err: %v\n", err)
	}
	fmt.Println("⚙️ 正在初始化 mySQL 的配置文件...")
	fmt.Println("✅ viper 读取到了 config 的配置文件(数据库路由): ", viper.Get("mysql")) // 打印得到的内容 => map[dns:root:123456@tcp(127.0.0.1:3306)/ginChat?charset=utf8mb4&parseTime=True&loc=UTC]
}



// 传入 【初始化配置】以连接数据库
func InitMySQL() {
	// 自定义日志模板, 打印查询数据库的 SQL 语句 => 方便调试
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io 流
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值, 默认是 100ms
			LogLevel: logger.Info, // 级别
			Colorful: true, // 是否彩色打印
		},
	)
	var err error
	// 打开数据库连接
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{
		Logger: newLogger, // 使用自定义的日志模板
	})
	fmt.Println("⚙️ 正在连接数据库...")

	if err != nil {
    fmt.Printf("❌ 数据库连接失败: %v\n", err)
	fmt.Println("✅ 连接数据库成功")
    return
}

	// 👇 测试下, 后续这个查询的动作就放在 model 层的 GetUserList 方法
	// 【方法一】使用 gorm 封装的查询语句, 新建一个人, 查询一个人________________________________________________________________________
	// userData := models.UserBasic{}
	// DB.Find(&userData)


	// 【方法二】使用 gorm 封装的查询语句, 新建一个切片, 查询其中一个人________________________________________________________________________
	// userData := []*models.UserBasic{}  // 👈 这样无法查询到数据
	// userData := make([]*models.UserBasic, 10) 
	// DB.Find(&userData[0])
	// fmt.Println("✅ 连接数据库成功, 数据库内的数据为: ", userData)


	//【方法三】同上也是新建切片，不过查询的是一组人
	// userData := []*models.UserBasic{}
	// DB.Find(&userData)
	// fmt.Println("✅ 连接数据库成功, 数据库内的数据为: ", userData)


	// 【方法四】使用 mySQL 查询语句________________________________________________________________________
	// result := DB.Raw("SELECT * FROM user_basic").Scan(&userData)
	// if result.Error != nil {
	// 	fmt.Println("❌ 执行原始查询时出错: ", result.Error)
	// }
	// fmt.Println("✅ 执行原始查询成功, 数据库内的数据为: ", userData)
}

