package api

import (
	"go_stream_api/api/docs"
	"go_stream_api/api/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Run 'swag init -o ./api/docs' everytime changes that affect swagger occurs
func Run() {
	// Disable debug log
	gin.SetMode(gin.ReleaseMode)
	docs.SwaggerInfo.BasePath = "/api/v1"

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiGroup := r.Group("/api/v1")
	registerRouters(apiGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}

}

func registerRouters(rg *gin.RouterGroup) {
	router.DemoTokenRouter(rg)
	router.RefreshTokenRouter(rg)
}
