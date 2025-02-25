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

	// actual creating the user
	tx := utils.CreateTransaction()
	var user *domain.User
	err = tx.A(func() error {
		_user, err := us.users.CreateOne(&pocketbase.CreateOneInput[domain.User]{
			Data: domain.User{
				Username:  input.Username,
				PushToken: input.PushToken,
				PublicKey: input.PublicKey,
			},
		})
		if err == nil {
			user = _user
		}
		return err
	}, func() error {
		return us.users.DeleteOne(&pocketbase.DeleteOneInput{
			Id: user.Id,
		})
	})
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occurred when creating a friendship record",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"input", input,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	// now we get auth token for the new user
	var token string
	err = tx.A(func() error {
		_token, err := domain.GetToken(user)
		if err == nil {
			token = _token
		}
		return err
	}, func() error { return nil })
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occurred when getting token",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"input", input,
			"user", user.Username,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	// encrypt token and user
	var marshalled []byte
	err = tx.A(func() error {
		_marshalled, err := json.Marshal(user)
		if err == nil {
			marshalled = _marshalled
		}
		return err
	}, func() error { return nil })
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occurred when marshalling user",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"input", input,
			"user", user.Username,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	var encryptedUser string
	err = tx.A(func() error {
		_encryptedUser, err := crypto.LongEncryptMessageRSA(
			string(marshalled),
			user.PublicKey,
		)
		if err == nil {
			encryptedUser = _encryptedUser
		}
		return err
	}, func() error { return nil })
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occurred when encrypting marshalled user",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"input", input,
			"user", user.Username,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	var encryptedToken string
	err = tx.A(func() error {
		_encryptedToken, err := crypto.LongEncryptMessageRSA(
			token,
			user.PublicKey,
		)
		if err == nil {
			encryptedToken = _encryptedToken
		}
		return err
	}, func() error { return nil })
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occurred when encrypting token",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"input", input,
			"user", user.Username,
		)
		return "", "", ErrSomethingWentWrongRegister
	}

	return encryptedUser, encryptedToken, nil
}
