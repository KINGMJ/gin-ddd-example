package test_suite

import (
	"context"
	"gin-ddd-example/pkg/container"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	container.Container
	Ctx context.Context
}

func (s *TestSuite) SetupSuite() {
	s.Ctx = context.Background()
	s.Container = container.NewContainer()
}
