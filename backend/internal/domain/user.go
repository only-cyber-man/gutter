package domain

import (
	"encoding/json"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

type User struct {
	pocketbase.PbItem

	Username  string `json:"username,omitempty"`
	PushToken string `json:"pushToken,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id        string `json:"id,omitempty"`
		Username  string `json:"username,omitempty"`
		PublicKey string `json:"publicKey,omitempty"`
	}{
		Id:        u.Id,
		Username:  u.Username,
		PublicKey: u.PublicKey,
	})
}
