package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoverWhenPanic(c *gin.Context, errMsg string) {
	if err := recover(); err != nil {
		log.Println(err)
		WrapWithBaseResponse(c, nil, errMsg, http.StatusInternalServerError)
	}
}
