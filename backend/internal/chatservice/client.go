package chatservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type Client struct {
	pbClient   *pocketbase.Client
	expoClient *expo.Client
	chats      *pocketbase.CollectionHandler[domain.Chat]
	messages   *pocketbase.CollectionHandler[domain.Message]
}

func New(pbClient *pocketbase.Client, expoClient *expo.Client) *Client {
	return &Client{
		pbClient:   pbClient,
		expoClient: expoClient,
		chats:      pocketbase.C[domain.Chat](pbClient, "gu_chats"),
		messages:   pocketbase.C[domain.Message](pbClient, "gu_messages"),
	}
}
