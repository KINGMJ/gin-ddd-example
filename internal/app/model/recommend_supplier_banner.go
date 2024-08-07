package model

import (
	"time"

	"gorm.io/datatypes"
)

// 持久化模型，包含模型定义和关联关系
type RecommendSupplierBannerPo struct {
	RecommendSupplierBanner
	// 定义Belongs To 关联
	RecommendSupplier RecommendSupplierPo `gorm:"foreignKey:SupplierID;references:ID" json:"recommend_supplier"`
}

type RecommendSupplierBanner struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"type:varchar(32);not null" json:"name"`
	TypeID     int64          `gorm:"not null;default:0" json:"type_id"`
	SupplierID int64          `gorm:"not null;default:0" json:"supplier_id"`
	Images     datatypes.JSON `gorm:"type:json;not null" json:"images"`
	LinkType   int8           `gorm:"type:tinyint;not null;default:1" json:"link_type"`
	LinkUrl    string         `gorm:"type:varchat(255);not null" json:"link_url"`
	Creator    int64          `gorm:"not null;default:0" json:"creator"`
	Status     int8           `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt  time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (RecommendSupplierBannerPo) TableName() string {
	return "recommend_supplier_banner"
}
