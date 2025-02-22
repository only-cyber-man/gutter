package domain

import "github.com/tomek7667/cyberman-go/pocketbase"

type FriendshipStatus string

const (
	FriendshipStatusRequestSent FriendshipStatus = "request sent"
	FriendsStatus               FriendshipStatus = "friends"
)

func (s FriendshipStatus) String() string {
	return string(s)
}

type Friendship struct {
	pocketbase.PbItem

	RequesterId string           `json:"requester"`
	InviteeId   string           `json:"invitee"`
	Status      FriendshipStatus `json:"status"`

	E struct {
		Requester User `json:"requester"`
		Invitee   User `json:"invitee"`
	} `json:"expand,omitempty"`
}
