package config

import (
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

var IsDebug = env.Get("DEBUG", "false") == "true"
var IsLocal = env.Get("APP_ENV", "local") == "local"
var IsProduction = env.Get("APP_ENV", "local") == "production"
var AppEnv = env.Get("APP_ENV", "local")

// var IsStaging = env.Get("APP_ENV", "local") == "staging"

// var LocationGlobal, _ = time.LoadLocation(env.Get("TZ", "UTC"))
var LocationBrl, _ = time.LoadLocation("America/Sao_Paulo")

var JwtExpiresIn = time.Duration(env.GetInt("JWT_EXPIRES_IN_SECONDS", "86400"))
