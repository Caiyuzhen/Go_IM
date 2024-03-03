package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

// 群 Model
type Thread struct {
	gorm.Model        // 继承 Gorm
	Name       string // 群名称
	OwnerId    uint   // 群主 ID
	Img        string // 群头像
	Desc       string // 群描述
}


// 创建群的普通方法-- 参数则是传入上面 类的实例
func CreateThread(thread Thread) (int, string) { // 返回 int 跟 string
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


// 显示群列表的方法
func LoadThreadModel(ownerId uint) ([]*Thread, string) { // 返回 【群数据集合】 跟 【string】
	// threadData := make([]*Thread, 10)
	var findalThreads []*Thread // 存放用户所创建的群
	var joinedThreads []Thread // 存放用户所加入的群
	var contactBasics []ContactBasic

	// // 首先，获取用户创建的所有群组  =>  去数据库中查询群列表
	utils.DB.Where("owner_id=?", ownerId).Find(&findalThreads) // 查询条件是 ownerId, 也就是过滤出属于谁的群


	// 其次，获取用户加入的所有群组的关系记录
	utils.DB.Where("owner_id=? AND type=?", ownerId, 2).Find(&contactBasics)
	if len(contactBasics) > 0 { // 如果记录 > 0, 说明用户加入了某些群
		for _, contact := range contactBasics {
			var thread Thread
			utils.DB.Where("id=?", contact.TargetId).First(&thread)
			// fmt.Println(contact) // 打印群的集合数据
			joinedThreads = append(joinedThreads, thread) // 把查询出来的群数据放到 threads 切片里边
		}
	}

	// 将用户加入的群组添加到最终的群组列表中
	for _, getJoinThread := range joinedThreads {
		findalThreads = append(findalThreads, &getJoinThread)
	}
	return findalThreads, "✅ 群列表查询成功"
}


// 添加群组
func JoinThreadModel(userId uint, threadId string) (int, string) {
	contact := ContactBasic{} // 创建一个 ContactBasic 的实例
	contact.OwnerId = userId
	contact.Type = 2

	thread := Thread{}

	// 👇 通过 id 去查找群
	utils.DB.Where("id=? or name=?", threadId, threadId).Find(&thread) // 🔥🔥【第一步】把 threadId 传入 thread 实例内
	if thread.Name == "" {
		return -1, "❌ 群不存在"
	}

	// 👇 通过 id、targetId、 类型 去判断是否加过群了
	utils.DB.Where("owner_id=? and target_id=? and type=2", userId, threadId).Find(&contact) // 通过数据库去查找某个人的群, 过滤出 contact
	if !contact.CreatedAt.IsZero() { // 如果 contact.CreatedAt 不为空, 就说明已经加入过群了
		return -1, "❌ 已经加入过群"
	} else {
		contact.TargetId = thread.ID // 🔥🔥【第二步】把查到的 threadId 传入到 contact 关系 model 里边  =>  建立 【哪个人 userId】 跟 【哪个群 threadId】 的关系  =>  【userId】 与 【threadId】
		utils.DB.Create(&contact) // 创建一条数据, 表示【某个人】与【某个群】的关系 => 加群成功
		return 0, "✅ 成功加入群聊"
	}
}