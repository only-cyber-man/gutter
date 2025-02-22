package userservice

import (
	"errors"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

func (us *Client) DeleteAll() error {
	var errs error
	users, err := us.users.GetFullList(&pocketbase.GetFullListInput[domain.User]{})
	if err != nil {
		return err
	}
	for _, user := range users {
		err = us.users.DeleteOne(&pocketbase.DeleteOneInput{
			Id: user.Id,
		})
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}
