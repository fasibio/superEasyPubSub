package properties

import (
	"os"
)

func GetTimezone() string {
	result := os.Getenv("TIMEZONE")
	if result == "" {
		return "Europe/Berlin"
	}
	return result
}

func getRequiredEnv(env string) string {
	result := os.Getenv(env)
	if len(result) == 0 {
		panic("NEED ENV " + env)
	}
	return result
}
