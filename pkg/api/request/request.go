package apirequest

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

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

func GetBodyAsMap(c *gin.Context) map[string]any {
	result := make(map[string]any)
	_ = json.Unmarshal(GetBodyAsBytes(c), &result)
	return result
}

func GetBodyAsRedacted(c *gin.Context) map[string]any {
	return utils.RedactKeys(GetBodyAsMap(c), nil)
}
