package route

import (
	"gin-ddd-example/internal/app/controller"
	"gin-ddd-example/internal/app/controller/ent_controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("signup", controller.Signup)

	v1 := r.Group("/v1")
	{
		v1.GET("/ent", ent_controller.List)
		v1.POST("/ent", ent_controller.Create)
		v1.DELETE("/ent/:id", ent_controller.Delete)
		v1.PUT("/ent/:id", ent_controller.Update)
	}

	return r
}
