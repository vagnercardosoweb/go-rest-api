package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
)

func RequestId(c *gin.Context) {
	requestId := uuid.New().String()
	c.Set(apicontext.RequestIdKey, requestId)

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
