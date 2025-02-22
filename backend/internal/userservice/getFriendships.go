package userservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

func (us *Client) GetFriendships(requester *domain.User) ([]domain.KeyExchange, error) {
	keyExchanges, err := us.keyExchanges.GetFullList(&pocketbase.GetFullListInput[domain.KeyExchange]{
		Filter: pocketbase.BuildFilter(
			"relatedChat.requester.id = {:userId} || relatedChat.invitee.id = {:userId}",
			map[string]interface{}{
				"userId": requester.Id,
			},
		),
		Expand: "requester,target,relatedChat",
	})
	if err != nil {
		return nil, err
	}
	return keyExchanges, nil
}
