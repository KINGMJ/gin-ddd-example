package test

import (
	"context"
	"fmt"
	"gin-ddd-example/example/ent_example/ent"
	"github.com/stretchr/testify/suite"
	"log"
)

// 公共的TestSuite，用于初始化db连接
type TestSuite struct {
	suite.Suite
	client *ent.Client
	ctx    context.Context
}

func (suite *TestSuite) SetupSuite() {
	suite.client = initClient()
	suite.ctx = context.Background()
}

func (suite *TestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Close()
	}
}

func initClient() *ent.Client {
	dsn := "root:123456@tcp(localhost:3306)/ent_test?parseTime=True"
	client, err := ent.Open("mysql", dsn,
		ent.Debug(), // 启动调试模式
		ent.Log(func(a ...any) {
			log.Printf("%s %s",
				"[SQL]",
				fmt.Sprint(a...),
			)
		}),
	)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	return client
}
