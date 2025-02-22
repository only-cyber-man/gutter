package userservice

import (
	"errors"
	"fmt"
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type InviteDto struct {
	Username string `json:"username" binding:"required"`
}

var (
	ErrSomethingWentWrongInvite = errors.New("something went wrong while inviting the user. Contact the administrator if the problem persists")
	ErrUserAlreadyInvited       = errors.New("the user is already invited")
	ErrUserAlreadyFriends       = errors.New("this user is already your friend")
	ErrCantInviteSelf           = errors.New("did you really try to invite yourself..? xD")
)

func (us *Client) Invite(requester *domain.User, input *InviteDto) error {
	users, err := us.usersCollection.GetFullList(&pocketbase.GetFullListInput[domain.User]{
		Filter: pocketbase.BuildFilter(
			"username = {:username}",
			map[string]interface{}{
				"username": input.Username,
			},
		),
	})
	if err != nil {
		slog.Error(
			"error occured when checking if the user exists in invite",
			"err", err,
			"tried to invite", input.Username,
		)
		return ErrSomethingWentWrongInvite
	}
	if len(users) == 0 {
		slog.Warn(
			"user tried to invite, but the invitee doesn't exist",
			"requester", requester.Username,
			"target", input.Username,
		)
		return ErrUserNotFound
	}
	invitee := users[0]
	if invitee.Id == requester.Id {
		return ErrCantInviteSelf
	}
	friendships, err := us.friendshipsCollection.GetFullList(&pocketbase.GetFullListInput[domain.Friendship]{
		Filter: pocketbase.BuildFilter(
			"requester.id = {:requesterId} && invitee.username = {:inviteeUsername}",
			map[string]interface{}{
				"requesterId":     requester.Id,
				"inviteeUsername": invitee.Username,
			},
		),
	})
	if err != nil {
		slog.Error(
			"error occured when getting friendships",
			"err", err,
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		return ErrSomethingWentWrongInvite
	}
	if len(friendships) > 0 {
		friendshipRecord := friendships[0]
		switch friendshipRecord.Status {
		case domain.FriendsStatus:
			return ErrUserAlreadyFriends
		case domain.FriendshipStatusRequestSent:
			return ErrUserAlreadyInvited
		default:
			slog.Error(
				"this shouldn't be possible. friendship record status is not available",
				"record", friendshipRecord,
				"requester", requester.Username,
				"invitee", invitee.Username,
			)
			return ErrSomethingWentWrongInvite
		}
	}
	_, err = us.friendshipsCollection.CreateOne(&pocketbase.CreateOneInput[domain.Friendship]{
		Data: domain.Friendship{
			Requester: requester.Id,
			Invitee:   invitee.Id,
			Status:    domain.FriendshipStatusRequestSent,
		},
	})
	if err != nil {
		slog.Error(
			"error occured when creating a friendship record",
			"err", err,
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		return ErrSomethingWentWrongInvite
	}

	// here create a new expo push token to the invitee
	if invitee.PushToken != "" {
		slog.Debug(
			"trying to send a notification regarding and invite",
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		err = us.expoClient.SendNotification(&expo.SendNotificationInput{
			Title: "New friendship request",
			Body: fmt.Sprintf(
				"%s sent you an invitation",
				requester.Username,
			),
		}, invitee.PushToken)
		if err != nil {
			slog.Error(
				"error occured when sending a notification about a new friendship",
				"err", err,
				"requester", requester.Username,
				"invitee", invitee.Username,
			)
		}
	}

	return nil
}
