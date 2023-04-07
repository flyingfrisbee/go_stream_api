package token

import (
	"fmt"
	"go_stream_api/api/common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := parseTokenFromRequestHeader(c)
		if err != nil {
			common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
			c.Abort()
			return
		}

		err = validateJWT(Authorization, token)
		if err != nil {
			common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

func JWTRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := parseTokenFromRequestHeader(c)
		if err != nil {
			common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusBadRequest)
			c.Abort()
			return
		}

		err = validateJWT(Refresh, token)
		if err != nil {
			common.WrapWithBaseResponse(c, nil, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}

func parseTokenFromRequestHeader(c *gin.Context) (string, error) {
	reqHeader := c.Request.Header.Get("Authorization")
	reqHeaderSlice := strings.Split(reqHeader, " ")
	if len(reqHeaderSlice) != 2 {
		return "", fmt.Errorf("wrong authorization format, should be like: 'Bearer {insert_your_token_here}'")
	}

	token := reqHeaderSlice[1]
	return token, nil
}
