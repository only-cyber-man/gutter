package userservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type Client struct {
	pbClient              *pocketbase.Client
	expoClient            *expo.Client
	usersCollection       *pocketbase.CollectionHandler[domain.User]
	friendshipsCollection *pocketbase.CollectionHandler[domain.Friendship]
}

func New(pbClient *pocketbase.Client, expoClient *expo.Client) *Client {
	return &Client{
		pbClient:              pbClient,
		expoClient:            expoClient,
		usersCollection:       pocketbase.C[domain.User](pbClient, "gu_users"),
		friendshipsCollection: pocketbase.C[domain.Friendship](pbClient, "gu_friendships"),
	}
}
