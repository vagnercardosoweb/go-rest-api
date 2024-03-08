package tests

import (
	"context"
	"os"

	"github.com/stretchr/testify/suite"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

var environments = map[string]string{
	"APP_ENV":        "test",
	"LOGGER_ENABLED": "false",
	"LOGGER_DEBUG":   "false",
	"DEBUG":          "false",
	"TZ":             "UTC",
}

type GlobalTestSuite struct {
	suite.Suite
	Logger *logger.Logger
	Ctx    context.Context
}

func (s *GlobalTestSuite) SetupSuite() {
	for k, v := range environments {
		_ = os.Setenv(k, v)
	}

	s.Ctx = context.Background()
	s.Logger = logger.New()
}
