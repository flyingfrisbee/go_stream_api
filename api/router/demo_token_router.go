package router

import (
	"go_stream_api/api/common"
	"go_stream_api/api/token"
	env "go_stream_api/environment"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DemoTokenRouter(rg *gin.RouterGroup) {
	rg.GET(env.RouterSecretPath, demoTokenHandler)
}

type demoTokenResponse struct {
	AuthToken string `json:"auth_token"`
}

// @Summary Generate auth token
// @Description This endpoint is intended for testing purpose only
// @Tags Token
// @Produce json
// @Success 200 {object} common.baseResponse{data=router.demoTokenResponse}
// @Router /huh/wrong/endpoint [get]
func demoTokenHandler(c *gin.Context) {
	authToken, err := token.GenerateJWT(token.Authorization)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, "Error while trying to generate authorization token", http.StatusInternalServerError)
		return
	}

	response := demoTokenResponse{
		AuthToken: authToken,
	}
	common.WrapWithBaseResponse(c, response, "Success generating authorization token", http.StatusOK)
}
