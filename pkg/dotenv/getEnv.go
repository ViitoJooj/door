package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

var JwtAccessTokenSecret string
var JwtRefreshTokenSecret string

func GetEnv() {
	godotenv.Load(".env")
	godotenv.Load("../.env")
	godotenv.Load("../../.env")
	godotenv.Load("../../../.env")

	JwtAccessTokenSecret = os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	if JwtAccessTokenSecret == "" {
		panic("JwtAccessTokenSecret is null")
	}

	JwtRefreshTokenSecret = os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	if JwtRefreshTokenSecret == "" {
		panic("JwtRefreshTokenSecret is null")
	}
}
