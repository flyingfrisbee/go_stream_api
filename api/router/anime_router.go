package router

import (
	"go_stream_api/api/controller"
	"go_stream_api/api/token"

	"github.com/gin-gonic/gin"
)

func AnimeRouter(rg *gin.RouterGroup) {
	group := rg.Group("/anime")
	group.Use(token.JWTAuthMiddleware())
	group.GET("/recent", controller.RecentAnimeHandler)
	group.GET("/detail/:id", controller.AnimeDetailHandler)
	group.GET("/search", controller.SearchTitleHandler)
	group.POST("/video-url", controller.VideoURLHandler)
	group.POST("/detail-alt", controller.AnimeDetailAlternativeHandler)
}
