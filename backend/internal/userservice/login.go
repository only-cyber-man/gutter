package userservice

import (
	"errors"
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/crypto"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type LoginDto struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	PushToken string `json:"pushToken,omitempty"`
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
)

func (us *Client) Login(input *LoginDto) (*domain.User, string, error) {
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
			"trial password", crypto.Obfuscate(input.Password, 5),
			"error", err,
		)
		return nil, "", ErrUserNotFound
	}
	if len(users) == 0 {
		return nil, "", ErrUserNotFound
	}

	retrievedUser := users[0]
	isCorrectPassword := retrievedUser.Compare(input.Password)
	if !isCorrectPassword {
		slog.Warn(
			"login invalid attempt",
			"username", input.Username,
			"trial password", crypto.Obfuscate(input.Password, 5),
			"error", err,
		)
		return nil, "", ErrIncorrectPassword
	}

	_, err = us.users.UpdateOne(&pocketbase.UpdateOneInput[domain.User]{
		Id: users[0].Id,
		Data: domain.User{
			PushToken: input.PushToken,
		},
	})
	if err != nil {
		slog.Error(
			"update one with push token and username of users failed",
			"input dto", input,
			"err", err,
		)
	}
	tk, err := domain.GetToken(&retrievedUser)
	return &users[0], tk, err
}
