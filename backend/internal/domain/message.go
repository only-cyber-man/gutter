package domain

import "github.com/tomek7667/cyberman-go/pocketbase"

type Message struct {
	pocketbase.PbItem

	// encrypted with the chat's public key
	EncryptedMessage     string `json:"encryptedMessage"`
	SenderUsernameAtTime string `json:"senderUsernameAtTime"`
	ChatId               string `json:"chat"`
	SenderId             string `json:"sender"`

	E struct {
		Chat   Chat `json:"chat"`
		Sender User `json:"sender"`
	} `json:"expand,omitempty"`
}
