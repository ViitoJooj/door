package security

import (
	"regexp"
	"strings"
)

var (
	sqlCommandRegex = regexp.MustCompile(`(?i)\b(select|insert|update|delete|drop|alter|create|truncate|union|grant|revoke)\b[\s\S]{0,80}\b(from|into|table|set|where|join|values)\b`)
	sqlExecRegex    = regexp.MustCompile(`(?i)\b(exec|execute)\b`)
)

func ContainsSQLCommand(input string) bool {
	value := strings.TrimSpace(input)
	if value == "" {
		return false
	}

	if sqlCommandRegex.MatchString(value) {
		return true
	}

	if sqlExecRegex.MatchString(value) && strings.Contains(strings.ToLower(value), "sql") {
		return true
	}

	return false
}
