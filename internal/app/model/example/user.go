package example

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(50);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex"`
	CompanyID uint   // 用户 belong to 公司

	// 定义关系
	Company Company `gorm:"foreignKey:CompanyID"` // Belongs To 关系
	Profile Profile `gorm:"foreignKey:UserID"`    // Has One 关系
	Posts   []Post  `gorm:"foreignKey:UserID"`    // Has Many 关系
}
