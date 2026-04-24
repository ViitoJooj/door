package initializer

import (
	"os"
	"strconv"

	"github.com/ViitoJooj/ward/pkg/database"
)

func Init_project() {
	database.Conn()
	EnsureMasterKey(database.DB)

	if LoadEnv(database.DB, "JWT_ACCESS_TOKEN_SECRET") == "" {
		SaveEnv(database.DB, "JWT_ACCESS_TOKEN_SECRET", randomHex(32))
	}

	if LoadEnv(database.DB, "JWT_REFRESH_TOKEN_SECRET") == "" {
		SaveEnv(database.DB, "JWT_REFRESH_TOKEN_SECRET", randomHex(32))
	}

	port := EnsureAppPort(database.DB)

	os.Setenv("JWT_ACCESS_TOKEN_SECRET", LoadEnv(database.DB, "JWT_ACCESS_TOKEN_SECRET"))
	os.Setenv("JWT_REFRESH_TOKEN_SECRET", LoadEnv(database.DB, "JWT_REFRESH_TOKEN_SECRET"))
	os.Setenv(appPortVarName, strconv.Itoa(port))
	InjectDefaultCors(database.DB)
}
