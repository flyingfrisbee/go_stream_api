package common

import (
	"github.com/gin-gonic/gin"
)

type baseResponse struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

func WrapWithBaseResponse(
	c *gin.Context,
	data interface{},
	message string,
	statusCode int,
) {
	c.JSON(statusCode, baseResponse{
		Data:       data,
		Message:    message,
		StatusCode: statusCode,
	})
}
