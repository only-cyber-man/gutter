package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/tomek7667/cyberman-go/rest"
)

func (c *Client) AddRoutes_ApiDebug() {
	slog.Warn("registering debug routes")
	c.RestClient.AddRoute("POST", "/api/debug/remove-users", func(ctx *gin.Context) {
		slog.Warn("removing all users")
		err := c.UserService.DeleteAll()
		rest.FailOrReturn(ctx, nil, err)
	})
}
