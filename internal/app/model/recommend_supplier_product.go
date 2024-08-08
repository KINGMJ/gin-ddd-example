package model

import (
	"time"

	"gorm.io/datatypes"
)

type RecommendSupplierProductPo struct {
	RecommendSupplierProduct
}

type RecommendSupplierProduct struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code       string         `gorm:"type:varchar(32);not null" json:"code"`
	Name       string         `gorm:"type:varchar(64);not null" json:"name"`
	TypeID     int64          `gorm:"not null;default:0" json:"type_id"`
	SupplierID int64          `gorm:"not null;default:0" json:"supplier_id"`
	Params     datatypes.JSON `gorm:"type:json;not null" json:"params"`
	Images     datatypes.JSON `gorm:"type:json;not null" json:"images"`
	Intro      string         `gorm:"type:text;not null" json:"intro"`
	Creator    int64          `gorm:"not null;default:0" json:"creator"`
	Ean13      string         `gorm:"type:varchar(255);not null" json:"ean13"`
	Price      int64          `gorm:"type:int;not null;default:0" json:"price"`
	Status     int8           `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt  time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (RecommendSupplierProductPo) TableName() string {
	return "recommend_supplier_product"
}
