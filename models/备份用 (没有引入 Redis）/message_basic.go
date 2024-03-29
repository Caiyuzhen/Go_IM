package models // 消息的结构
import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"ginchat/utils" // 引入 utils 内的方法, 在下面通过 InitConfig 进行调用

	"github.com/gorilla/websocket"
	// "github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)


/*
	sendProc_websocketMsg_Personal   		🌟 把消息发送给特定某人
	receiveProc_websocketMsg_Personal 		🌟 接收自己发送的消息

	broadCastMsg_BeenSave 					保存消息的方法
	udpSendProc_Podcast						【📢 广播消息到局域网内的方法】用于处理 UDP 广播消息的发送, 从 udpSendChan 通道中读取消息, 并通过 UDP 协议将这些消息广播到局域网内
	udpReceiveProc_Podcast 					【📢 接收广播消息】, 责监听 UDP 广播消息, 当局域网内有消息广播时, 这个协程会接收到这些消息并进行获取
	dispatchMsg_Podcast 					【📢 把消息转发给谁的调度】的调度逻辑 => 判断要把拿到的局域网消息分发消息到【单聊】、【群聊】还是【系统消息】等, 看业务需求

	sendMsg_Podcast 						🌟 拿到前端发来的消息, 存入管道

*/


// 消息结构体 => 定义后可以去 testGorm.go 中去生成一张表
type MessageBasic struct {
	gorm.Model        // 继承 Gorm
	FromId     int64  // 消息发送者 ID
	TargetId   int64  // 消息接收者 ID
	Type       int    // 消息类型 (1.私聊、2.群聊、3广播(比如欢迎加入 XXX 群聊))  => 用 1 2 3 来表示
	Media      int    // 消息媒体类型 (1.文本、2.图片、3.表情包、4.音频、5.视频、6.文件)  =>  后续可以扩展出红包、名片等更多类型
	Content    string // 消息内容
	Pic        string // 图片地址
	Audio      string // 音频地址
	Url        string // 链接地址
	Desc       string // 描述
	Amount     int    // 文件大小等其他数字统计类型
}

// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *MessageBasic) TableName() string { // TableName 为数据表, 用于指定表名
	return "message_basic" // 在 db 中的表名
}

// 🚀 关系节点的结构体, 包含用户关系、消息数据以及群组 ——————————————————————————————————————————————————————————————————————————————————————————————————————————————
type Node struct {
	Conn      *websocket.Conn // 🚀客户端的 WebSocket 连接, 用于与客户端通信  => 用户的连接数据, 用于发送消息, 知道要发送给谁
	DataQueue chan []byte     // 🔥 一个管道, 用于存放待发送给客户端的数据
	GroupSets set.Interface   // ⚡️  一个集合, 用于存储该客户端所加入的群组 => 使用 set 库 来存储用户所在的群组, 可以构造更安全的线程
}

// 🔥 一个全局变量, 用于存储所有连接到服务器的客户端节点, 键是客户端的唯一标识符（如用户 ID），值是对应的 Node 结构体实例  =>  存放映射关系（绑定用户 ID 和 Node）
var clientMap map[int64]*Node = make(map[int64]*Node, 0) // 用于存储用户的连接信息

// 读写锁
var rwLocker sync.RWMutex // 读写锁

// 【🔥🔥 聊天需要的字段 - 前端需要发送（发送者 ID、接收者 ID 、消息类型、发送的内容、登录 token 校验）】聊天室的总的公共方法(处理客户端连接请求的函数, 当客户端尝试建立 WebSocket 连接时会创建一个 Node 实例, 将其添加到 clientMap 中, 并启动发送（sendProc）和接收（receiveProc）协程) => 单聊、群聊、广播都需要获取一些参数等等 -> 发送消息, 需要 【发送者 ID】、【接收者 ID】、【消息类型】、【消息内容】
func Chat(writer http.ResponseWriter, request *http.Request) {
	//【第一步】 从 URL 中获取参数
	query := request.URL.Query()

	// 获取 Chat 路由内的参数
	// token := query.Get("token")
	Id := query.Get("userId")                   // 是 string 类型, 但是上面的 clientMap 是 int64 类型, 因此需要转换一下数据格式
	userId, err := strconv.ParseInt(Id, 10, 64) // 10 表示十进制, 64 表示 int64 类型
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	// msgType := query.Get("type")
	isValida := true // 临时变量, 用于校验参数是否合法, 后续传入数据库进行校验 checkToken(token)

	//【第二步】升级为 websocket 并校验请求来源, 防止跨域攻击
	conn, err := (&websocket.Upgrader{
		// 校验 Token (能否聊天)
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println("❌ 升级为 websocket 失败", err)
		return
	}

	//【第三步】初始化 node 来获取用户关系 Conn
	node := &Node{
		Conn:      conn,                    // 客户端的 WebSocket 连接, 用于与客户端通信
		DataQueue: make(chan []byte, 50),   //  一个管道, 用于存放待发送给客户端的数据 => 初始化 50 个消息
		GroupSets: set.New(set.ThreadSafe), //  一个集合, 用于存储该客户端所加入的群组 => 初始化一个线程安全的 set 集合
	}

	//【第四步】判断用户关系

	//【第五步】将 userId 跟 node 进行绑定并【加锁】, 用于后续的消息推送
	rwLocker.Lock()          // 加锁
	clientMap[userId] = node // 将 userId 跟 node 进行绑定, 建立关系, 用于后续的消息推送
	rwLocker.Unlock()        // 解锁

	//【第六步】调用消息发送的方法
	go sendProc_websocketMsg_Personal(node) // 从管道中取出数据

	//【第七步】调用接收消息的方法
	go receiveProc_websocketMsg_Personal(node)   // 接收消息的协程
	sendMsg_Podcast(userId, []byte("🚀 欢迎加入聊天室")) // 连接后, 默认给前端发送一条消息
}



// 👇 发送消息的具体方法  ——————————————————————————————————————————————————————————————————————————————————————————————————————————————
//	(🌟  第四步) 发送【websocketMsg_Persona 双向消息】的方法 (从管道中取出数据) => 这条调用了后, 接收方（对方）才能收到消息!
func sendProc_websocketMsg_Personal(node *Node) {
	for {
		select {
		case data := <-node.DataQueue: // 从管道中获取数据 🔥
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 发送消息
			if err != nil {
				fmt.Println("❌ 发送消息失败 (sendProc_websocketMsg_Personal)", err)
				return
			}
			fmt.Println("📮 【第四步】发送消息成功 (sendProc_websocketMsg_Personal) >>>>>>", string(data))
			fmt.Println("————————————————————————————————————————————————————————")
		}
	}
}

// (🌟 第一步) 接收消息并写入管道
func receiveProc_websocketMsg_Personal(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 接收消息, 返回值有三个, 第一个是消息类型, 第二个是消息内容, 第三个是错误信息
		if err != nil {
			fmt.Println("❌ 接收(自己发送的)消息失败 (receiveProc_websocketMsg_Personal)", err)
			return
		}
		broadCastMsg_BeenSave(data) // 🔥 把消息保存到 udpSendChan_SaveMsg 上
		fmt.Println("✅ 【第一步】接收(自己发送的)消息成功 (receiveProc_websocketMsg_Personal) >>> ", string(data))
	}
}


// 🌟 全局变量, 用来保存【消息内容】 => 然后可以在下面的广播消息中进行调用
var udpSendChan_SaveMsg chan []byte = make(chan []byte, 1024) // 用于存放消息的管道, 1024 表示最多存放 1024 个消息


// 进行消息保存的方法 (写入管道)
func broadCastMsg_BeenSave(data []byte) {
	udpSendChan_SaveMsg <- data // 把数据加入管道, 然后再去给下面的 📢 广播消息 进行发送
}




// 【初始化广播协程, 自动执行】Go 语言会在程序启动时自动执行该（🌟 初始化函数）。在这里，它用于启动处理 UDP 📢 广播消息发送（udpSendProc）和 📢接收广播消息（udpReceiveProc）的协程（数据的发送与接收读取） ——————————————————————————————————————————————————————————————————————————————————————————————————————————————
func init() {
	go udpSendProc_Podcast()    // 广播消息到局域网的协程
	go udpReceiveProc_Podcast() // 从局域网中接收消息的协程
}




// 【📢 广播消息到局域网内的方法】用于处理 UDP 广播消息的发送, 从 udpSendChan 通道中读取消息, 并通过 UDP 协议将这些消息广播到局域网内
func udpSendProc_Podcast() { // 👈 也可以用来广播群消息
	// fmt.Println("🌟🌟🌟🌟🌟🌟🌟🌟🌟🌟到这里🌟🌟🌟🌟🌟🌟🌟🌟🌟🌟🌟")
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{ // DialUDP 为 net 包中的方法, 用于发送 udp 数据
		IP:   net.IPv4(192, 168, 0, 255), // 广播到局域网内, 传入以太网 ip (路由的网关地址)
		// Port: viper.GetInt("port.udpPort"),
		Port: 3000,
	})

	if err != nil {
		fmt.Println("❌ 广播消息失败 (udpSendProc_Podcast)", err)
		return
	}

	defer con.Close() // 关闭连接, 避免内存泄漏

	for {
		select {
		case data := <- udpSendChan_SaveMsg: // 🌟 从广播到局域网内的消息取出数据
			fmt.Println("✅ 【前置一】广播消息到局域网 (udpSendProc_Podcast) >>>>>> ", string(data))
			_, err := con.Write(data) // 写入消息
			if err != nil {
				fmt.Println("❌ 广播消息失败 (udpSendProc_Podcast)", err)
				return
			}
		}
	}
}

// 【📢 接收广播消息】, 责监听 UDP 广播消息, 当局域网内有消息广播时, 这个协程会接收到这些消息并进行获取
func udpReceiveProc_Podcast() { // 👈 也可以用来广播群消息
	con, err := net.ListenUDP("udp", &net.UDPAddr{ // ListenUDP 为 net 包中的方法, 用于接收 udp 数据
		IP:   net.IPv4zero,             // IPv4ero  (0,0,0,0)  => 表示所有 ip 端口都可以接受
		// Port: viper.GetInt("port.udpPort"), // 配置在 app.yml 的端口号
		Port: 3000, // 写死的端口号
	})

	if err != nil {
		fmt.Println("❌ 接收广播消息失败 (udpReceiveProc_Podcast)", err)
		return
	}

	defer con.Close() // 关闭连接

	for {
		var buf [512]byte
		n, err := con.Read(buf[0:]) // 读取管道中的消息数据
		if err != nil {
			fmt.Println("❌ 读取管道消息失败", err)
			return
		}
		fmt.Println("接收到了局域网内的广播消息 (udpReceiveProc_Podcast) >>>>>> ", string(buf[0:n]))
		dispatchMsg_Podcast(buf[0:n]) // 读取消息的调度逻辑
	}
}


// (🌟 第三步) 拿到前端发来的消息, 存入管道
func sendMsg_Podcast(userId int64, msg []byte) { // 传入 userId 和 msg
	fmt.Println("🚀 【第三步】后台发送消息了 (sendMsg_Podcast) >>>>>> 消息发送者: ", userId, " 消息内容:", string(msg))
	rwLocker.RLock()              // 加锁 => 读锁
	node, ok := clientMap[userId] // 获取用户的连接信息, 用于发送消息
	rwLocker.RUnlock()            // 解锁 => 读锁

	if ok {
		node.DataQueue <- msg // 把消息加入管道, 然后给下面的 sendProc_websocketMsg_Personal 进行判断, 如果是单聊消息则发送给对应的用户
	}
}


// (🌟 第二步)【把消息转发给谁的 (调度）】=> 判断要把拿到的局域网消息分发消息到【单聊】、【群聊】还是【系统消息】等, 看业务需求******************************************************************
func dispatchMsg_Podcast(data []byte) {
	// 首先检查数据是否为有效的 JSON 格式
	if utils.IsValidJSON(data) {
		// testData := []byte(`{"FromId": 1, "TargetId": 2, "Type": 1, "Content": "测试消息"}`)
		msg := MessageBasic{}             // 初始化消息结构体
		err := json.Unmarshal(data, &msg) // 解析数据, 因为 data 是二进制数据, 需要解析成结构体
		if err != nil {
			fmt.Println("❌ 解析 JSON 消息失败", err)
			return
		} else {
			fmt.Println("✅ 【第二步】解析静态JSON成功", msg)
		}

		// 🌟 根据消息类型进行分发
		switch msg.Type {
		case 1: // 私聊
			sendMsg_Podcast(msg.TargetId, data) // 发送【📢 广播消息 - 把消息转发给另外一个人】的具体方法, 如果后续消息量大可能要【做缓存】、【做集群】
		// case 2: // 群聊
			// sendGroupMsg(msg)
		// case 3: // 广播
			// sendAllMsg(msg)
		// case 4:
			// return
		}
	} else {
		// 数据不是有效的JSON，进行其他处理, 比如是 String 类型的数据
		fmt.Println("❌ 数据不是有效的 JSON 格式", string(data))

		// 返回错误
		return
	}
}
