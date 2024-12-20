package example

import (
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gin-ddd-example/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ExampleTest struct {
	suite.Suite
	db *gorm.DB
}

func (s *ExampleTest) SetupTest() {
	config.InitConfig()
	// 日志初始化
	logs.InitLog(*config.Conf)
	s.db = db.InitDb().DB
	// 执行迁移语句
	// 1. 先执行数据库迁移
	//s.db.DisableForeignKeyConstraintWhenMigrating = true
	err := s.db.AutoMigrate(&Company{}, &User{}, &Profile{}, &Post{}, &Tag{}, &PostTag{})
	s.NoError(err)
}

// 清理测试数据
func (s *ExampleTest) TearDownTest() {
	//s.db.Migrator().DropTable(&User{}, &Company{}, &Profile{}, &Post{}, &Tag{}, &PostTag{})
}

func TestExampleTest(t *testing.T) {
	suite.Run(t, new(ExampleTest))
}

func (s *ExampleTest) createTestData() {
	company1 := Company{
		Name:    "腾讯",
		Address: "深圳",
	}
	company2 := Company{
		Name:    "阿里",
		Address: "杭州",
	}

	s.db.Create(&company1)
	s.db.Create(&company2)

	users := []User{
		{
			Name:      "张三",
			Email:     "zhangsan@qq.com",
			CompanyID: company1.ID,
			Profile: Profile{
				Age:     16,
				Address: "广东省深圳市",
			},
		},
		{
			Name:      "李四",
			Email:     "lisi@alibaba.com",
			CompanyID: company2.ID,
			Profile: Profile{
				Age:     21,
				Address: "上海市浦东新区",
			},
		},
	}

	for _, user := range users {
		s.db.Create(&user)
	}
}

func (s *ExampleTest) TestPreload() {
	s.createTestData()
	var user User
	err := s.db.Preload("Company").First(&user, 1).Error
	s.NoError(err)
	utils.PrettyJson(user)
}

func (s *ExampleTest) TestJoinQuery() {
	//s.createTestData()
	var user User
	err := s.db.Joins("Company", s.db.Where(&Company{Name: "腾讯"})).
		First(&user, "users.company_id = ?", 2).
		Order("companies.name ASC").Error

	s.NoError(err)
	utils.PrettyJson(user)
}

func (s *ExampleTest) TestHasManyCreate() {
	user := User{
		Name:      "王舞",
		CompanyID: 1,
		Profile: Profile{
			Age:     21,
			Address: "上海市徐汇区",
		},
		Posts: []Post{
			{
				Title:   "文章1",
				Content: "内容1",
			},
			{
				Title:   "文章2",
				Content: "内容2",
			},
			{
				Title:   "文章3",
				Content: "内容3",
			},
		},
	}

	err := s.db.Create(&user).Error
	s.NoError(err)
	s.Equal(3, len(user.Posts))
}

func (s *ExampleTest) TestHasManyAppend() {
	// 已存在的用户
	var user User
	err := s.db.First(&user, 3).Error
	s.NoError(err)

	posts := []Post{
		{
			Title:   "文章3",
			Content: "内容3",
			// 不需要手动设置 UserID
		},
		{
			Title:   "文章4",
			Content: "内容4",
		},
	}
	err = s.db.Model(&user).Association("Posts").Append(posts)
	s.NoError(err)
}

func (s *ExampleTest) TestMany2ManyCreate() {
	user := User{
		Name:      "小刘子",
		CompanyID: 1,
		Profile: Profile{
			Age:     30,
			Address: "上海市闵行区",
		},
		Posts: []Post{
			{
				Title:   "文章A",
				Content: "内容A",
				//Tags: []Tag{
				//	{Name: "经济"},
				//	{Name: "哲学"},
				//},
			},
			{
				Title:   "文章B",
				Content: "内容B",
				//Tags: []Tag{
				//	{Name: "历史"},
				//	{Name: "艺术"},
				//},
			},
		},
	}
	err := s.db.Create(&user).Error
	s.NoError(err)
}

func (s *ExampleTest) TestMany2Many2Create() {
	post := Post{
		Title:   "文章4",
		Content: "内容4",
		UserID:  1,
		PostTags: []PostTag{
			{
				Tag:  Tag{Name: "标签1"},
				Sort: 1,
			},
			{
				Tag:  Tag{Name: "标签2"},
				Sort: 2,
			},
			{
				Tag:  Tag{Name: "标签3"},
				Sort: 3,
			},
		},
	}

	err := s.db.Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Create(&post).Error
	s.NoError(err)
}

func (s *ExampleTest) TestMany2Many2Preload() {
	// 查询文章以及标签
	var post Post
	err := s.db.Preload("PostTags", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC")
	}).Preload("PostTags.Tag").Find(&post, 3).Error
	s.NoError(err)
	utils.PrettyJson(post)
}
