package test

import (
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

// 公共的TestSuite，用于初始化db连接
type RepoTestSuite struct {
	suite.Suite
	db  *db.Database
	db1 *gorm.DB
}

func (suite *RepoTestSuite) SetupTest() {
	config.InitConfig()
	// 日志初始化
	logs.InitLog(*config.Conf)
	suite.db = db.InitDb()
	suite.db1 = suite.db.DB
}
