package test

import (
	"gin-ddd-example/example/ent_example/ent/group"
	"gin-ddd-example/pkg/utils"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

// 多对多模型测试
type N2NTestSuit struct {
	TestSuite
}

func (s *N2NTestSuit) SetupTest() {
	s.TestSuite.SetupSuite()
}

func TestN2NTestSuit(t *testing.T) {
	suite.Run(t, new(N2NTestSuit))
}

func (s *N2NTestSuit) TestCreateGraph() {
	// create users
	ariel, err := s.client.User.
		Create().
		SetAge(30).
		SetName("Ariel").
		Save(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}

	neta, err := s.client.User.
		Create().
		SetAge(30).
		SetName("neta").
		Save(s.ctx)

	if err != nil {
		s.T().Fatal(err)
	}

	// create cars
	err = s.client.Car.
		Create().
		SetModel("tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(ariel).
		Exec(s.ctx)

	if err != nil {
		s.T().Fatal(err)
	}
	err = s.client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		// Attach this car to Ariel.
		SetOwner(ariel).
		Exec(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		// Attach this car to Neta.
		SetOwner(neta).
		Exec(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}

	// create groups
	err = s.client.Group.
		Create().
		SetName("Gitlab").
		AddUsers(neta, ariel).
		Exec(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	log.Println("The graph was created successfully")
}

// 获取Gitlab 组下面所有用户的车辆
func (s *N2NTestSuit) TestQueryGitlab() {
	cars, err := s.client.Group.
		Query().
		Where(group.Name("Gitlab")).
		QueryUsers().
		QueryCars().
		All(s.ctx)

	if err != nil {
		s.T().Fatal(err)
	}
	utils.PrettyJson(cars)
}
