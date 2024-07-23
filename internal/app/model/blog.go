package model

// 嵌套结构体，DDD中的值对象可以用这种方式建模
type Author struct {
	Name  string
	Email string
}

type Blog struct {
	BaseModel `gorm:"embedded"`
	Author    Author `gorm:"embedded;embeddedPrefix:author_" json:"author"`
	Upvotes   int64
}

func (table *Blog) TableName() string {
	return "blog"
}
