package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

func newRedactTestContext(target string) *gin.Context {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, target, nil)

	return c
}

func TestRedactedQueryParams(t *testing.T) {
	c := newRedactTestContext("/test?name=John&bearerToken=jwt-token&tag=go&tag=api&accessToken=access-token")

	result := redactedQueryParams(c)
	assert.Equal(t, map[string]any{
		"name":        "John",
		"bearerToken": utils.RedactedValue,
		"accessToken": utils.RedactedValue,
		"tag":         []any{"go", "api"},
	}, result)
}

func TestRedactedHeaders(t *testing.T) {
	c := newRedactTestContext("/test")

	c.Request.Header.Set("Authorization", "Bearer jwt-token")
	c.Request.Header.Set("X-Api-Key", "api-key")
	c.Request.Header.Set("User-Agent", "go-test")
	c.Request.Header.Add("Accept", "application/json")
	c.Request.Header.Add("Accept", "application/problem+json")

	result := redactedHeaders(c)
	assert.Equal(t, map[string]any{
		"authorization": utils.RedactedValue,
		"x-api-key":     utils.RedactedValue,
		"accept":        []any{"application/json", "application/problem+json"},
		"user-agent":    "go-test",
	}, result)
}

func TestStringValuesAsLogValue(t *testing.T) {
	assert.Equal(t, "single", stringValuesAsLogValue([]string{"single"}))
	assert.Equal(t, []any{"first", "second"}, stringValuesAsLogValue([]string{"first", "second"}))
}
