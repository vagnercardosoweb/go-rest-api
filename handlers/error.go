package handlers

import (
	"net/http"

	"rest-api/config"
	"rest-api/errors"
	"rest-api/shared"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context) {
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
	logger := shared.GetLogger()

	var responseError = errors.New(&errors.AppError{
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
		AddMetadata("path", c.Request.URL.Path).
		AddMetadata("method", c.Request.Method).
		AddMetadata("query", c.Request.URL.Query()).
		AddMetadata("body", c.Request.Form).
		AddMetadata("headers", c.Request.Header).
		AddMetadata("error", responseError).
		Error("Error to Json")

	c.JSON(responseError.StatusCode, responseError)
	return
}
