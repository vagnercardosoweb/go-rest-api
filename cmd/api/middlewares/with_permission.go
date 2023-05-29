package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
	"net/http"
	"strings"
)

type PermissionRow struct {
	ID   uuid.UUID
	Name string
}

func WithPermission(c *gin.Context) {
	path := c.Request.URL.String()
	method := strings.ToUpper(c.Request.Method)

	if routePath := c.FullPath(); routePath != "" {
		path = routePath
	}

	db := c.MustGet(config.PgConnectCtxKey).(*postgres.Connection)
	tokenOutput := c.MustGet(config.TokenPayloadCtxKey).(*token.Output)
	pathWithMethod := fmt.Sprintf("%s %s", method, path)

	var permissionRow PermissionRow
	if err := db.QueryOne(&permissionRow, `
		SELECT
			permissions.id,
			permissions.name
		FROM permissions
				 INNER JOIN users_permissions ON permissions.id = users_permissions.permission_id
		WHERE users_permissions.user_id = $1
		  AND permissions.name = $2
		ORDER BY permissions.created_at DESC
		LIMIT 1;
	`, tokenOutput.Subject, pathWithMethod); err != nil {
		abortError := errors.New(errors.Input{
			Message:       "You do not have permission to access this resource. Please contact administrators.",
			StatusCode:    http.StatusForbidden,
			OriginalError: err,
			Metadata: errors.Metadata{
				"userId":         tokenOutput.Subject,
				"permissionName": pathWithMethod,
			},
		})
		c.AbortWithError(abortError.StatusCode, abortError)
		return
	}

	c.Next()
}
