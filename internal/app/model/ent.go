package model

import "gorm.io/gorm"

// 企业持久化模型
type Ent struct {
	gorm.Model
	Name              string `gorm:"unique;not null"`
	Desc              string
	PmoEnable         uint8 `gorm:"default:1"`
	ProviderScope     uint8 `gorm:"default:1"`
	AssignedProviders string
	Edition           uint8 `gorm:"default:1"`
}
