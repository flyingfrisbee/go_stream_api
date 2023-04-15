package api

import (
	"context"
	"fmt"
	"go_stream_api/api/docs"
	"go_stream_api/api/router"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	srv *http.Server
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

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	go runServer(srv)
}

func Stop() {
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func registerRouters(rg *gin.RouterGroup) {
	router.DemoTokenRouter(rg)
	router.RefreshTokenRouter(rg)
	router.AnimeRouter(rg)
	router.UserRouter(rg)
}

func runServer(srv *http.Server) {
	defer log.Println("Stopping API...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
