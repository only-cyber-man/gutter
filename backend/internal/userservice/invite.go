package userservice

import (
	"errors"
	"fmt"
	"log/slog"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
	"github.com/tomek7667/cyberman-go/utils"
)

type InviteDto struct {
	Username            string `json:"username" binding:"required"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey" binding:"required"`
	ChatPublicKey       string `json:"chatPublicKey" binding:"required"`
}

var (
	ErrSomethingWentWrongInvite = errors.New("something went wrong while inviting the user. Contact the administrator if the problem persists")
	ErrUserAlreadyInvited       = errors.New("the user is already invited")
	ErrUserAlreadyFriends       = errors.New("this user is already your friend")
	ErrCantInviteSelf           = errors.New("did you really try to invite yourself..? xD")
)

func (us *Client) Invite(requester *domain.User, input *InviteDto) (*domain.Chat, error) {
	tx := utils.CreateTransaction()
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
		return nil, ErrSomethingWentWrongInvite
	}
	if len(users) == 0 {
		slog.Warn(
			"user tried to invite, but the invitee doesn't exist",
			"requester", requester.Username,
			"target", input.Username,
		)
		return nil, ErrUserNotFound
	}
	invitee := users[0]
	if invitee.Id == requester.Id {
		return nil, ErrCantInviteSelf
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
		return nil, ErrSomethingWentWrongInvite
	}
	if len(friendships) > 0 {
		friendshipRecord := friendships[0]
		switch friendshipRecord.Status {
		case domain.FriendsStatus:
			return nil, ErrUserAlreadyFriends
		case domain.FriendshipStatusRequestSent:
			return nil, ErrUserAlreadyInvited
		default:
			slog.Error(
				"this shouldn't be possible. friendship record status is not defined",
				"record", friendshipRecord,
				"requester", requester.Username,
				"invitee", invitee.Username,
			)
			return nil, ErrSomethingWentWrongInvite
		}
	}

	var friendship *domain.Friendship
	err = tx.A(func() error {
		_friendship, err := us.friendships.CreateOne(&pocketbase.CreateOneInput[domain.Friendship]{
			Data: domain.Friendship{
				RequesterId: requester.Id,
				InviteeId:   invitee.Id,
				Status:      domain.FriendshipStatusRequestSent,
			},
		})
		if err == nil {
			friendship = _friendship
		}
		return err
	}, func() error {
		return us.friendships.DeleteOne(&pocketbase.DeleteOneInput{
			Id: friendship.Id,
		})
	})
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occured when creating a friendship record",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		return nil, ErrSomethingWentWrongInvite
	}

	var chat *domain.Chat
	err = tx.A(func() error {
		_chat, err := us.chats.CreateOne(&pocketbase.CreateOneInput[domain.Chat]{
			Data: domain.Chat{
				CreatorId:       requester.Id,
				ParticipantsIds: []string{invitee.Id},
				PublicKey:       input.ChatPublicKey,
			},
		})
		if err == nil {
			chat = _chat
		}
		return err
	}, func() error {
		return us.chats.DeleteOne(&pocketbase.DeleteOneInput{
			Id: chat.Id,
		})
	})
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occured when creating a chat record",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		return nil, ErrSomethingWentWrongInvite
	}

	var keyExchange *domain.KeyExchange
	err = tx.A(func() error {
		_keyExchange, err := us.keyExchanges.CreateOne(&pocketbase.CreateOneInput[domain.KeyExchange]{
			Data: domain.KeyExchange{
				EncryptedPrivateKey: input.EncryptedPrivateKey,
				RequesterId:         requester.Id,
				TargetId:            invitee.Id,
				RelatedChatId:       chat.Id,
			},
		})
		if err == nil {
			keyExchange = _keyExchange
		}
		return err
	}, func() error {
		return us.keyExchanges.DeleteOne(&pocketbase.DeleteOneInput{
			Id: keyExchange.Id,
		})
	})
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		slog.Error(
			"error occured when creating a key exchange record",
			"main reason", reason,
			"rollback errors", rollbackErrors,
			"chat", chat,
			"requester", requester.Username,
			"invitee", invitee.Username,
		)
		return nil, ErrSomethingWentWrongInvite
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

	return chat, nil
}
