package router

import (
	"gin-ddd-example/internal/app/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ApiRouter struct {
	entController  *controller.EntController
	authController *controller.AuthController
}

func NewApiRouter(
	entController *controller.EntController,
	authController *controller.AuthController,
) *ApiRouter {
	return &ApiRouter{
		entController, authController,
	}
}

// 设置 route
func (ar *ApiRouter) SetupRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/signup", ar.authController.Signup)

	v1 := r.Group("/v1")
	{
		v1.GET("/ents", ar.entController.List)
		v1.POST("/ents", ar.entController.Create)
	}
}
