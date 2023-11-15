//go:build wireinject
// +build wireinject

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
		controller.NewAuthController,
		service.NewEntService,
		service.NewAuthService,
		repo.NewEntRepo,
		repo.NewUserRepo,
		// 接口与实现类绑定
		wire.Bind(new(service.EntService), new(*service.EntServiceImpl)),
		wire.Bind(new(repo.EntRepo), new(*repo.EntRepoImpl)),
		wire.Bind(new(service.AuthService), new(*service.AuthServiceImpl)),
		wire.Bind(new(repo.UserRepo), new(*repo.UserRepoImpl)),
	)
	return nil
}
