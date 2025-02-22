package handler

import (
	"log/slog"
	"net/http"

	"gutter/internal/domain"
	"gutter/internal/userservice"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/tomek7667/cyberman-go/rest"
)

func (c *Client) AddRoutes_ApiAuth() {
	slog.Info("registering auth api")
	c.RestClient.AddRateLimitedRoute("POST", "/api/auth/login", ratelimit.InMemoryOptions{}, func(ctx *gin.Context) {
		dto := userservice.LoginDto{}
		err := ctx.ShouldBind(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		user, token, err := c.UserService.Login(&dto)
		rest.FailOrReturn(ctx, map[string]interface{}{
			"user":  user,
			"token": token,
		}, err)
	})

	c.RestClient.AddRateLimitedRoute("POST", "/api/auth/register", ratelimit.InMemoryOptions{}, func(ctx *gin.Context) {
		dto := userservice.RegisterDto{}
		err := ctx.ShouldBind(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		user, token, err := c.UserService.Register(&dto)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		rest.FailOrReturn(ctx, map[string]interface{}{
			"user":  user,
			"token": token,
		}, err)
	})

	c.RestClient.AddRoute("DELETE", "/api/auth/account", c.WithAuth(), func(ctx *gin.Context) {
		user, _ := ctx.Get("user")
		err := c.UserService.DeleteAccount(user.(*domain.User))
		rest.FailOrReturn(ctx, nil, err)
	})
}
