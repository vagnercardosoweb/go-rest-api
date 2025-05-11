package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
)

type ResponseTimer struct {
	gin.ResponseWriter
	Start time.Time
}

func (w *ResponseTimer) WriteHeader(statusCode int) {
	duration := time.Since(w.Start)
	w.Header().Set("X-Response-Time", duration.String())
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseTimer) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func ResponseTime(c *gin.Context) {
	start := time.Now()
	c.Set(apicontext.StartTimeKey, start)
	c.Writer = &ResponseTimer{c.Writer, start}
	c.Next()
}
