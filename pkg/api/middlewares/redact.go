package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

var sensitiveLogKeys = []string{
	"cookie",
	"set-cookie",
	"x-api-key",
	"bearerToken",
	"authorization",
	"refreshToken",
	"x-internal-key",
	"x-refresh-token",
	"accessToken",
	"x-id-token",
	"idToken",
	"token",
}

func redactedQueryParams(c *gin.Context) map[string]any {
	return redactedStringValues(c.Request.URL.Query(), false)
}

func redactedHeaders(c *gin.Context) map[string]any {
	return redactedStringValues(c.Request.Header, true)
}

func redactedStringValues(values map[string][]string, normalizeKeys bool) map[string]any {
	data := make(map[string]any, len(values))

	for key, values := range values {
		if normalizeKeys {
			key = strings.ToLower(key)
		}

		data[key] = stringValuesAsLogValue(values)
	}

	return utils.RedactKeys(data, sensitiveLogKeys)
}

func stringValuesAsLogValue(values []string) any {
	if len(values) == 1 {
		return values[0]
	}

	items := make([]any, len(values))
	for i, value := range values {
		items[i] = value
	}

	return items
}
