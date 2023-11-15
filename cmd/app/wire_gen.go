// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gin-ddd-example/internal/app/controller"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/router"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/db"
)

// Injectors from wire.go:

// func InitApp(database *db.Database) *router.ApiRouter {
// 	entRepoImpl := repo.NewEntRepo(database)
// 	entServiceImpl := service.NewEntService(entRepoImpl)
// 	entController := controller.NewEntController(entServiceImpl)
// 	userRepoImpl := repo.NewUserRepo(database)
// 	authServiceImpl := service.NewAuthService(userRepoImpl)
// 	authController := controller.NewAuthController(authServiceImpl)
// 	apiRouter := router.NewApiRouter(entController, authController)
// 	return apiRouter
// }
