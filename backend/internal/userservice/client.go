package userservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type Client struct {
	pbClient     *pocketbase.Client
	expoClient   *expo.Client
	users        *pocketbase.CollectionHandler[domain.User]
	friendships  *pocketbase.CollectionHandler[domain.Friendship]
	keyExchanges *pocketbase.CollectionHandler[domain.KeyExchange]
	chats        *pocketbase.CollectionHandler[domain.Chat]
}

func New(pbClient *pocketbase.Client, expoClient *expo.Client) *Client {
	return &Client{
		pbClient:     pbClient,
		expoClient:   expoClient,
		users:        pocketbase.C[domain.User](pbClient, "gu_users"),
		friendships:  pocketbase.C[domain.Friendship](pbClient, "gu_friendships"),
		keyExchanges: pocketbase.C[domain.KeyExchange](pbClient, "gu_key_exchanges"),
		chats:        pocketbase.C[domain.Chat](pbClient, "gu_chats"),
	}
}
