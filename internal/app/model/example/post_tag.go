package example

import "gorm.io/gorm"

// PostTag 中间表模型
type PostTag struct {
	gorm.Model
	PostID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
	Sort   int  `gorm:"not null;default:0"`
	// Belongs To 关系
	Post Post `gorm:"foreignKey:PostID"`
	Tag  Tag  `gorm:"foreignKey:TagID"`
}

//google.type.decimal
//google.type.date
//google.type.timestamp
