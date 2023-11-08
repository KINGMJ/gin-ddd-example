package route

import (
	"gin-ddd-example/internal/app/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("signup", controller.Signup)
	r.GET("login", controller.Login)

	v1 := r.Group("/v1")
	{
		v1.GET("/ents", controller.List)
		v1.POST("/ents", controller.Create)
		v1.DELETE("/ents/:id", controller.Delete)
		v1.PUT("/ents/:id", controller.Update)
	}
	return r
}
