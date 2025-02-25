package handler

import (
	"log/slog"
	"net/http"

	"gutter/internal/chatservice"
	"gutter/internal/domain"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/tomek7667/cyberman-go/rest"
)

func (c *Client) AddRoutes_ApiChats() {
	slog.Info("registering chats api")

	c.RestClient.AddRateLimitedRoute("POST", "/api/chats/:chatId/messages", ratelimit.InMemoryOptions{}, c.WithAuth(), func(ctx *gin.Context) {
		user, _ := ctx.Get("user")
		chatId := ctx.Param("chatId")
		dto := chatservice.SendMessageDto{}
		err := ctx.ShouldBind(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		err = c.ChatService.SendMessage(
			user.(*domain.User),
			chatId,
			&dto,
		)
		rest.FailOrReturn(ctx, nil, err)
	})
}
