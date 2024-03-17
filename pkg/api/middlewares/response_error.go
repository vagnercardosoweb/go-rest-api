package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"
)

type ResponseErrorOutput struct {
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	RequestId   string           `json:"requestId"`
	StatusCode  int              `json:"statusCode"`
	Validations []map[string]any `json:"validations"`
	Message     string           `json:"message"`
}

func ResponseError(c *gin.Context) {
	c.Next()

	requestErrors := c.Errors
	hasRequestError := len(requestErrors) > 0
	isAborted := c.IsAborted()

	if !hasRequestError && !isAborted {
		return
	}

	path := c.Request.URL.Path
	statusCode := c.Writer.Status()
	method := c.Request.Method

	if isAborted && statusCode == http.StatusOK {
		statusCode = http.StatusInternalServerError
	}

	var metadata = make(map[string]any)

	metadata["ip"] = c.ClientIP()
	metadata["time"] = time.Since(utils.GetRequestStartTime(c)).String()
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
	metadata["body"] = utils.GetBodyAsJson(c)

	if forwardedUser := c.GetHeader("X-Forwarded-User"); forwardedUser != "" {
		metadata["forwardedUser"] = forwardedUser
	}

	if forwardedEmail := c.GetHeader("X-Forwarded-Email"); forwardedEmail != "" {
		metadata["forwardedEmail"] = forwardedEmail
	}

	var appError *errors.Input
	logger := utils.GetLogger(c)
	requestId := logger.GetId()

	if !hasRequestError {
		logger.WithMetadata(metadata).Error("HTTP_REQUEST_ERROR")
		return
	}

	originalError := requestErrors[0].Err
	if valueAsError, ok := originalError.(*errors.Input); ok {
		appError = valueAsError
	} else {
		appError = errors.New(errors.Input{
			Message:    originalError.Error(),
			StatusCode: statusCode,
		})
	}

	appError.ErrorId = requestId
	metadata["error"] = appError

	if *appError.Logging {
		logger.WithMetadata(metadata).Error("HTTP_REQUEST_ERROR")
	}

	if env.GetAppEnv() == env.AppLocal {
		c.JSON(appError.StatusCode, appError)
		return
	}

	if *appError.SendToSlack {
		go slack.
			NewAlert().
			WithRequestError(method, path, appError).
			Send()
	}

	errorMessage := appError.Message
	if appError.StatusCode == http.StatusInternalServerError {
		errorMessage = fmt.Sprintf(
			`An internal error occurred, contact the developers and enter the code "%s".`,
			appError.ErrorId,
		)
	}

	validations := make([]map[string]any, 0)
	if v, ok := appError.Metadata["validations"]; ok {
		validations = v.([]map[string]any)
	}

	c.JSON(appError.StatusCode, &ResponseErrorOutput{
		Code:        appError.Code,
		Name:        appError.Name,
		RequestId:   requestId,
		StatusCode:  appError.StatusCode,
		Validations: validations,
		Message:     errorMessage,
	})
}
