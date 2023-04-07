package environment

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	APISecretKey             string
	AuthTokenValidityDays    string
	RefreshTokenValidityDays string
	RouterSecretPath         string
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	APISecretKey = os.Getenv("API_SECRET_KEY")
	AuthTokenValidityDays = os.Getenv("AUTH_TOKEN_VALIDITY_DAYS")
	RefreshTokenValidityDays = os.Getenv("REFRESH_TOKEN_VALIDITY_DAYS")
	RouterSecretPath = os.Getenv("ROUTER_SECRET_PATH")

}
