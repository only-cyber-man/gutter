package domain

import (
	"time"

	"github.com/tomek7667/cyberman-go/crypto"
	"github.com/tomek7667/cyberman-go/utils"
)

func GetToken(user *User) (string, error) {
	return crypto.JwtEncode(
		map[string]interface{}{
			"sub":       user.GetId(),
			"id":        user.GetId(),
			"iat":       time.Now().Unix(),
			"username":  user.Username,
			"createdAt": user.CreatedAt,
		},
		utils.Getenv(
			"JWT_SECRET",
			"8c7fafb856380624fa60b22e7baf311d",
		),
	)
}

func VerifyToken(token string) (string, error) {
	data, err := crypto.JwtVerify(token, utils.Getenv(
		"JWT_SECRET",
		"8c7fafb856380624fa60b22e7baf311d",
	))
	if err != nil {
		return "", err
	}
	return data["id"].(string), nil
}
