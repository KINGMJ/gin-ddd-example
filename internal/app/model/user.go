package model

import "gorm.io/gorm"

// 用户模型
type User struct {
	gorm.Model
	Name     string
	Email    string
	Phone    string
	Password string
	Sex      uint8
}

func (table *User) TableName() string {
	return "user"
}
