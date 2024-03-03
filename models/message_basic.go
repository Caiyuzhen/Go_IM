package models // 消息的结构
import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"ginchat/utils" // 引入 utils 内的方法, 在下面通过 InitConfig 进行调用
	// "github.com/google/uuid" // 给消息生成唯一的 id, 用于标识消息的唯一性, 避免重复发消息
	"github.com/gorilla/websocket"
	// "github.com/redis/go-redis"
	"github.com/redis/go-redis/v9"
	// "github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)


// 消息结构体 => 定义后可以去 testGorm.go 中去生成一张表
type MessageBasic struct {
	gorm.Model        // 继承 Gorm
	UUID       string 	  // 消息的唯一标识符
	UserId     int64  // 消息发送者 ID
	TargetId   int64  // 消息接收者 ID
	Type       int    // 消息类型 (1.私聊、2.群聊、3广播(比如欢迎加入 XXX 群聊))  => 用 1 2 3 来表示
	Media      int    // 消息媒体类型 (1.文本、2.图片、3.表情包、4.音频、5.视频、6.文件)  =>  后续可以扩展出红包、名片等更多类型
	Content    string // 消息内容
	Pic        string // 图片地址
	Audio      string // 音频地址
	Url        string // 链接地址
	ReadTime   uint64 //读取时间
	Desc       string // 描述
	Amount     int    // 文件大小等其他数字统计类型
	CreateTime uint64 // 创建时间
}

// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *MessageBasic) TableName() string { // TableName 为数据表, 用于指定表名
	return "message_basic" // 在 db 中的表名
}


// 🚀 关系节点的结构体, 包含用户关系、消息数据以及群组 ——————————————————————————————————————————————————————————————————————————————————————————————————————————————
type Node struct {
	Conn          *websocket.Conn // 🚀 客户端的 WebSocket 连接, 用于与客户端通信  => 用户的连接数据, 用于发送消息, 知道要发送给谁
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          // 💗 用户的心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     // 🔥 消息 (一个管道, 用于存放待发送给客户端的消息)
	GroupSets     set.Interface   // ⚡️  好友 / 群 => 使用 set 库存储该客户端所加入的群组的集合, 可以构造更安全的线程
}



var clientMap map[int64]*Node = make(map[int64]*Node, 0) // 用于存储用户的连接信息  =>  🔥 存放映射关系（绑定用户 ID 和 Node）的全局变量  =>  存储所有连接到服务器的客户端节点, 键是客户端的唯一标识符（如用户 ID），值是对应的 Node 结构体实例
var rwLocker sync.RWMutex // 读写锁
var udpSendChan_SaveMsgFromUDP chan []byte = make(chan []byte, 1024) // 🌟 全局变量, 用来保存【UDP 协议接收到的广播消息内容】 => 然后可以在下面的广播消息中进行调用  =>  用于存放消息的管道, 1024 表示最多存放 1024 个消息
// var processedMsgIDs = make(map[string]bool)// 全局变量, 记录已发送的消息 ID
// var nn int = 0 // debug 看函数执行了几次



// 【🔥 聊天需要的字段 - 前端需要发送（发送者 ID、接收者 ID 、消息类型、发送的内容、登录 token 校验）】聊天室的总的公共方法(处理客户端连接请求的函数, 当客户端尝试建立 WebSocket 连接时会创建一个 Node 实例, 将其添加到 clientMap 中, 并启动发送（sendProc）和接收（receiveProc）协程) => 单聊、群聊、广播都需要获取一些参数等等 -> 发送消息, 需要 【发送者 ID】、【接收者 ID】、【消息类型】、【消息内容】
func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query() //【☝️第一步】 从 URL 中获取参数
	Id := query.Get("userId")                 // 是 string 类型, 但是上面的 clientMap 是 int64 类型, 因此需要转换一下数据格式
	userId, _ := strconv.ParseInt(Id, 10, 64) // 10 表示十进制, 64 表示 int64 类型
	isValida := true                          // 临时变量, 用于校验参数是否合法, 后续传入数据库进行校验 checkToken(token)

	conn, err := (&websocket.Upgrader{ //【☝️第二步】升级为 websocket 并校验请求来源, 防止跨域攻击
		// 校验 Token (能否聊天)
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("❌ 升级为 websocket 失败", err)
		return
	}

	currentTime := uint64(time.Now().Unix()) //【第三步】初始化 node 来获取用户关系 Conn
	node := &Node{
		Conn:          conn,                       // 客户端的 WebSocket 连接, 用于与客户端通信
		Addr:          conn.RemoteAddr().String(), //客户端地址
		LoginTime:     currentTime,                //首次连接时间
		HeartbeatTime: currentTime,                // 💗 用户的心跳时间
		DataQueue:     make(chan []byte, 50),      //  一个管道, 用于存放待发送给客户端的数据 => 初始化 50 个消息
		GroupSets:     set.New(set.ThreadSafe),    //  一个集合, 用于存储该客户端所加入的群组 => 初始化一个线程安全的 set 集合
	}

	// *****   //【第四步】判断用户关系
	
	rwLocker.Lock()          //【第五步】将 userId 跟 node 进行绑定并【加锁】, 用于后续的消息推送udpReceiveProc_Podcast
	clientMap[userId] = node // 将 userId 跟 node 进行绑定, 建立关系, 用于后续的消息推送
	rwLocker.Unlock()        // 解锁

	go sendProc_WebsocketMsg_Personal(node) //【第六步】调用消息发送的方法 (☎️☎️ 给他人发!)  =>  从管道中取出数据

	go receiveProc_WebsocketMsg_Personal(node) //【第七步】调用接收消息的方法  (☎️☎️ (给自己发）!)   =>   接收消息的协程

	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(4)*time.Hour) // 【第八步】把在线用户的消息加到 Redis 缓存中  =>  SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
	// sendMsg_ToME(userId, []byte("🚀 欢迎加入聊天室")) // 连接后, 默认给前端发送一条消息
}





// ************************************************************************************************************************************************






// 【初始化广播协程, 自动执行】Go 语言会在程序启动时自动执行该（初始化函数） => 在这里，它用于启动处理 UDP 📢 广播消息发送（udpSendProc）和 📢 接收广播消息（udpReceiveProc）的协程（数据的发送与接收读取） ——————————————————————————————————————————————————————————————————————————————————————————————————————————————
func init() {
	go DUP_SendProc_BroadcastMessage()    // 调度发送消息的协程
	go DUP_ReceiveProc_BroadcastMessage() // 调度接收消息的协程
}



// 发送广播消息的具体方法
// 【📢 广播消息到局域网内的方法】用于处理 UDP 广播消息的发送, 从 udpSendChan 通道中读取消息, 并通过 UDP 协议将这些消息广播到局域网内
func DUP_SendProc_BroadcastMessage() { // 👈 也可以用来广播群消息
	fmt.Println("🌟🌟🌟🌟🌟🌟🌟🌟🌟🌟到这里🌟🌟🌟🌟🌟🌟🌟🌟🌟🌟🌟")
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{ // DialUDP 为 net 包中的方法, 用于发送 udp 数据
		IP:   net.IPv4(192, 168, 0, 255), // 广播到局域网内, 传入以太网 ip (路由的网关地址)  // Port: viper.GetInt("port.udpPort"),
		Port: 3000,
	})
	defer con.Close() // 关闭连接, 避免内存泄漏
	if err != nil {
		fmt.Println("❌ 广播消息失败", err)
	}
	for {
		select {
		case data := <- udpSendChan_SaveMsgFromUDP: // 🌟 从广播到局域网内的消息内【取出数据】
			// fmt.Println("🛜 广播消息到局域网 >>>>>> ", string(data))
			fmt.Println("🛜 广播消息到局域网 >>>>>> ")
			_, err := con.Write(data) // 写入消息
			if err != nil {
				fmt.Println("❌ 广播消息失败", err)
				return
			}
		}
	}
}





// 【📢 接收广播消息】, 责监听 UDP 广播消息, 当局域网内有消息广播时, 这个协程会接收到这些消息并进行获取
func DUP_ReceiveProc_BroadcastMessage() { // 👈 也可以用来广播群消息	
	con, err := net.ListenUDP("udp", &net.UDPAddr{ // ListenUDP 为 net 包中的方法, 用于接收 udp 数据
		IP:   net.IPv4(192, 168, 0, 255),                 // IPv4ero  (0,0,0,0)  => 表示所有 ip 端口都可以接受
		Port: 3000, // 写死的端口号  // Port: viper.GetInt("port.udpPort"), // 配置在 app.yml 的端口号
	})
	defer con.Close() // 关闭连接

	if err != nil {
		fmt.Println("❌ 接收广播消息失败", err)
	}
	
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:]) // 读取管道中的消息数据
		if err != nil {
			fmt.Println("❌ 读取管道消息失败", err)
			return
		}

		// fmt.Println("🛜 接收到了局域网内的广播消息 >>>>>> ", string(buf[0:n]))
		fmt.Println("🛜 接收到了局域网内的广播消息 >>>>>> ")
		fmt.Println("————————————————————————————————————————————————————————")
		fmt.Println(" ")
		DispatchMsg_Podcast(buf[0:n]) // 👈👈 调度消息的调度逻辑
	}
}







// ************************************************************************************************************************************************






//【把消息转发给谁的 (调度）】=> 判断要把拿到的局域网消息分发消息到【单聊】、【群聊】还是【系统消息】等, 看业务需求******************************************************************
func DispatchMsg_Podcast(data []byte) {
	// testData := []byte(`{"FromId": 1, "TargetId": 2, "Type": 1, "Content": "测试消息"}`)
	msg := MessageBasic{}                     

	// 解析数据, 因为 data 是二进制字符数据, 需要解析成 Json 数据（反序列化）
	err := json.Unmarshal(data, &msg)          
	if err != nil {
		fmt.Println("❌ 解析 JSON 消息失败", err)
		return
	}

	// 检查数据是否为有效的 JSON 格式
	if utils.IsValidJSON(data) {
		switch msg.Type { // 🌟 根据消息类型进行分发
			case 1: // 私聊
				fmt.Println(" ")
				fmt.Println("————————————————————————————————————————————————————————")
				fmt.Println("******【💬 第一步】开始分发消息给自己  >>>>>> ", msg.Content)
				sendMsg_ToME(msg.TargetId, data) // 发送【📢 广播消息 - 把消息转发给另外一个人】的具体方法, 如果后续消息量大可能要【做缓存】、【做集群】
			case 2: // 群聊
				sendGroupMsg(msg.TargetId, data) // 群发消息
		}
	} else {
		fmt.Println("❌ 数据不是有效的 JSON 格式", string(data)) // 数据不是有效的JSON，进行其他处理, 比如是 String 类型的数据
		return
	}
}




// ************************************************************************************************************************************************




// 👇 发送 Websocket 广播消息的具体方法
//	发送【websocketMsg_Persona 双向消息】的方法 (从管道中取出数据) => 这条调用了后, 接收方（对方）才能收到消息!
func sendProc_WebsocketMsg_Personal(node *Node) {
	for {
		select {
		case data := <- node.DataQueue: // 从管道中获取数据 🔥
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 写入 Conn 来发送消息
			if err != nil {
				fmt.Println("❌ 发送消息失败", err)
				return
			}
			fmt.Println("******【💬 第三步】从管道中获取数据")
		}
	}
}



// 接收【websocketMsg_Persona 双向消息】（发送方也会接收到自己发送的消息, 接收到数据后可以广播给其他地方）
func receiveProc_WebsocketMsg_Personal(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 接收消息, 返回值有三个, 第一个是消息类型, 第二个是消息内容, 第三个是错误信息
		if err != nil {
			fmt.Println("❌ 接收(自己发送的)消息失败", err)
			return
		}
		msg := MessageBasic{}
		err = json.Unmarshal(data, &msg) // 解析数据（反序列化接收到的数据）, 因为 data 是二进制数据, 需要解析成结构体
		if err != nil {
			fmt.Println("❌ 解析 JSON 消息失败", err)
		}

		// 心跳检测
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.UpdateUserHeartbeat(currentTime)
		} else {
			DispatchMsg_Podcast(data)   // 分发消息 -> sendMsg_ToME
			save_broadCastMsg(data) // 把消息保存到管道
			fmt.Println("******【💬 第四步】接收(自己发送的)消息成功  >>>>>> ", msg.Content)
			fmt.Println("————————————————————————————————————————————————————————")
			fmt.Println(" ")
		}
	}
}




// 拿到前端发来的消息, 存入管道, 并存入 Redis
func sendMsg_ToME(TargetId int64, originalMsg []byte) { // 传入 userId 和 msg
	rwLocker.RLock()              // 加锁 => 读锁
	node, ok := clientMap[TargetId] // 获取用户的连接信息, 用于发送消息
	rwLocker.RUnlock()

	// 【⭕️ zRedis 缓存 - 1】 前期处理, 消息序列化等工作
	jsonMsg := MessageBasic{}
	json.Unmarshal(originalMsg, &jsonMsg) // 👈 Unmarshal 用来编码 json 数据, 将 JSON 字符串反序列化到 MessageBasic结 构体实例
	
	ctx := context.Background()   // 创建一个空的 context.Context 对象, 用来初始化
	targetIdStr := strconv.Itoa(int(TargetId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix()) // 拿到时间戳


	res, err := utils.RedisDB.Get(ctx, "online_" + userIdStr).Result() // 比如 online_23, 表示用户 23 的在线状态
	if err != nil {
		fmt.Println("❌ 没有找到历史消息", err)
	}
	if res != "" {
		if ok {
			fmt.Println("******【💬 第二步】后台发送消息了 >>>>>> 消息接收者: ", TargetId, " 消息内容:", string(jsonMsg.Content))
			node.DataQueue <- originalMsg // 把消息内容加入管道, 然后再去给下面的 📢 广播消息 进行发送(到某个人)
		}
	}
	var key string
	if TargetId > jsonMsg.UserId { // 表示消息的递增 (确定两个用户 ID 的顺序, 如果你有两个用户ID，1 和 2，不管是用户1发送给用户 2 的消息, 还是用户 2 发送给用户 1 的消息, 都希望它们能够存储在同一个 Redis 键下)
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	// fmt.Println("🌟🌟🌟 测试: ", "userIdStr:", userIdStr, "targetIdStr:", targetIdStr) // 打印检查


	// 【⭕️ zRedis 缓存 - 2】 真正去做消息的缓存
	res2, err2 := utils.RedisDB.ZRevRange(ctx, key, 0, -1).Result() // 先查询下缓存的消息, 看下怎么排序
	if err2 != nil {
		fmt.Println("❌ 查询缓存消息失败!", err2)
	}

	score := float64(cap(res2)) + 1 // 根据长度 + 1, 表示消息的递增
	r, e := utils.RedisDB.ZAdd(ctx, key, redis.Z{score, originalMsg}).Result() // 添加消息缓存
	if e != nil {
		fmt.Println("❌ Redis 缓存添加失败", e)
	}
	fmt.Println("✅ Redis 缓存添加成功!", r)	 // r 表示添加成功的数量
}




// ***********************************************************************************************************************************************




// 进行消息保存的方法 (写入管道)
func save_broadCastMsg(data []byte) {
	udpSendChan_SaveMsgFromUDP <- data // 把数据加入管道, 然后再去给下面的 📢 广播消息 进行发送
}




// ************************************************************************************************************************************************



// 🚀 群发消息的方法
func sendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("✈️ 开始群发消息")
	userIds := SearchUserByGroupId(uint(targetId)) // 根据群内的用户 id 找到用户
	for i := 0; i < len(userIds); i++ {
		if targetId != int64(userIds[i]) { // 排除给自己消息
			sendMsg_ToME(int64(userIds[i]), msg)
		}
	}
}




// ************************************************************************************************************************************************




// 将 msg 转为 byte 类型 (类方法)
func (originalMsg MessageBasic) MarshalBinary() ([]byte, error) {
	return json.Marshal(originalMsg)
}



// 获取 Redis 缓存里边的消息
func RedisMsgModel(userIdA int64, userIdB int64, start int64, end int64, isRevRange bool) []string {
	rwLocker.RLock()
	// node, ok := clientMap[userIdA]
	rwLocker.RUnlock()

	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))

	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	// fmt.Println("🌟🌟 开始查询 Redis 缓存消息, 拼接的 redis key 是: ", key)

	var rels []string
	var err error

	if isRevRange { // 判断是否倒序
		rels, err = utils.RedisDB.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.RedisDB.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println("❌ 没有找到 Redis 消息缓存", err)
	}
	return rels
}


// 💗 更新用户的心跳时间 (🔥类方法)
func (node *Node) UpdateUserHeartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}


// 检测用户心跳是否超时 (🔥类方法)
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) { // 返回 timeout 这个参数（ bool 类型）
	// if node.HeartbeatTime+viper.GetUint64("heartbeat.timeout") < currentTime { // ⚡️ app.yml 内配置的超时时间
		if node.HeartbeatTime + 100000 < currentTime {  // ⚡️如果大于了 100000, 则表示超时
		fmt.Println("🛜 检测到心跳超时了, 即将关闭 socket 连接...", node)
		timeout = true
	}
	return
}


// 🧹 清理超时连接
func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("❌ 清理超时连接失败", r)
		}
	}()
	currentTime := uint64(time.Now().Unix())
	for i := range clientMap { // 表示遍历 clientMap, 用于清理超时连接
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("💔 心跳超时, 已关闭 socket 连接...", node)
			node.Conn.Close() // 关闭连接
		}
	}
	return result
}
























