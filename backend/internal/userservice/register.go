package userservice

import (
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
	ErrUsernameAlreadyTaken = errors.New("this username is already taken")
	ErrSomethingWentWrong   = errors.New("registration failed, please try again later or contact the administrator")
)

func (us *Client) Register(input *RegisterDto) (*domain.User, string, error) {
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
		return nil, "", err
	}
	if len(users) > 0 {
		return nil, "", ErrUsernameAlreadyTaken
	}
	encryptedPassword := crypto.EncryptAES256(
		input.Password,
		utils.Getenv("AES_KEY", "ba7816bf8f01cfea414140de5dae2223"),
	)
	user, err := us.users.CreateOne(&pocketbase.CreateOneInput[domain.User]{
		Data: domain.User{
			Username:          input.Username,
			EncryptedPassword: encryptedPassword,
			PushToken:         input.PushToken,
			PublicKey:         input.PublicKey,
		},
	})
	if err != nil {
		slog.Warn(
			"create user failed",
			"err", err,
			"user", user,
		)
		return nil, "", err
	}
	token, err := domain.GetToken(user)
	if err != nil {
		slog.Error(
			"get token failed",
			"err", err,
			"user", user,
		)
		return nil, "", err
	}
	return user, token, nil
}
