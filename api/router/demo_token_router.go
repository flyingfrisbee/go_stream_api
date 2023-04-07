package router

import (
	"go_stream_api/api/controller"
	env "go_stream_api/environment"

	"github.com/gin-gonic/gin"
)

func DemoTokenRouter(rg *gin.RouterGroup) {
	rg.GET(env.RouterSecretPath, controller.DemoTokenHandler)
}
