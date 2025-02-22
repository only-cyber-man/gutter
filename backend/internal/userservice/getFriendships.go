package userservice

import (
	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/pocketbase"
)

type user struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type outputFriendship struct {
	FriendshipId string                  `json:"friendshipId"`
	Requester    user                    `json:"requester"`
	Invitee      user                    `json:"invitee"`
	Status       domain.FriendshipStatus `json:"status"`
}

func (us *Client) GetFriendships(requester *domain.User) ([]outputFriendship, error) {
	friendships, err := us.friendshipsCollection.GetFullList(&pocketbase.GetFullListInput[domain.Friendship]{
		Filter: pocketbase.BuildFilter(
			"requester.id = {:userId} || invitee.id = {:userId}",
			map[string]interface{}{
				"userId": requester.Id,
			},
		),
		Expand: "requester,invitee",
	})
	if err != nil {
		return nil, err
	}
	var output []outputFriendship
	for _, f := range friendships {
		output = append(output, outputFriendship{
			FriendshipId: f.Id,
			Requester: user{
				Id:       f.E.Requester.Id,
				Username: f.E.Requester.Username,
			},
			Invitee: user{
				Id:       f.E.Invitee.Id,
				Username: f.E.Invitee.Username,
			},
			Status: f.Status,
		})
	}
	return output, nil
}
