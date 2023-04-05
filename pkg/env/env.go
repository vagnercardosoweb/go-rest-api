package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadFromLocal() {
	if os.Getenv("IS_AWS_LAMBDA") == "true" {
		log.Println("Skipping load the environment, the environment is being executing with lambda")
		return
	}
	if Get("APP_ENV", "local") != "local" {
		log.Println("Only APP_ENV=local will be loaded the Environments")
		return
	}
	err := godotenv.Load(".env.local")
	if err != nil {
		panic(err)
	}
}

func Get(name string, defaultValue ...string) string {
	value := os.Getenv(name)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func Required(name string, defaultValue ...string) string {
	value := Get(name, defaultValue...)
	if len(value) == 0 {
		panic(fmt.Sprintf("Environment [%s] not exist", name))
	}
	return value
}
