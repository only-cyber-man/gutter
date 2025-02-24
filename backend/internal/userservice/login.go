package userservice

import (
	"encoding/json"
	"errors"
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/crypto"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type LoginDto struct {
	Username string `json:"username" binding:"required"`
}

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrSomethingWentWrongLogin = errors.New("something went wrong while logging in")
)

func (us *Client) Login(input *LoginDto) (string, string, error) {
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
			"login get full list failed",
			"username", input.Username,
			"error", err,
		)
		return "", "", ErrUserNotFound
	}
	if len(users) == 0 {
		return "", "", ErrUserNotFound
	}

	retrievedUser := users[0]
	token, err := domain.GetToken(&retrievedUser)
	if err != nil {
		return "", "", ErrSomethingWentWrongLogin
	}

	// encrypt token and user
	marshalled, err := json.Marshal(retrievedUser)
	if err != nil {
		slog.Error(
			"failed marshalling user",
			"err", err,
		)
		return "", "", ErrSomethingWentWrongLogin
	}

	// encrypting user/token
	encryptedUser, err := crypto.LongEncryptMessageRSA(
		string(marshalled),
		retrievedUser.PublicKey,
	)
	if err != nil {
		slog.Error(
			"failed encrypting marshalled user",
			"err", err,
		)
		return "", "", ErrSomethingWentWrongLogin
	}

	encryptedToken, err := crypto.LongEncryptMessageRSA(
		token,
		retrievedUser.PublicKey,
	)
	if err != nil {
		slog.Error(
			"failed encrypting token",
			"err", err,
			"tk", token,
			"user id", retrievedUser.Id,
			"public key", retrievedUser.PublicKey,
		)
		return "", "", ErrSomethingWentWrongLogin
	}
	return encryptedUser, encryptedToken, nil
}
