package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
)

func RequestId(c *gin.Context) {
	requestId := uuid.New().String()
	logger := config.GetLoggerFromCtx(c).WithID(requestId)

	c.Set(config.PgClientCtxKey, config.GetPgClientFromCtx(c).WithLogger(logger))
	c.Set(config.RequestLoggerCtxKey, logger)
	c.Set(config.RequestIdCtxKey, requestId)

	injectAwsRequestIdToHeader(c)

	c.Header("X-Request-Id", requestId)
	c.Next()
}

func injectAwsRequestIdToHeader(c *gin.Context) {
	awsRequestId := c.GetHeader("x-amzn-trace-id")

	if awsRequestId == "" {
		awsRequestId = c.GetHeader("x-amzn-requestid")
	}

	if awsRequestId != "" {
		c.Header("X-Aws-RequestId", awsRequestId)
	}
}
