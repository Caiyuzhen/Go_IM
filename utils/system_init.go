package utils

import (
	"fmt"
	"ginchat/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var DB *gorm.DB // 定义一个全局的 db 变量, 用于接收初始化后的数据库连接


// 应用的初始化配置
func InitConfig() { // 用 viper 读取配置文件内的流式数据
	viper.SetConfigName("app") // 设置配置文件名的名称 (不带后缀)
	viper.AddConfigPath("config") // 设置配置文件的路径 => gincgat 是项目的根目录
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		fmt.Printf("❌ viper read config failed, err: %v\n", err)
	}
	fmt.Println("✅ viper read config success, 读取的信息为: ", viper.Get("app")) // 打印得到的内容
	fmt.Println("✅ viper read config success, 读取的信息为: ", viper.Get("mysql")) // 打印得到的内容
}



// 连接数据库的初始化配置
func InitMySQL() {
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	user := models.UserBasic{}
	DB.Find(&user)
	fmt.Println("✅ 连接数据库成功, 数据库内的数据为: ", user)
}

