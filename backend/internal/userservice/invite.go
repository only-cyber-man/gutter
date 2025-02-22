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
	Username            string `json:"username" binding:"required"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey" binding:"required"`
	PublicKey           string `json:"plaintextPublicKey" binding:"required"`
}

var (
	ErrSomethingWentWrongInvite = errors.New("something went wrong while inviting the user. Contact the administrator if the problem persists")
	ErrUserAlreadyInvited       = errors.New("the user is already invited")
	ErrUserAlreadyFriends       = errors.New("this user is already your friend")
	ErrCantInviteSelf           = errors.New("did you really try to invite yourself..? xD")
)

func (us *Client) Invite(requester *domain.User, input *InviteDto) error {
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
	friendships, err := us.friendships.GetFullList(&pocketbase.GetFullListInput[domain.Friendship]{
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
				"this shouldn't be possible. friendship record status is not defined",
				"record", friendshipRecord,
				"requester", requester.Username,
				"invitee", invitee.Username,
			)
			return ErrSomethingWentWrongInvite
		}
	}
	_, err = us.friendships.CreateOne(&pocketbase.CreateOneInput[domain.Friendship]{
		Data: domain.Friendship{
			RequesterId: requester.Id,
			InviteeId:   invitee.Id,
			Status:      domain.FriendshipStatusRequestSent,
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
	chat, err := us.chats.CreateOne(&pocketbase.CreateOneInput[domain.Chat]{
		Data: domain.Chat{
			CreatorId:       requester.Id,
			ParticipantsIds: []string{invitee.Id},
			PublicKey:       input.PublicKey,
		},
	})
	if err != nil {
		slog.Error(
			"error occured when creating a chat record",
			"err", err,
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		return ErrSomethingWentWrongInvite
	}
	_, err = us.keyExchanges.CreateOne(&pocketbase.CreateOneInput[domain.KeyExchange]{
		Data: domain.KeyExchange{
			EncryptedPrivateKey: input.EncryptedPrivateKey,
			RequesterId:         requester.Id,
			TargetId:            invitee.Id,
			RelatedChatId:       chat.Id,
		},
	})
	if err != nil {
		slog.Error(
			"error occured when creating a key exchange record",
			"err", err,
			"chat", chat,
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
