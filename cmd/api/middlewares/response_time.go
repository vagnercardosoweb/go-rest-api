package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

type XResponseTimer struct {
	gin.ResponseWriter
	start time.Time
}

func (w *XResponseTimer) WriteHeader(statusCode int) {
	duration := time.Since(w.start)
	w.Header().Set("X-Response-Time", duration.String())
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *XResponseTimer) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func responseTimer(c *gin.Context) {
	blw := &XResponseTimer{ResponseWriter: c.Writer, start: time.Now()}
	c.Writer = blw
	c.Next()
}
