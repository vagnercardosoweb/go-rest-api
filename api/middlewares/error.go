package middlewares

import (
	"net/http"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context) {
	c.Next()

	requestErrors := c.Errors
	hasError := len(requestErrors) > 0
	isAborted := c.IsAborted()

	if !hasError && !isAborted {
		return
	}

	statusCode := c.Writer.Status()

	if (isAborted && statusCode == http.StatusOK) || hasError {
		statusCode = http.StatusInternalServerError
	}

	statusText := http.StatusText(statusCode)
	logger := c.MustGet(config.LoggerContextKey).(*logger.Input)

	var responseError = errors.New(errors.AppError{
		Code:        statusText,
		StatusCode:  statusCode,
		Message:     statusText,
		SendToSlack: true,
		Logging:     true,
	})

	if hasError {
		if _, ok := requestErrors[0].Err.(*errors.AppError); !ok {
			responseError.Message = requestErrors[0].Error()
			if config.IsLocal {
				responseError.AddMetadata("errors", requestErrors.Errors())
			}
			logger.AddMetadata("errors", requestErrors.Errors())
		} else {
			responseError = requestErrors[0].Err.(*errors.AppError)
		}
	}

	logger.
		AddMetadata("ip", c.ClientIP()).
		AddMetadata("path", c.Request.URL.Path).
		AddMetadata("route_path", c.FullPath()).
		AddMetadata("method", c.Request.Method).
		AddMetadata("query", c.Request.URL.Query()).
		AddMetadata("version", c.Request.Proto).
		AddMetadata("referer", c.GetHeader("referer")).
		AddMetadata("agent", c.Request.UserAgent()).
		AddMetadata("body", c.Request.Form).
		AddMetadata("headers", c.Request.Header).
		AddMetadata("error", responseError)

	if config.IsDebug {
		if forwardedUser := c.GetHeader("X-Forwarded-User"); forwardedUser != "" {
			logger.AddMetadata("forwarded_user", forwardedUser)
		}
		if forwardedEmail := c.GetHeader("X-Forwarded-Email"); forwardedEmail != "" {
			logger.AddMetadata("forwarded_email", forwardedEmail)
		}
	}

	logger.Error("error")

	c.JSON(responseError.StatusCode, responseError)
	return
}
