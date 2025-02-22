package userservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

func (us *Client) DeleteAccount(requester *domain.User) error {
	return us.users.DeleteOne(&pocketbase.DeleteOneInput{
		Id: requester.Id,
	})
}
