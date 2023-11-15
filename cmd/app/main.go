package main

import (
	"gin-ddd-example/docs"
	"gin-ddd-example/internal/app/controller"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/internal/app/router"
	"gin-ddd-example/internal/app/service"
	"gin-ddd-example/pkg/cache"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
	"gin-ddd-example/pkg/logs"
	"gin-ddd-example/pkg/rabbitmq"
	"gin-ddd-example/pkg/validate"

	"github.com/gin-gonic/gin"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// 初始化操作
	config.InitConfig()
	cache.InitRedis(*config.Conf)
	rabbitmq.InitRabbitmq(*config.Conf)
	// 日志初始化
	logs.InitLog(*config.Conf)
	logs.Log.Info("log init success!")

	database := db.InitDb()

	// swagger 配置
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	// 注册自定义的格式校验
	validate.RegisterCustomValidate()
	// 使用wire初始化应用
	apiRouter := InitApp(database)
	// 创建路由
	apiRouter.SetupRoutes(r)
	// 运行服务
	r.Run()

}

// Injectors from wire.go:

func InitApp(database *db.Database) *router.ApiRouter {
	entRepoImpl := repo.NewEntRepo(database)
	entServiceImpl := service.NewEntService(entRepoImpl)
	entController := controller.NewEntController(entServiceImpl)
	userRepoImpl := repo.NewUserRepo(database)
	authServiceImpl := service.NewAuthService(userRepoImpl)
	authController := controller.NewAuthController(authServiceImpl)
	apiRouter := router.NewApiRouter(entController, authController)
	return apiRouter
}
