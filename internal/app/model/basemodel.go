package model

import (
	"gin-ddd-example/internal/app/model/ctype"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt ctype.NullTime `json:"created_at"`
	UpdatedAt ctype.NullTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
