package router

import (
	"go_stream_api/api/controller"
	"go_stream_api/api/token"

	"github.com/gin-gonic/gin"
)

func RefreshTokenRouter(rg *gin.RouterGroup) {
	group := rg.Group("/token")
	group.Use(token.JWTRefreshMiddleware())
	group.GET("/refresh", controller.RefreshTokenHandler)
}
