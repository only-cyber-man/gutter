package domain

import (
	"encoding/json"
	"log/slog"

	"github.com/tomek7667/cyberman-go/crypto"
	"github.com/tomek7667/cyberman-go/pocketbase"
	"github.com/tomek7667/cyberman-go/utils"
)

type User struct {
	pocketbase.PbItem

	Username          string `json:"username,omitempty"`
	PushToken         string `json:"pushToken,omitempty"`
	EncryptedPassword string `json:"encryptedPassword,omitempty"`
	PublicKey         string `json:"publicKey,omitempty"`
}

func (u *User) Compare(currentPassword string) bool {
	decryptedPassword, err := crypto.DecryptAES256(
		u.EncryptedPassword,
		utils.Getenv("AES_KEY", "ba7816bf8f01cfea414140de5dae2223"),
	)
	if err != nil {
		slog.Info(
			"comparing password error",
			"current password", crypto.Obfuscate(currentPassword, 3),
			"key", crypto.Obfuscate(
				utils.Getenv("AES_KEY", "ba7816bf8f01cfea414140de5dae2223"),
				4,
			),
			"err", err,
		)
	}
	return decryptedPassword == currentPassword
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
