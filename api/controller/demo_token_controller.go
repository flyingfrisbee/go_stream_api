package controller

import (
	"go_stream_api/api/common"
	"go_stream_api/api/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type demoTokenResponse struct {
	AuthToken string `json:"auth_token"`
}

func DemoTokenHandler(c *gin.Context) {
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
