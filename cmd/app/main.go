package main

import (
	"gin-ddd-example/docs"
	"gin-ddd-example/internal/app/route"
	"gin-ddd-example/pkg/config"
	"gin-ddd-example/pkg/db"
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
	// // 初始化操作
	config.InitConfig()
	db.InitDb()

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := route.InitRouter()

	// 运行服务
	r.Run()

	// db.Db.AutoMigrate(&model.User{})
}
