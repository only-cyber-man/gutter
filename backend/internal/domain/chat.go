package domain

import "github.com/tomek7667/cyberman-go/pocketbase"

type Chat struct {
	pocketbase.PbItem

	CreatorId       string   `json:"creator"`
	ParticipantsIds []string `json:"participants"`
	PublicKey       string   `json:"publicKey"`

	E struct {
		Creator      User   `json:"creator"`
		Participants []User `json:"participants"`
	} `json:"expand,omitempty"`
}
