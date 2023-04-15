package middlewares

import (
	"fmt"
	"net/http"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func responseError(c *gin.Context) {
	c.Next()

	reqErrors := c.Errors
	hasError := len(reqErrors) > 0
	isAborted := c.IsAborted()

	if !hasError && !isAborted {
		return
	}

	statusCode := c.Writer.Status()

	if (isAborted && statusCode == http.StatusOK) || hasError {
		statusCode = http.StatusInternalServerError
	}

	var metadata = make(logger.Metadata, 0)
	statusText := http.StatusText(statusCode)
	requestId := c.MustGet(config.RequestIdCtxKey)

	var resError = errors.New(errors.Input{
		Code:        statusText,
		StatusCode:  statusCode,
		Message:     statusText,
		SendToSlack: true,
	})

	if hasError {
		if _, ok := reqErrors[0].Err.(*errors.Input); !ok {
			metadata["errors"] = reqErrors.Errors()
			resError.Message = reqErrors[0].Error()
		} else {
			resError = reqErrors[0].Err.(*errors.Input)
		}
	}

	metadata["ip"] = c.ClientIP()
	metadata["path"] = c.Request.URL.String()

	routePath := c.FullPath()
	if routePath != "" {
		metadata["route_path"] = c.FullPath()
	}

	metadata["method"] = c.Request.Method
	metadata["query"] = c.Request.URL.Query()
	metadata["version"] = c.Request.Proto
	metadata["referrer"] = c.GetHeader("referer")
	metadata["agent"] = c.Request.UserAgent()
	metadata["body"] = c.Request.Form
	metadata["headers"] = c.Request.Header
	metadata["error"] = resError

	if config.IsDebug {
		if forwardedUser := c.GetHeader("X-Forwarded-User"); forwardedUser != "" {
			metadata["forwarded_user"] = forwardedUser
		}
		if forwardedEmail := c.GetHeader("X-Forwarded-Email"); forwardedEmail != "" {
			metadata["forwarded_email"] = forwardedEmail
		}
	}

	logger.Log(logger.Input{
		Id:       fmt.Sprintf("REQ:%s", requestId),
		Level:    logger.ERROR,
		Message:  "REQUEST_ERROR",
		Metadata: metadata,
	})

	c.JSON(resError.StatusCode, gin.H{
		"code":       resError.Code,
		"errorId":    resError.ErrorId,
		"statusCode": resError.StatusCode,
		"message":    resError.Message,
	})
}
