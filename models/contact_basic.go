package models
import "gorm.io/gorm"

// 人员关系表
type ContactBasic struct {
	gorm.Model
	OwnerId uint // 关系拥有者 ID, 类型 uint 要与 MessageBasic 中的 FromId、ToId 一致
	TargetId uint // 关系的目标 ID, 类型 uint 要与 MessageBasic 中的 FromId、ToId 一致
	Type int // 关系类型 (好友、关注、粉丝、黑名单)  => 用 0 1 3 来表示
	Desc string // 描述信息(备注、标签), 预留字段
}

// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *ContactBasic) TableName() string { // TableName 为数据表, 用于指定表名
	return "contact_basic" // 在 db 中的表名
}
	