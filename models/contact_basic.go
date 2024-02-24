package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

// 人员关系表
type ContactBasic struct { // (contact_basic 这个 model 表示的是一组关系, owner_id 表示这个好友是谁的, target_id 好友是谁, 比如 24 是 3 owner 的好友) 
	gorm.Model
	OwnerId uint // 关系拥有者 ID, 类型 uint 要与 MessageBasic 中的 FromId、ToId 一致
	TargetId uint // 关系的目标 ID, 类型 uint 要与 MessageBasic 中的 FromId、ToId 一致
	Type int // 关系类型 (好友、关注、粉丝、黑名单)  => 用 1 2 3 来表示 (1: 好友, 2: 群组, 3: 黑名单), 后面可以扩展 4: 粉丝等等预留字段
	Desc string // 描述信息(备注、标签), 预留字段
}



// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *ContactBasic) TableName() string { // TableName 为数据表, 用于指定表名
	return "contact_basic" // 在 db 中的表名
}




// 查找【某个人】的好友 (contact_basic 这个 model 表示的是一组关系, owner_id 表示这个好友是谁的, target_id 好友是谁, 比如 24 是 3 owner 的好友) 
// 也可以直接 SQL 查询  =>  SELECT * FROM `user_basic` WHERE id in (20,21) AND `user_basic`.`deleted_at` IS NULL
func SearchFriend(userId uint) []UserBasic { // 传入 userID, 返回好友的具体信息
	contacts := make([]ContactBasic, 0) // ContactBasic 类型的切片, 用来储存一组好友
	objIDS := make([]uint64, 0) // uint 类型的切片, 用来存储好友的 ID, 然后再在下面的 for 循环中去取出好友

	// 通过数据库去查找这个人的好友, 过滤 contact
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts) // owner_id 表示某个人的好友, type = 1 写死是好友关系的类型

	// 取出好友
	for _, v := range contacts {
		fmt.Println("✅ 查到了好友 ID:", v.TargetId) // 打印出好友的 ID
		objIDS = append(objIDS, uint64(v.TargetId)) // 把好友的 ID 存储到 objIDS 切片中
	}

	users := make([]UserBasic, 0) // 用来存储好友的信息
	utils.DB.Where("id in ?", objIDS).Find(&users) // 通过好友的 ID 去查找好友的信息, 需要使用 In 查询, 取一定的范围
	fmt.Println("✅ 查到了好友的信息:", users)
	
	return users // 返回好友的信息
}