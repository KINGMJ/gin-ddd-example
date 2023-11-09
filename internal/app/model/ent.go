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

// 创建企业，请求参数

// 企业列表加载 dto

type UpdateEntReq struct {
	EntName      string `form:"ent_name"`
	EntDesc      string `form:"ent_desc"`
	ContactName  string `form:"contact_name"`
	ContactEmail string `form:"contact_email"`
	ContactPhone string `form:"contact_phone"`
}

// func (req AddEntReq) ToEnt() Ent {
// 	return Ent{
// 		Name: req.EntName,
// 		Desc: req.EntDesc,
// 	}
// }
