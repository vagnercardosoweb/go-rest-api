package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	loggerpkg "github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
)

func RequestId(c *gin.Context) {
	requestId := uuid.New().String()
	logger := utils.GetLogger(c).WithId(requestId)

	c.Set(postgres.CtxKey, utils.GetPgClient(c).WithLogger(logger))
	c.Set(loggerpkg.CtxKey, logger)
	c.Set(utils.RequestIdCtxKey, requestId)

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
