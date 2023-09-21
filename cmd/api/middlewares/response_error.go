package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack_alert"
)

func ResponseError(c *gin.Context) {
	c.Next()

	requestErrors := c.Errors
	hasRequestError := len(requestErrors) > 0
	isAborted := c.IsAborted()

	if !hasRequestError && !isAborted {
		return
	}

	path := c.Request.URL.String()
	statusCode := c.Writer.Status()
	method := c.Request.Method

	if (isAborted && statusCode == http.StatusOK) || hasRequestError {
		statusCode = http.StatusInternalServerError
	}

	var metadata = make(map[string]any, 0)
	var appError *errors.Input

	if hasRequestError {
		originalError := requestErrors[0].Err
		if valueAsError, ok := originalError.(*errors.Input); ok {
			appError = valueAsError
		} else {
			appError = errors.New(errors.Input{
				OriginalError: originalError,
				StatusCode:    statusCode,
			})
		}
	}

	metadata["ip"] = c.ClientIP()
	metadata["time"] = time.Since(c.Writer.(*XResponseTimer).Start).String()
	metadata["path"] = path

	if routePath := c.FullPath(); routePath != "" {
		metadata["routePath"] = routePath
	}

	params := make(map[string]string)
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}
	metadata["params"] = params

	headers := make(map[string]string)
	for key, value := range c.Request.Header {
		headers[strings.ToLower(key)] = value[0]
	}
	metadata["headers"] = headers

	metadata["method"] = method
	metadata["query"] = c.Request.URL.Query()
	metadata["version"] = c.Request.Proto
	metadata["body"] = GetBodyAsJson(c)
	metadata["error"] = appError

	if forwardedUser := c.GetHeader("X-Forwarded-User"); forwardedUser != "" {
		metadata["forwardedUser"] = forwardedUser
	}

	if forwardedEmail := c.GetHeader("X-Forwarded-Email"); forwardedEmail != "" {
		metadata["forwardedEmail"] = forwardedEmail
	}

	logger := config.GetLoggerFromCtx(c)
	appError.ErrorId = config.GetRequestIdFromCtx(c)

	if *appError.Logging {
		logger.WithMetadata(metadata).Error("HTTP_REQUEST_ERROR")
	}

	if config.IsLocal {
		c.JSON(appError.StatusCode, appError)
		return
	}

	if *appError.SendToSlack {
		go slack_alert.NewClient().WithRequestError(method, path, appError).Send()
	}

	errorMessage := appError.Message
	if appError.StatusCode == http.StatusInternalServerError {
		errorMessage = fmt.Sprintf(
			"An internal error occurred, contact the developers and enter the code [%s].",
			appError.ErrorId,
		)
	}

	var validations []map[string]any
	if v, ok := appError.Metadata["validations"]; ok {
		validations = v.([]map[string]any)
	}

	c.JSON(appError.StatusCode, gin.H{
		"name":        appError.Name,
		"code":        appError.Code,
		"errorId":     appError.ErrorId,
		"statusCode":  appError.StatusCode,
		"validations": validations,
		"message":     errorMessage,
	})
}

func GetBodyAsBytes(c *gin.Context) []byte {
	bodyAsBytes := []byte("{}")
	if val, ok := c.Get(gin.BodyBytesKey); ok && val != nil {
		bodyAsBytes = val.([]byte)
	} else {
		b, _ := io.ReadAll(c.Request.Body)
		if len(b) > 0 {
			c.Set(gin.BodyBytesKey, b)
			bodyAsBytes = b
		}
	}
	return bodyAsBytes
}

func GetBodyAsJson(c *gin.Context) map[string]any {
	bodyAsBytes := GetBodyAsBytes(c)
	result := make(map[string]any)
	_ = json.Unmarshal(bodyAsBytes, &result)
	return result
}
