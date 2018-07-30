package properties

import (
	"log"
	"os"
)

func GetTimezone() string {
	result := os.Getenv("TIMEZONE")
	log.Println("HIER", result == "")
	if result == "" {
		return "Europe/Berlin"
	}
	return result
}

func GetMongoDbPath() string {
	return getRequiredEnv("mongoDbPath")
}

func getRequiredEnv(env string) string {
	result := os.Getenv(env)
	if len(result) == 0 {
		panic("NEED ENV " + env)
	}
	return result
}
