package example

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Address string `gorm:"type:varchar(200)"`
}
