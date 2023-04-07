package token

import (
	"fmt"
	env "go_stream_api/environment"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType int

const (
	Authorization TokenType = iota
	Refresh
)

func (tt TokenType) getString() string {
	switch tt {
	case Authorization:
		return "authorization"
	default:
		return "refresh"
	}
}

// Returns authToken / refreshToken depending on tokenType, error if there's any
func GenerateJWT(expectedTokenType TokenType) (string, error) {
	var token *jwt.Token

	switch expectedTokenType {
	case Authorization:
		expiryDays, _ := strconv.Atoi(env.AuthTokenValidityDays)
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"token_type": Authorization,
			"exp":        time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour).Unix(),
		})
	case Refresh:
		expiryDays, _ := strconv.Atoi(env.RefreshTokenValidityDays)
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"token_type": Refresh,
			"exp":        time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour).Unix(),
		})
	}

	tokenString, err := token.SignedString([]byte(env.APISecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(expectedTokenType TokenType, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.APISecretKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token")
	}

	tokenType := claims["token_type"]
	tokenTypeAsFloat, ok := tokenType.(float64)
	if !ok {
		log.Fatal(fmt.Errorf(`expecting "token_type" from JWT as float64 type but got %T`, tokenType))
	}

	actualTokenType := TokenType(tokenTypeAsFloat)
	if actualTokenType != expectedTokenType {
		return fmt.Errorf(
			"token type mismatch, actualType: %s, expectedType: %s",
			actualTokenType.getString(),
			expectedTokenType.getString(),
		)
	}

	return nil
}
