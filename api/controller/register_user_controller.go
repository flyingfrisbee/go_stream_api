package controller

import (
	"go_stream_api/api/common"
	"go_stream_api/api/token"
	db "go_stream_api/repository/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	UserToken string `json:"user_token" binding:"required"`
}

type registerUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RegisterUserHandler(c *gin.Context) {
	var request registerUserRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.Conn.Pg.InsertUser(request.UserToken)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	authToken, err := token.GenerateJWT(token.Authorization)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := token.GenerateJWT(token.Refresh)
	if err != nil {
		common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	response := registerUserResponse{
		AccessToken:  authToken,
		RefreshToken: refreshToken,
	}
	common.WrapWithBaseResponse(c, response, "Success register user", http.StatusOK)
}
