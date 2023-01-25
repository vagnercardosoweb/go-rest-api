package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvLoadFromFile(filenames ...string) {
	if os.Getenv("IS_AWS_LAMBDA") == "true" {
		return
	}
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Println("shared.EnvLoadFromFile")
		panic(err)
	}
}

func EnvGetByName(name string, defaultValue ...string) string {
	value := os.Getenv(name)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func EnvRequiredByName(name string, defaultValue ...string) string {
	value := EnvGetByName(name, defaultValue...)
	if len(value) == 0 {
		panic(fmt.Sprintf("Environment [%s] not exist", name))
	}
	return value
}
