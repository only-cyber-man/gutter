package userservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

func (us *Client) Verify(token string) (*domain.User, error) {
	userId, err := domain.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	user, err := us.users.GetOne(&pocketbase.GetOneInput[domain.User]{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
