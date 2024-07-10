package env

import (
	"log"
	"os"
)

func GetEnvValue(envName string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		log.Fatalf("Missing required env var: %s", envName)
	}
	return envValue
}
