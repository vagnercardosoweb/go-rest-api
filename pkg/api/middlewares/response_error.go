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
	var logData = make(map[string]any)

	logData["ip"] = c.ClientIP()
	logData["request"] = fmt.Sprintf("%s %s", method, path)
	logData["time"] = time.Since(apicontext.StartTime(c)).String()

	if routePath := c.FullPath(); routePath != "" {
		logData["routePath"] = routePath
	}

	logData["statusCode"] = statusCode
	logData["params"] = getParams(c)
	logData["queryParams"] = getQueryParams(c)
	logData["headers"] = getHeaders(c)
	logData["body"] = apirequest.GetBodyAsRedacted(c)

	if !hasRequestError {
		logger.
			WithFields(logData).
			Error("HTTP_REQUEST_ERROR")
		return
	}

	var appError *errors.Input
	firstRequestError := requestErrors[0].Err

	if valueAsAppError, ok := firstRequestError.(*errors.Input); ok {
		appError = valueAsAppError
	} else {
		appError = errors.New(errors.Input{
			Message:    firstRequestError.Error(),
			StatusCode: statusCode,
		})
	}

	appError.RequestId = logger.GetId()
	logData["error"] = appError

	if env.IsLocal() {
		c.JSON(appError.StatusCode, appError)
		return
	}

	if *appError.Logging {
		logger.
			WithFields(logData).
			Error("HTTP_REQUEST_ERROR")
	}

	if *appError.SendAlert {
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
