package controller

import (
	"go_stream_api/api/common"
	db "go_stream_api/repository/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	UserToken string `json:"user_token" binding:"required"`
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

	common.WrapWithBaseResponse(c, nil, "Success register user", http.StatusOK)
}
