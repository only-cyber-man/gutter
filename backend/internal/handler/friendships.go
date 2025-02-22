package handler

import (
	"log/slog"
	"net/http"

	"gutter/internal/domain"
	"gutter/internal/userservice"

	"github.com/gin-gonic/gin"
	"github.com/tomek7667/cyberman-go/rest"
)

func (c *Client) AddRoutes_ApiFriendships() {
	slog.Info("registering friendships api")
	c.RestClient.AddRoute("GET", "/api/friendships", c.WithAuth(), func(ctx *gin.Context) {
		user, _ := ctx.Get("user")
		output, err := c.UserService.GetFriendships(user.(*domain.User))
		rest.FailOrReturn(ctx, output, err)
	})
	c.RestClient.AddRoute("POST", "/api/friendships/invite", c.WithAuth(), func(ctx *gin.Context) {
		// TODO: add the gu_key_exchanges collection record
		user, _ := ctx.Get("user")
		dto := userservice.InviteDto{}
		err := ctx.ShouldBind(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		err = c.UserService.Invite(user.(*domain.User), &dto)
		rest.FailOrReturn(ctx, nil, err, "if the user exists, an invite has been sent")
	})
	c.RestClient.AddRoute("POST", "/api/friendships/answer", c.WithAuth(), func(ctx *gin.Context) {
		user, _ := ctx.Get("user")
		dto := userservice.AnswerDto{}
		err := ctx.ShouldBind(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		err = c.UserService.Answer(user.(*domain.User), &dto)
		rest.FailOrReturn(ctx, nil, err)
	})
}
