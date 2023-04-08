package router

import (
	"go_stream_api/api/controller"
	env "go_stream_api/environment"

	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup) {
	rg.POST(env.RouterSecretPath2, controller.RegisterUserHandler)
}
