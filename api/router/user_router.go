package router

import (
	"go_stream_api/api/controller"
	"go_stream_api/api/token"
	env "go_stream_api/environment"

	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup) {
	rg.POST(env.RouterSecretPath2, controller.RegisterUserHandler)
	group := rg.Group("/bookmark")
	group.Use(token.JWTAuthMiddleware())
	group.POST("/", controller.SaveBookmarkHandler)
	group.DELETE("/", controller.DeleteBookmarkHandler)
}
