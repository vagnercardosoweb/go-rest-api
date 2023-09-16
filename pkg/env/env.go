package env

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() {
	if os.Getenv("IS_AWS_LAMBDA") == "true" {
		log.Println("Skipping load the environment, the environment is being executing with lambda")
		return
	}
	if Get("LOAD_ENV", "true") == "false" {
		log.Println("Skipping load the environment")
		return
	}
	err := godotenv.Load(fmt.Sprintf(".env.%s", Get("APP_ENV", "local")))
	if err != nil {
		panic(err)
	}
}

func Get(name string, defaultValue ...string) string {
	value, exist := os.LookupEnv(name)
	if !exist && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func GetInt(name string, defaultValue ...string) int {
	value, err := strconv.Atoi(Get(name, defaultValue...))
	if err != nil {
		panic(err)
	}
	return value
}

func Required(name string) string {
	value, exist := os.LookupEnv(name)
	if !exist {
		panic(fmt.Sprintf("Environment [%s] not exist", name))
	}
	return value
}
