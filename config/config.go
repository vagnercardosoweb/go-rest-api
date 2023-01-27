package config

import (
	"os"
	"rest-api/shared"
)

var hostname, _ = os.Hostname()
var appEnv = shared.EnvGetByName("APP_ENV", "local")

var (
	Pid      = os.Getpid()
	AppEnv   = appEnv
	Hostname = hostname

	IsLocal      = appEnv == "local"
	IsStaging    = appEnv == "staging"
	IsProduction = appEnv == "production"
	IsDebug      = shared.EnvGetByName("DEBUG", "false") == "true"

	LoggerContextKey    = "logger"
	RequestIdContextKey = "requestId"
)
