package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	if os.Getenv("IS_AWS_LAMBDA") == "true" || IsLocal() {
		log.Println("skipping load the environment")
		return
	}

	envName := fmt.Sprintf(".env.%s", GetAsString("APP_ENV", "development"))
	if err := godotenv.Load(envName); err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
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
		panic(fmt.Errorf(`environment "%s" is not a boolean`, name))
	}

	return value
}

func GetAsInt(name string, defaultValue ...string) int {
	value, err := strconv.Atoi(GetAsString(name, defaultValue...))

	if err != nil {
		panic(fmt.Errorf(`environment "%s" is not a integer`, name))
	}

	return value
}

func GetAsFloat64(name string, defaultValue ...string) float64 {
	value, err := strconv.ParseFloat(GetAsString(name, defaultValue...), 64)

	if err != nil {
		panic(fmt.Errorf(`environment "%s" is not a float64`, name))
	}

	return value
}

func Required(name string) string {
	value, exist := os.LookupEnv(name)

	if !exist {
		panic(fmt.Errorf(`environment "%s" is required`, name))
	}

	return value
}

func GetAppEnv() string {
	return GetAsString("APP_ENV", "development")
}

func GetSchedulerSleep() time.Duration {
	return time.Duration(GetAsInt("SCHEDULER_SLEEP", "60"))
}

func GetRedactKeys() []string {
	keys := strings.Split(GetAsString("REDACT_KEYS", ""), ",")

	if len(keys) == 0 {
		keys = []string{"password", "passwordConfirm", "x-internal-key", "x-api-key"}
	}

	return keys
}

func IsTest() bool {
	return GetAppEnv() == Test
}

func IsDevelopment() bool {
	return GetAppEnv() == Development
}

func IsProduction() bool {
	return GetAppEnv() == Production
}

func IsLocal() bool {
	return GetAsBool("IS_LOCAL", "false")
}

func IsSchedulerEnabled() bool {
	return GetAsBool("SCHEDULER_ENABLED", "false")
}

func IsAlertOnServerStart() bool {
	return GetAsBool("SLACK_ALERT_ON_SERVER_START", "false")
}

func IsAlertOnServerClose() bool {
	return GetAsBool("SLACK_ALERT_ON_SERVER_CLOSE", "false")
}
