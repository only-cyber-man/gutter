package domain

import (
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type KeyExchange struct {
	pocketbase.PbItem

	// encrypted with the target's public key
	EncryptedPrivateKey string `json:"encryptedPrivateKey"`
	RequesterId         string `json:"requester"`
	TargetId            string `json:"target"`
	RelatedChatId       string `json:"relatedChat"`

	E struct {
		Requester   User `json:"requester"`
		Target      User `json:"target"`
		RelatedChat Chat `json:"relatedChat"`
	} `json:"expand,omitempty"`
}
