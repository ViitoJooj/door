package dotenv

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var JwtAccessTokenSecret string
var JwtRefreshTokenSecret string
var CorsOriginsMap map[string]struct{}
var IP2LocationBin string

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

	CorsOriginsMap = make(map[string]struct{})
	origins := os.Getenv("ALLOWED_ORIGINS")
	for _, o := range strings.Split(origins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			CorsOriginsMap[o] = struct{}{}
		}
	}

	IP2LocationBin = os.Getenv("IP2LOCATION_BIN")
}
