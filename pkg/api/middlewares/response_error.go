package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	apirequest "github.com/vagnercardosoweb/go-rest-api/pkg/api/request"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"
)

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

	logger := apicontext.Logger(c)
	var metadata = make(map[string]any)

	metadata["ip"] = c.ClientIP()
	metadata["time"] = time.Since(apicontext.StartTime(c)).String()
	metadata["statusCode"] = statusCode
	metadata["path"] = fmt.Sprintf("%s %s", method, path)

	if routePath := c.FullPath(); routePath != "" {
		metadata["routePath"] = routePath
	}

	metadata["params"] = getParams(c)
	metadata["queryParams"] = getQueryParams(c)
	metadata["headers"] = getHeaders(c)
	metadata["body"] = apirequest.GetBodyAsRedacted(c)

	if !hasRequestError {
		logger.
			WithMetadata(metadata).
			Error("HTTP_REQUEST_ERROR")
		return
	}

	var appError *errors.Input
	originalError := requestErrors[0].Err

	if valueAsAppError, ok := originalError.(*errors.Input); ok {
		appError = valueAsAppError
	} else {
		appError = errors.New(errors.Input{
			Message:    originalError.Error(),
			StatusCode: statusCode,
		})
	}

	appError.RequestId = logger.GetId()
	metadata["error"] = appError

	if *appError.Logging {
		logger.
			WithMetadata(metadata).
			Error("HTTP_REQUEST_ERROR")
	}

	if env.IsLocal() {
		c.JSON(appError.StatusCode, appError)
		return
	}

	if *appError.SendToSlack {
		go func() {
			_ = slack.
				NewAlert().
				WithRequestError(method, path, appError).
				Send()
		}()
	}

	errorMessage := appError.Message
	if appError.StatusCode == http.StatusInternalServerError {
		errorMessage = fmt.Sprintf(
			`An internal error occurred, contact the developers and enter the code "%s".`,
			appError.RequestId,
		)
	}

	validations := make([]map[string]any, 0)
	if v, ok := appError.Metadata["validations"]; ok {
		validations = v.([]map[string]any)
	}

	c.JSON(appError.StatusCode, &response{
		Code:        appError.Code,
		Name:        appError.Name,
		RequestId:   appError.RequestId,
		StatusCode:  appError.StatusCode,
		Validations: validations,
		Message:     errorMessage,
	})
}

func getParams(c *gin.Context) map[string]string {
	params := make(map[string]string)
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}
	return params
}

func getQueryParams(c *gin.Context) map[string]string {
	queryParams := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		queryParams[key] = value[0]
	}
	return queryParams
}

func getHeaders(c *gin.Context) map[string]string {
	headers := make(map[string]string)
	for key, value := range c.Request.Header {
		headers[strings.ToLower(key)] = value[0]
	}
	return headers
}

type response struct {
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	RequestId   string           `json:"requestId"`
	Validations []map[string]any `json:"validations"`
	StatusCode  int              `json:"statusCode"`
	Message     string           `json:"message"`
}
