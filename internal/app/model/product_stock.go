package model

import "time"

type ProductStockPo struct {
	ProductStock
}

type ProductStock struct {
	ID        int64     `gorm:"primaryKey,autoIncrement" json:"id"`
	StoreID   int64     `gorm:"not null;default:0;comment:商户|仓库id" json:"store_id"`
	ProductID int64     `gorm:"not null;default:0;uniqueIndex;comment:商品id" json:"product_id"`
	Count     int64     `gorm:"not null;default:0;comment:商品数量" json:"count"`
	Created   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created"`
	Updated   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间" json:"updated"`
}

func (table *ProductStockPo) TableName() string {
	return "product_stock"
}

// ProductStocks 自定义类型
type ProductStocks []*ProductStock
