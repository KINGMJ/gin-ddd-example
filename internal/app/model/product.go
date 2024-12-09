package model

import (
	"time"
)

type Product struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Unit        string    `gorm:"column:unit;type:varchar(100);not null;default:''" json:"unit"`
	EAN13       string    `gorm:"column:ean13;type:varchar(255);not null" json:"ean13"`
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Images      string    `gorm:"column:images;type:json" json:"images"`
	MerchantID  int64     `gorm:"column:merchant_id;not null" json:"merchant_id"`
	StoreID     int64     `gorm:"column:store_id;not null;default:0" json:"store_id"`
	Type        int64     `gorm:"column:type;not null;default:0" json:"type"`
	Brand       string    `gorm:"column:brand;type:varchar(255)" json:"brand"`
	BuyPrice    int64     `gorm:"column:buy_price;not null;default:0" json:"buy_price"`
	Price       int64     `gorm:"column:price;not null;default:0" json:"price"`
	SettlePrice int64     `gorm:"column:settle_price;not null;default:0" json:"settle_price"`
	Discount    int8      `gorm:"column:discount;not null;default:100" json:"discount"`
	Created     time.Time `gorm:"column:created;not null;default:CURRENT_TIMESTAMP" json:"created"`
	Updated     time.Time `gorm:"column:updated;not null;default:CURRENT_TIMESTAMP" json:"updated"`
	Specs       string    `gorm:"column:specs;type:varchar(255);not null" json:"specs"`
	ExpireTime  string    `gorm:"column:expire_time;type:varchar(100);not null;default:''" json:"expire_time"`
	Measurement int8      `gorm:"column:measurement;not null;default:1" json:"measurement"`
	Source      int8      `gorm:"column:source;not null;default:1" json:"source"`
	Status      int8      `gorm:"column:status;not null;default:3" json:"status"`
}

func (t *Product) TableName() string {
	return "product"
}
