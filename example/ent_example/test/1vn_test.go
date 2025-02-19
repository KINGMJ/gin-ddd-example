package test

import (
	"gin-ddd-example/example/ent_example/ent/car"
	"gin-ddd-example/example/ent_example/ent/user"
	"gin-ddd-example/pkg/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"log"
	"net/url"
	"testing"
	"time"
)

type One2NTestSuit struct {
	TestSuite
}

func (s *One2NTestSuit) SetupTest() {
	s.TestSuite.SetupSuite()
}

func TestOne2NTestSuit(t *testing.T) {
	suite.Run(t, new(One2NTestSuit))
}

func (s *One2NTestSuit) TestCreateCar() {
	jerry, err := s.client.User.
		Query().
		Where(user.ID(1)).
		Only(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	// 创建关联车辆
	price, _ := decimal.NewFromString("12.23")
	bmw, err := s.client.Car.
		Create().
		SetModel("BMW X5").
		SetURL(&url.URL{Host: "https://bwm.com", Path: "/product/x5"}).
		SetRegisteredAt(time.Now()).
		SetPrice(price).
		SetStatus(car.StatusCancelled).
		SetOwner(jerry). // 关键关联操作
		Save(s.ctx)
	if err != nil {
		s.T().Fatal("创建车辆失败:", err)
	}
	utils.PrettyJson(bmw)
}

func (s *One2NTestSuit) TestCreateCars() {
	tesla, err := s.client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(s.ctx)

	if err != nil {
		s.T().Fatal(err)
	}
	log.Println("car was created: ", tesla)

	ford, err := s.client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	log.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := s.client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	log.Println("user was created: ", a8m)
}

func (s *One2NTestSuit) TestQueryCars() {
	// 查询 a8m 用户
	a8m, err := s.client.User.
		Query().
		Where(user.Name("a8m")).
		Only(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}

	// 获取该用户下的所有 cars
	cars, err := a8m.QueryCars().All(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	log.Println("cars are: ", cars)

	// 查询该用户的 cars
	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}

	utils.PrettyJson(ford)

	log.Println("cars are: ", ford)
}

// 通过 Car 定义的反向关系查询
// 存在 N+1 问题
func (s *One2NTestSuit) TestQueryCarUsers() {
	cars, err := s.client.Car.
		Query().
		All(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	for _, c := range cars {
		owner, err := c.QueryOwner().Only(s.ctx)
		if err != nil {
			s.T().Fatal(err)
		}
		log.Printf("car %q owner: %q\n", c.Model, owner.Name)
	}
}

// 通过预加载的方式一次查询
// 避免 N+1 问题
func (s *One2NTestSuit) TestQueryCarUsers2() {
	cars, err := s.client.Car.
		Query().
		WithOwner(). // 预加载owner
		All(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	for _, c := range cars {
		owner := c.Edges.Owner // 直接获取已加载的关联
		log.Printf("Car %s owner: %s", c.Model, owner.Name)
	}
}
