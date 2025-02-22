package userservice

import (
	"errors"
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

var ErrSomethingWentWrongGet = errors.New("something went wront while getting users")

func (us *Client) GetOneByUsername(username string) (*domain.User, error) {
	users, err := us.users.GetFullList(&pocketbase.GetFullListInput[domain.User]{
		Filter: pocketbase.BuildFilter(
			"username = {:username}",
			map[string]interface{}{
				"username": username,
			},
		),
	})
	if err != nil {
		slog.Error(
			"error occured when getting users by username",
			"err", err,
			"username", username,
		)
		return nil, ErrSomethingWentWrongGet
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return &users[0], nil
}
