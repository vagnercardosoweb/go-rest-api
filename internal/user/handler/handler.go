package user

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/database/postgres"
	"github.com/vagnercardosoweb/go-rest-api/sqlc/store"
	"net/http"
)

type handler struct{}

func (h handler) List(c *gin.Context) {
	conn := c.MustGet(config.PgConnectCtxKey).(*postgres.Connection)
	queries := store.New(conn.GetSqlx())
	users, err := queries.GetUsers(c.Request.Context(), 10)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var results []map[string]any
	for _, user := range users {
		var confirmedEmailAt any
		if user.ConfirmedEmailAt.Valid {
			confirmedEmailAt = user.ConfirmedEmailAt.Time
		}
		results = append(results, map[string]any{
			"id":                 user.ID,
			"name":               user.Name,
			"email":              user.Email,
			"confirmed_email_at": confirmedEmailAt,
			"code_to_invite":     user.CodeToInvite,
			"birth_date":         user.BirthDate.Format("2006-01-02"),
		})
	}
	c.JSON(http.StatusOK, results)
}

func New() *handler {
	return &handler{}
}
