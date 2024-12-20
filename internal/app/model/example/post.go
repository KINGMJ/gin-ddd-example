package example

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null"`
	Content string `gorm:"type:text"`
	UserID  uint   `gorm:"not null;index"` // 外键字段
	// Has Many 关系
	PostTags []PostTag `gorm:"foreignKey:PostID"`
	// 便捷访问 Tags（通过 PostTags）
	//Tags []Tag `gorm:"many2many:post_tags;joinForeignKey:PostID;References:TagID"`
}
