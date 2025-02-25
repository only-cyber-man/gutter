package chatservice

import (
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

type SendMessageDto struct {
	EncryptedMessage string `json:"encryptedMessage" binding:"required"`
}

func (cs *Client) SendMessage(requester *domain.User, chatId string, input *SendMessageDto) error {
	chats, err := cs.chats.GetFullList(&pocketbase.GetFullListInput[domain.Chat]{
		Filter: pocketbase.BuildFilter(
			"chat.id = {:chatId} && (creator.id = {:userId} || participants.id ?= {:userId})",
			map[string]any{
				"chatId": chatId,
				"userId": requester.Id,
			},
		),
	})
	slog.Info("chats", "c", chats, "err", err)
	return nil
}
