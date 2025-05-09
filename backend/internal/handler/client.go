package handler

import (
	"log/slog"

	"gutter/internal/chatservice"
	"gutter/internal/userservice"

	"github.com/gin-gonic/gin"
	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
	"github.com/tomek7667/cyberman-go/rest"
)

type Client struct {
	RestClient  *rest.Client
	PbClient    *pocketbase.Client
	UserService *userservice.Client
	ChatService *chatservice.Client
}

func New(
	restClient *rest.Client,
	pbClient *pocketbase.Client,
	expoClient *expo.Client,
) *Client {
	return &Client{
		RestClient: restClient,
		PbClient:   pbClient,
		UserService: userservice.New(
			pbClient,
			expoClient,
		),
		ChatService: chatservice.New(
			pbClient,
			expoClient,
		),
	}
}

func (c *Client) Start() {
	slog.Info(
		"Starting gutter handler",
		"port", c.RestClient.Port,
		"mode", gin.Mode(),
	)
	c.AddRoutes_ApiBase()
	c.AddRoutes_ApiAuth()
	c.AddRoutes_ApiFriendships()
	c.AddRoutes_ApiChats()
	if gin.Mode() != gin.ReleaseMode {
		// these are not protected at all
		c.AddRoutes_ApiDebug()
	}
	c.RestClient.Serve()
}
