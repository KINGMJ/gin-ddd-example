package main

import (
	"gin-ddd-example/internal/app/controller"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/router"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/db"

	"github.com/google/wire"
)

func InitApp(database *db.Database) *router.ApiRouter {
	wire.Build(
		router.NewApiRouter,
		controller.NewEntController,
		service.NewEntService,
		repo.NewEntRepo,
	)
	return nil
}
