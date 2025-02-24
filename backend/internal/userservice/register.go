package userservice

import (
	"encoding/json"
	"errors"
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/crypto"
	"github.com/tomek7667/cyberman-go/pocketbase"
	"github.com/tomek7667/cyberman-go/utils"
)

type RegisterDto struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	PushToken string `json:"pushToken,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

var (
	ErrUsernameAlreadyTaken       = errors.New("this username is already taken")
	ErrSomethingWentWrongRegister = errors.New("registration failed, please try again later or contact the administrator")
)

func (us *Client) Register(input *RegisterDto) (string, string, error) {
	if err := utils.CheckUsername(input.Username); err != nil {
		return "", "", err
	}
	users, err := us.users.GetFullList(&pocketbase.GetFullListInput[domain.User]{
		Filter: pocketbase.BuildFilter(
			"username = {:username}",
			map[string]interface{}{
				"username": input.Username,
			},
		),
	})
	if err != nil {
		slog.Error(
			"get full list to check whether the user already exists failed",
			"err", err,
			"input", input,
		)
		return "", "", err
	}
	if len(users) > 0 {
		return "", "", ErrUsernameAlreadyTaken
	}
	user, err := us.users.CreateOne(&pocketbase.CreateOneInput[domain.User]{
		Data: domain.User{
			Username:  input.Username,
			PushToken: input.PushToken,
			PublicKey: input.PublicKey,
		},
	})
	if err != nil {
		slog.Warn(
			"create user failed",
			"err", err,
			"user", user,
		)
		return "", "", err
	}
	token, err := domain.GetToken(user)
	if err != nil {
		slog.Error(
			"get token failed",
			"err", err,
			"user", user,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	// encrypt token and user
	marshalled, err := json.Marshal(user)
	if err != nil {
		slog.Error(
			"failed marshalling user",
			"err", err,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	// encrypting user/token
	encryptedUser, err := crypto.LongEncryptMessageRSA(
		string(marshalled),
		user.PublicKey,
	)
	if err != nil {
		slog.Error(
			"failed encrypting marshalled user",
			"err", err,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	encryptedToken, err := crypto.LongEncryptMessageRSA(
		token,
		user.PublicKey,
	)
	if err != nil {
		slog.Error(
			"failed encrypting token",
			"err", err,
			"tk", token,
			"user id", user.Id,
			"public key", user.PublicKey,
		)
		return "", "", ErrSomethingWentWrongRegister
	}
	return encryptedUser, encryptedToken, nil
}
