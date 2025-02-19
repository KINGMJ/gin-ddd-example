package test

import (
	_ "gin-ddd-example/example/ent_example/ent/runtime"
	"gin-ddd-example/example/ent_example/ent/schema"
	"gin-ddd-example/example/ent_example/ent/user"
	"gin-ddd-example/pkg/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type BaseTestSuit struct {
	TestSuite
}

func (s *BaseTestSuit) SetupTest() {
	s.TestSuite.SetupSuite()
}

func TestBaseTestSuit(t *testing.T) {
	suite.Run(t, new(BaseTestSuit))
}

func (s *BaseTestSuit) TestCreateUser() {
	u, err := s.client.User.Create().
		SetAge(30).
		SetName("a8m").
		SetActive(true).
		Save(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	log.Printf("user was created: %+v", u)
}

func (s *BaseTestSuit) TestQueryUser() {
	u, err := s.client.User.Query().
		Where(user.Name("a8m")).
		// 必须返回 1 条，否则错误
		Only(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	log.Printf("user was found: %+v", u)
}

// 测试软删除
func (s *BaseTestSuit) TestSoftDeleteUser() {
	err := s.client.User.DeleteOneID(1).Exec(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
}

// 查询，包含软删除
func (s *BaseTestSuit) TestFindSoftDeleteUser() {
	users, err := s.client.User.Query().All(schema.SkipSoftDelete(s.ctx))
	if err != nil {
		s.T().Fatal(err)
	}
	utils.PrettyJson(users)
}

// 物理删除，也是通过 schema.SkipSoftDelete(s.ctx) 来实现的
func (s *BaseTestSuit) TestHardDeleteUser() {
	err := s.client.User.DeleteOneID(1).Exec(schema.SkipSoftDelete(s.ctx))
	if err != nil {
		s.T().Fatal(err)
	}
}
