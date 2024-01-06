package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	if os.Getenv("IS_AWS_LAMBDA") == "true" {
		log.Println("Skipping load the environment, the environment is being executing with lambda")
		return
	}
	if GetAsString("LOAD_ENV", "true") == "false" {
		log.Println("Skipping load the environment")
		return
	}
	err := godotenv.Load(fmt.Sprintf(".env.%s", GetAsString("APP_ENV", "local")))
	if err != nil {
		panic(err)
	}
}

func GetAsString(name string, defaultValue ...string) string {
	value, exist := os.LookupEnv(name)
	if !exist && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func GetAsBool(name string, defaultValue ...string) bool {
	value, err := strconv.ParseBool(GetAsString(name, defaultValue...))
	if err != nil {
		panic(fmt.Sprintf(`Environment "%s" is not a boolean`, name))
	}
	return value
}

func GetAsInt(name string, defaultValue ...string) int {
	value, err := strconv.Atoi(GetAsString(name, defaultValue...))
	if err != nil {
		panic(fmt.Sprintf(`Environment "%s" is not a integer`, name))
	}
	return value
}

func GetAsFloat64(name string, defaultValue ...string) float64 {
	value, err := strconv.ParseFloat(GetAsString(name, defaultValue...), 64)
	if err != nil {
		panic(fmt.Sprintf(`Environment "%s" is not a float64`, name))
	}
	return value
}

func Required(name string) string {
	value, exist := os.LookupEnv(name)
	if !exist {
		panic(fmt.Sprintf(`Environment "%s" is required`, name))
	}
	return value
}

func GetAppEnv() string {
	return GetAsString("APP_ENV", "local")
}

func IsSchedulerEnabled() bool {
	return GetAsString("SCHEDULER_ENABLED", "true") == "true"
}

func GetSchedulerSleep() time.Duration {
	return time.Duration(GetAsInt("SCHEDULER_SLEEP", "60"))
}
