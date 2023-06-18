package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack_alert"
	"net/http"
	"strings"
	"time"
)

func responseError(c *gin.Context) {
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
	var validations []any

	if hasRequestError {
		if _, ok := requestErrors[0].Err.(*errors.Input); !ok {
			appError = errors.New(errors.Input{
				StatusCode:  statusCode,
				SendToSlack: true,
			})

			if errs, ok := requestErrors[0].Err.(validator.ValidationErrors); ok {
				appError.Message = "Some fields are invalid"

				if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
					lang := en.New()
					trans, _ := ut.New(lang, lang).GetTranslator("en")
					_ = enTranslation.RegisterDefaultTranslations(val, trans)

					for _, e := range errs {
						validations = append(validations, map[string]any{
							"tag":       e.Tag(),
							"field":     e.Field(),
							"namespace": e.Namespace(),
							"value":     e.Value(),
							"param":     e.Param(),
						})
					}

					appError.Message = errs[0].Translate(trans)
					_ = appError.AddMetadata("validations", validations)
				}

				appError.StatusCode = http.StatusBadRequest
				appError.SendToSlack = false
			} else {
				appError.OriginalError = requestErrors.Last().Error()
			}
		} else {
			appError = requestErrors[0].Err.(*errors.Input)
		}
	}

	metadata["ip"] = c.ClientIP()
	metadata["time"] = time.Since(c.Writer.(*XResponseTimer).Start).String()
	metadata["path"] = path

	if routePath := c.FullPath(); routePath != "" {
		metadata["route_path"] = routePath
	}

	params := make(map[string]string)
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}
	metadata["params"] = params

	headers := make(map[string]string)
	for key, value := range c.Request.Header {
		valueAsString := value[0]
		if key == "Authorization" {
			tokenType := strings.Split(valueAsString, " ")[0]
			valueAsString = fmt.Sprintf("%s %s", tokenType, "***")
		}
		headers[strings.ToLower(key)] = valueAsString
	}
	metadata["headers"] = headers

	metadata["method"] = method
	metadata["query"] = c.Request.URL.Query()
	metadata["version"] = c.Request.Proto
	metadata["body"] = getRequestBody(c)
	metadata["error"] = appError

	if forwardedUser := c.GetHeader("X-Forwarded-User"); forwardedUser != "" {
		metadata["forwarded_user"] = forwardedUser
	}

	if forwardedEmail := c.GetHeader("X-Forwarded-Email"); forwardedEmail != "" {
		metadata["forwarded_email"] = forwardedEmail
	}

	appError.ErrorId = c.GetString(config.RequestIdCtxKey)
	c.MustGet(config.RequestLoggerCtxKey).(*logger.Logger).
		WithMetadata(metadata).
		Error("HTTP_REQUEST_ERROR")

	if config.IsLocal {
		c.JSON(appError.StatusCode, appError)
		return
	}

	if appError.SendToSlack {
		go slack_alert.New().WithRequestError(method, path, appError).Send()
	}

	errorMessage := appError.Message
	if appError.StatusCode == http.StatusInternalServerError {
		errorMessage = fmt.Sprintf(
			"An internal error occurred, contact the developers and enter the code [%s].",
			appError.ErrorId,
		)
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

func getRequestBody(c *gin.Context) map[string]any {
	result := make(map[string]any)
	if val, ok := c.Get(gin.BodyBytesKey); ok && val != nil {
		_ = json.Unmarshal(val.([]byte), &result)
	}
	return result
}
