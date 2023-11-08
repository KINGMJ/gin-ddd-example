package db

import (
	"fmt"
	"gin-ddd-example/pkg/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() {
	c := config.Conf.DBConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		c.Host,
		c.User,
		c.Password,
		c.Dbname,
		c.Port,
		c.TimeZone,
	)
	databases, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	Db = databases
}
