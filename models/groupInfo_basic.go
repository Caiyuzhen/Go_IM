package models
import "gorm.io/gorm"

// 群信息表
type GroupBasic struct {
	gorm.Model
	Name string // 群名称
	OwnerId uint // 群主 ID
	AvatarIcon string // 群头像
	Desc string // 群描述
	Type int // 群类型 (普通群、活动群、学习群、游戏群)  => 用 1 2 3 4 来表示
}

// ⚠️ => 类方法, 从数据库中获取表名的方法
func (table *GroupBasic) TableName() string { // TableName 为数据表, 用于指定表名
	return "group_basic" // 在 db 中的表名
}
