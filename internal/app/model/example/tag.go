package example

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null;uniqueIndex"`
	// Has Many 关系
	PostTags []PostTag `gorm:"foreignKey:TagID"`
	// 便捷访问 Posts（通过 PostTags）
	//Posts []Post `gorm:"many2many:post_tags;joinForeignKey:TagID;joinReferences:PostID"`
}
