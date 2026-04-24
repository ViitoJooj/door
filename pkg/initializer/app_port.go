package initializer

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	appPortVarName  = "APP_PORT"
	defaultAppPort  = 7171
	maxPortDistance = 65535
)

func IsAppPortVar(name string) bool {
	return strings.EqualFold(strings.TrimSpace(name), appPortVarName)
}

func ParseAppPort(value string) (int, error) {
	port, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("APP_PORT must be a valid number")
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("APP_PORT must be between 1 and 65535")
	}
	return port, nil
}

func IsPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	_ = listener.Close()
	return true
}

func FindNearestAvailablePort(preferred int) (int, error) {
	if preferred < 1 || preferred > 65535 {
		return 0, fmt.Errorf("preferred port must be between 1 and 65535")
	}

	if IsPortAvailable(preferred) {
		return preferred, nil
	}

	for distance := 1; distance <= maxPortDistance; distance++ {
		high := preferred + distance
		if high <= 65535 && IsPortAvailable(high) {
			return high, nil
		}

		low := preferred - distance
		if low >= 1 && IsPortAvailable(low) {
			return low, nil
		}
	}

	return 0, errors.New("no available port found")
}

func EnsureAppPort(db *sql.DB) int {
	configuredValue := LoadEnv(db, appPortVarName)
	if configuredValue != "" {
		configuredPort, err := ParseAppPort(configuredValue)
		if err == nil && IsPortAvailable(configuredPort) {
			return configuredPort
		}
	}

	port, err := FindNearestAvailablePort(defaultAppPort)
	if err != nil {
		panic(err)
	}

	SaveEnv(db, appPortVarName, strconv.Itoa(port))
	return port
}
