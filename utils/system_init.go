package utils // 导出为 utils 包

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// _ "ginchat/models"
	// "ginchat/models"
)

// 🌟 定义一个全局的 db 变量, 用于接收初始化后的数据库连接 ————————————————————————————————————————————————
var (
	DB *gorm.DB  // => 在 model 层会调用到 DB 这个全局变量！
	RedisDB * redis.Client // => 在 model 层会调用到 RedisDB 这个全局变量！
)


// 应用的初始化配置 ————————————————————————————————————————————————
func InitConfig() { // 🌟 用 viper 读取配置文件内的流式数据, viper 为 GORM 内置的方法, 用于读取配置文件
	viper.SetConfigName("app") // 设置配置文件名的名称 (不带后缀)
	viper.AddConfigPath("config") // 设置配置文件的路径 => ginChat 是项目的根目录
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		fmt.Printf("❌ viper read config failed, err: %v\n", err)
	}
	fmt.Println("⚙️ 正在初始化 mySQL 的配置文件...")
	fmt.Println("✅ viper 读取到了 config 的配置文件(数据库路由): ", viper.Get("mysql")) // 打印得到的内容 => map[dns:root:123456@tcp(127.0.0.1:3306)/ginChat?charset=utf8mb4&parseTime=True&loc=UTC]
}




// 初始化 【读取 app.yml 的 MySQL 配置】以连接数据库 ————————————————————————————————————————————————
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
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{ // 🌟 用 viper 读取配置文件内的流式数据, viper 为 GORM 内置的方法, 用于读取配置文件, GetString 为读取字符串数据
		Logger: newLogger, // 使用自定义的日志模板
	})
	fmt.Println("⚙️ 正在连接数据库...")

	if err != nil {
    fmt.Printf("❌ 数据库连接失败: %v\n", err)
	fmt.Println("✅ 连接数据库成功")
    return
}

	// 👇 测试下, 后续这个查询的动作就放在 model 层的 GetUserList 方法
	// 【方法一】使用 gorm 封装的查询语句, 新建一个人, 查询一个人
	// userData := models.UserBasic{}
	// DB.Find(&userData)


	// 【方法二】使用 gorm 封装的查询语句, 新建一个切片, 查询其中一个人
	// userData := []*models.UserBasic{}  // 👈 这样无法查询到数据
	// userData := make([]*models.UserBasic, 10) 
	// DB.Find(&userData[0])
	// fmt.Println("✅ 连接数据库成功, 数据库内的数据为: ", userData)


	//【方法三】同上也是新建切片，不过查询的是一组人
	// userData := []*models.UserBasic{}
	// DB.Find(&userData)
	// fmt.Println("✅ 连接数据库成功, 数据库内的数据为: ", userData)


	// 【方法四】使用 mySQL 查询语句
	// result := DB.Raw("SELECT * FROM user_basic").Scan(&userData)
	// if result.Error != nil {
	// 	fmt.Println("❌ 执行原始查询时出错: ", result.Error)
	// }
	// fmt.Println("✅ 执行原始查询成功, 数据库内的数据为: ", userData)
}



// 初始化 【🌟 cd 到 Ridis 文件夹并 redis-server => 启动 Ridis 数据库, 然后读取 app.yml 的 Redis 配置】以连接数据库 ————————————————————————————————————————————————
func InitRedis() {
	RedisDB = redis.NewClient(&redis.Options {
		Addr: viper.GetString("redis.Addr"), // 读取配置的数据
		Password: viper.GetString("redis.PassWord"), // 读取配置的数据
		DB: viper.GetInt("redis.DB"), // 读取配置的数据
		PoolSize: viper.GetInt("redis.PoolSize"), // 读取配置的数据
		MinIdleConns: viper.GetInt("redis.MinIdleConns"), // 读取配置的数据
	})

	pong, err := RedisDB.Ping(context.Background()).Result() // 调试 redis 数据库的连接, 从 go-redis 版本 8 开始，Ping 方法需要一个 context.Context 参数, 用于【没有】上下文的情况下的 Redis 连接
	if err != nil {
		fmt.Printf("❌ Redis 数据库初始化失败...: %v\n", err)
	} else {
		fmt.Printf("✅ Redis 数据库初始化成功...: %v\n", pong)
	}
}


// WebSocket + Redis 的消息通讯功能 ————————————————————————————————————————————————
const (
	PublishKey = "websocket" // 发布消息的 key, 用来形成【管道】
)

// 发布消息到 Redis 的 WebSocket
func PublishMsgToRedis(ctx context.Context, channel string, msg string) error { // ctx 表示请求过来的东西, channel 为管道(效率更高一些), message 为消息 ｜ 返回值为 error 是因为可能会出现错误
	// var err error
	err := RedisDB.Publish(ctx, channel, msg).Err() // Publish() 为 redis 的发布消息的方法, Err() 为 ridis 的捕获错误的方法
	if err != nil {
		fmt.Println("❌ 发布消息到 Redis 的 WebSocket 失败...: ", err)
		return err
	}
	fmt.Println("✅ 发布消息到 Redis 的 WebSocket 成功...: ", msg)
	return err
}

// 订阅 Redis 消息的 WebSocket 推送（可以打印到控制台）
func SubMsgToRedis(ctx context.Context, channel string) (string, error) { // 🌟 ctx 表示前端请求过来的东西, channel 为管道(效率更高一些), message 为消息 ｜ 返回值为订阅的【消息字符串】跟【错误】
	fmt.Println("👍 接收到了前端传来的 ctx: ", ctx)
	sub := RedisDB.Subscribe(ctx, channel) // 订阅消息, Subscribe 为 rdis 的订阅消息的方法
	msg, err := sub.ReceiveMessage(ctx) // ReceiveMessage 为 ridis 的储存订阅的消息的方法

	if err != nil {
		fmt.Println("❌ 订阅 Redis 的 WebSocket 失败...: ", err)
		return "❌ Error", err
	}
	fmt.Println("✅ 订阅 Redis 的 WebSocket 成功...: ", msg.Payload)
	return msg.Payload, err // 🌟 Payload 为转化为 ridis 内把【消息】转化为【字符串】的方法
}

