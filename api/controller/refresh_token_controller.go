package controller

import (
	"go_stream_api/api/common"
	"go_stream_api/api/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshTokenResponse struct {
	AuthToken string `json:"auth_token"`
}

// @Summary Generate auth token by providing refresh token
// @Description Refresh token will have later expiry date than auth token so user can reissue auth token whenever expired
// @Tags Token
// @Produce json
// @Success 200 {object} common.baseResponse{data=controller.refreshTokenResponse}
// @Router /token/refresh [get]
// @Param Authorization header string true "Insert your refresh token" default(Bearer <Add access token here>)
func RefreshTokenHandler(c *gin.Context) {
	authToken, err := token.GenerateJWT(token.Authorization)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, "Error while trying to refresh authorization token", http.StatusInternalServerError)
		return
	}

	response := refreshTokenResponse{
		AuthToken: authToken,
	}
	common.WrapWithBaseResponse(c, response, "Success refreshing authorization token", http.StatusOK)
}
