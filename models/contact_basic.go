package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

// 人员关系表
type ContactBasic struct { // (contact_basic 这个 model 表示的是一组关系, owner_id 表示这个好友是谁的, target_id 好友是谁, 比如 24 是 3 owner 的好友)
	gorm.Model // 继承 Gorm
	OwnerId  uint   // 关系拥有者 ID, 类型 uint 要与 MessageBasic 中的 FromId、ToId 一致
	TargetId uint   // 关系的目标 ID, 类型 uint 要与 MessageBasic 中的 FromId、ToId 一致
	Type     int    // 关系类型 (好友、群、关注、粉丝、黑名单)  => 用 1 2 3 来表示 (1: 好友, 2: 群组, 3: 黑名单), 后面可以扩展 4: 粉丝等等预留字段
	Desc     string // 描述信息(备注、标签), 预留字段
}

// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *ContactBasic) TableName() string { // TableName 为数据表, 用于指定表名
	return "contact_basic" // 在 db 中的表名
}






// 查找【某个人】的好友 (contact_basic 这个 model 表示的是一组关系, owner_id 表示这个好友是谁的, target_id 好友是谁, 比如 24 是 3 owner 的好友)
// 也可以直接 SQL 查询  =>  SELECT * FROM `user_basic` WHERE id in (20,21) AND `user_basic`.`deleted_at` IS NULL
func SearchFriend(userId uint) []UserBasic { // 传入 userID, 返回好友的具体信息
	contacts := make([]ContactBasic, 0) // ContactBasic 类型的切片, 用来储存一组好友
	objIDS := make([]uint64, 0)         // uint 类型的切片, 用来存储好友的 ID, 然后再在下面的 for 循环中去取出好友

	// 通过数据库去查找这个人的好友, 过滤 contact
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts) // owner_id 表示某个人的好友, type = 1 写死是好友关系的类型

	// 取出好友
	for _, v := range contacts {
		fmt.Println("✅ 查到了好友 ID:", v.TargetId)      // 打印出好友的 ID
		objIDS = append(objIDS, uint64(v.TargetId)) // 把好友的 ID 存储到 objIDS 切片中
	}

	users := make([]UserBasic, 0)                  // 用来存储好友的信息
	utils.DB.Where("id in ?", objIDS).Find(&users) // 通过好友的 ID 去查找好友的信息, 需要使用 In 查询, 取一定的范围
	fmt.Println("✅ 查到了好友的信息:", users)

	return users // 返回好友的信息
}



// 😄 添加好友 -  通过 ID 添加好友 (好友是双向的, A 加了 B, A 同时也被 B 加了)
func AddFriend(userId uint, targetId uint) (int, string) { // 返回数字 + 字符串  =>  比如 0 + "添加成功", -1 + "添加失败"
	user := UserBasic{} // 创建一个 user 的实例

	if targetId != 0 { // 如果没传入目标用户的 id
		fmt.Println("👍 拿到了前端传来的 userID: ", userId, "跟 targetID: ", targetId)
		user = FindUserByID(targetId) // 传入要找的 id, 找到某个用户

		if user.Salt != "" { // 如果要添加的好友不为空 (判断 Identity 或 Salt 不为空都行)

			// 判断不能自己加自己为好友
			if userId == user.ID {
				return -1, "❌ 不能添加自己为好友"
			}

			// 不能添加已经加过的好友
			contact := ContactBasic{} // 创建一个 ContactBasic 的实例
			utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, targetId).First(&contact) // 通过数据库去查找这个人的好友, 过滤 contact
			if contact.ID != 0 { // 如果 contact.ID 不为空, 就说明已经添加过好友了 (因为在联系人表中有这个人)
				return -1, "❌ 不能重复添加好友"
			}

			// 【事物】GORM 的【事务】可以保证数据的一致性 （比如一张表要同时写入两次), 【事务】默认是开启的
			tx := utils.DB.Begin() // 💼 开启事务 *************

			defer func() { // 处理事务中如果出错了, 就会自动回滚
				if r := recover(); r != nil {
					tx.Rollback() // 💼 回滚事务 *************
				}
			}()

			contact2 := ContactBasic{}
			contact2.OwnerId = userId
			contact2.TargetId = targetId
			contact2.Type = 1          // ContactBasic 结构体的定义, 加好友, 类型为 1
			if err := utils.DB.Create(&contact2).Error; err != nil {  //【⚡️ 传入实例】, 新建一条数据表的数据
				tx.Rollback() // 💼 回滚事务 *************
				return -1, "❌ 好友添加失败"
			}

			contact3 := ContactBasic{}
			contact3.OwnerId = targetId
			contact3.TargetId = userId
			contact3.Type = 1
			if err := utils.DB.Create(&contact3).Error; err != nil {  //【⚡️ 传入实例】, 新建一条数据表的数据
				tx.Rollback() // 💼 回滚事务 *************
				return -1, "❌ 好友添加失败"
			}

			tx.Commit() // 💼 提交事务 *************
			return 0, "✅ 好友添加成功"
		}
		return -1, "❌ 没有找到此用户" // 否则为空, 就说明找不到这个用户
	}

	return -1, "❌ 好友 ID 不能为空" // 如果没有传入 targetId, 就返回 -1
}



// 👥 通过群来找到人的 ID
func SearchUserByGroupId(threadId uint) []uint {
	contacts := make([]ContactBasic, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", threadId).Find(&contacts)
	// 拿到了群 id 跟 群 好友的 id
	fmt.Println("🌟 查到了群的好友:", contacts, "群的 ID:", threadId)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId)) // 把好友的 ID 存储到 objIDS 切片中
	}
	return objIds
}

