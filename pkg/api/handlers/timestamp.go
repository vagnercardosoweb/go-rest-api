package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apiutils "github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

func Timestamp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"utc":      time.Now().UTC(),
		"duration": time.Since(apiutils.GetRequestStartTime(c)).String(),
		"brl":      time.Now().In(utils.LocationBrl),
	})
}
