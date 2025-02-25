package handler

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomek7667/cyberman-go/rest"
)

func (c *Client) AddRoutes_ApiBase() {
	slog.Info("registering base api")

	c.RestClient.AddRoute("GET", "/api", func(ctx *gin.Context) {
		rest.FailOrReturn(ctx, map[string]any{
			"success":   true,
			"timestamp": time.Now(),
		}, nil)
	})
}
