package model

import "gorm.io/gorm"

// 企业模型
type Ent struct {
	gorm.Model
	Name              string `gorm:"unique;not null"`
	Desc              string
	PmoEnable         uint8 `gorm:"default:1"`
	ProviderScope     uint8 `gorm:"default:1"`
	AssignedProviders string
	Edition           uint8 `gorm:"default:1"`
}

// 创建企业，请求参数
type AddEntReq struct {
	EntName      string `form:"ent_name" binding:"required"`
	EntDesc      string `form:"ent_desc"`
	ContactName  string `form:"contact_name" binding:"required"`
	ContactEmail string `form:"contact_email" binding:"required"`
	ContactPhone string `form:"contact_phone" binding:"required"`
}

func (req AddEntReq) ToEnt() Ent {
	return Ent{
		Name: req.EntName,
		Desc: req.EntDesc,
	}
}
