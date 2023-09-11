package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

type XResponseTimer struct {
	gin.ResponseWriter
	Start time.Time
}

func (w *XResponseTimer) WriteHeader(statusCode int) {
	duration := time.Since(w.Start)
	w.Header().Set("X-Response-Time", duration.String())
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *XResponseTimer) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func ResponseTime(c *gin.Context) {
	blw := &XResponseTimer{ResponseWriter: c.Writer, Start: time.Now()}
	c.Writer = blw
	c.Next()
}
