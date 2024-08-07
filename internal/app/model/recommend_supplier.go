package model

import (
	"time"

	"gorm.io/datatypes"
)

// 持久化模型，包含db模型，关联关系
type RecommendSupplierPo struct {
	RecommendSupplier
	Type RecommendSupplierTypePo `json:"type" gorm:"foreignKey:TypeID;references:ID"` // 推荐供应商类型
}

// db 模型
type RecommendSupplier struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement" json:"id"`                    // 主键ID
	Name                string         `gorm:"type:varchar(64);not null" json:"name"`                 // 名称
	Code                string         `gorm:"type:varchar(32);not null;default:''" json:"code"`      // 编号
	TypeID              int            `gorm:"not null" json:"type_id"`                               // 分类ID
	ProvinceID          int            `gorm:"not null" json:"province_id"`                           // 省ID
	CityID              int            `gorm:"not null" json:"city_id"`                               // 市ID
	DistrictID          int            `gorm:"not null" json:"district_id"`                           // 区ID
	Logo                datatypes.JSON `gorm:"type:json;not null" json:"logo"`                        // 商标或头像
	CreditCode          string         `gorm:"type:varchar(255);not null" json:"credit_code"`         // 社会信用代码
	CreateDate          string         `gorm:"type:varchar(16);not null" json:"create_date"`          // 成立日期
	LegalRepresentative string         `gorm:"type:varchar(16);not null" json:"legal_representative"` // 法人代表
	RegisteredCapital   int            `gorm:"not null" json:"registered_capital"`                    // 注册资本
	Website             string         `gorm:"type:varchar(255);not null" json:"website"`             // 网址
	Intro               string         `gorm:"type:varchar(512);not null" json:"intro"`               // 简介
	Contact             datatypes.JSON `gorm:"column:contact;type:json;not null" json:"contact"`      // 联系人信息
	Qualification       string         `gorm:"type:varchar(255);not null" json:"qualification"`       // 资质图片文件
	Sort                int            `gorm:"not null" json:"sort"`                                  // 排序
	SpecialIntro        datatypes.JSON `gorm:"type:json;not null" json:"special_intro"`               // 特殊介绍
	SpecialMark         datatypes.JSON `gorm:"type:json;not null" json:"special_mark"`                // 特殊标识
	Creator             int            `gorm:"not null" json:"creator"`                               // 创建人
	Status              int8           `gorm:"type:tinyint;not null;default:1" json:"status"`         // 状态 1.启用,0.禁用
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`                      // 创建时间
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                      // 更新时间
}

func (RecommendSupplierPo) TableName() string {
	return "recommend_supplier"
}
