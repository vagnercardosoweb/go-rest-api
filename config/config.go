package config

import (
	"os"
	"rest-api/shared"
)

var hostname, _ = os.Hostname()
var appEnv = shared.EnvGetByName("APP_ENV", "local")

var (
	Pid          = os.Getpid()
	IsDebug      = shared.EnvGetByName("DEBUG", "false") == "true"
	IsLocal      = appEnv == "local"
	IsStaging    = appEnv == "staging"
	IsProduction = appEnv == "production"
	AppEnv       = appEnv
	Hostname     = hostname
)
