package model

import "time"

type SupplierStorePo struct {
	SupplierStore
}

type SupplierStore struct {
	ID         int64     `gorm:"primaryKey,autoIncrement" json:"id"`
	SupplierID int64     `gorm:"not null;default:0" json:"supplier_id"`
	StoreID    int64     `gorm:"not null;default:0" json:"store_id"`
	JointRate  int64     `gorm:"not null;default:0" json:"joint_rate"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (table *SupplierStorePo) TableName() string {
	return "supplier_store"
}
