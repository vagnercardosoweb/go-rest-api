package apirequest

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func BodyAsBytes(c *gin.Context) []byte {
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

func BodyAsJson(c *gin.Context) map[string]any {
	bodyAsBytes := BodyAsBytes(c)
	result := make(map[string]any)
	_ = json.Unmarshal(bodyAsBytes, &result)
	return result
}
