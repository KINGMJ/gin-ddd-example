package model

import "time"

type RecommendSupplierTypePo struct {
	RecommendSupplierType
}

type RecommendSupplierType struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(32);not null" json:"name"`
	Status    int8      `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (RecommendSupplierTypePo) TableName() string {
	return "recommend_supplier_type"
}
