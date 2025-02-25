package userservice

import (
	"errors"
	"fmt"
	"gutter/internal/domain"
	"log/slog"
	"strings"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
	"github.com/tomek7667/cyberman-go/utils"
)

type AnswerDto struct {
	FriendshipId string `json:"friendshipId" binding:"required"`
	Accept       bool   `json:"accept"`
}

var (
	ErrFriendshipNotFound       = errors.New("friend request not found")
	ErrCantAcceptOwnRequest     = errors.New("you can't accept your own friend request")
	ErrSomethingWentWrongAnswer = errors.New("something went wrong while answering an invite")
)

func (us *Client) Answer(requester *domain.User, input *AnswerDto) error {
	var err, notificationErr error
	friendships, err := us.friendships.GetFullList(&pocketbase.GetFullListInput[domain.Friendship]{
		Filter: pocketbase.BuildFilter(
			"id = {:friendshipId} && (requester.id = {:userId} || invitee.id = {:userId})",
			map[string]interface{}{
				"friendshipId": input.FriendshipId,
				"userId":       requester.Id,
			},
		),
		Expand: "requester,invitee",
	})
	if err != nil {
		return err
	}
	if len(friendships) == 0 {
		return ErrFriendshipNotFound
	}
	foundFriendship := friendships[0]
	if foundFriendship.RequesterId == requester.Id && input.Accept {
		return ErrCantAcceptOwnRequest
	}
	tx := utils.CreateTransaction()
	if input.Accept {
		// user has accepted the friendship - updating the friendship status to friends
		previousStatus := foundFriendship.Status
		err = tx.A(func() error {
			foundFriendship.Status = domain.FriendsStatus
			_, err = us.friendships.SaveOne(&pocketbase.SaveOneInput[domain.Friendship]{
				Data: foundFriendship,
			})
			return err
		}, func() error {
			foundFriendship.Status = previousStatus
			_, err = us.friendships.SaveOne(&pocketbase.SaveOneInput[domain.Friendship]{
				Data: foundFriendship,
			})
			return err
		})
		if err != nil {
			reason, rollbackErrors := tx.R(err)
			slog.Error(
				"error occured when saving a friendship record with friends status",
				"main reason", reason,
				"rollback errors", rollbackErrors,
				"requester", requester.Username,
			)
			return ErrSomethingWentWrongAnswer
		}

		// inform opponent user about the friendship
		notificationErr = us.expoClient.SendNotification(&expo.SendNotificationInput{
			Title: "Friendship request accepted",
			Body: fmt.Sprintf(
				"%s accepted your invitation.",
				foundFriendship.E.Invitee.Username,
			),
		}, foundFriendship.E.Requester.PushToken)

		// add partipant id to the chat
		kRecords, err := us.keyExchanges.GetFullList(&pocketbase.GetFullListInput[domain.KeyExchange]{
			Filter: pocketbase.BuildFilter(
				"friendship.id = {:friendshipId}",
				map[string]interface{}{
					"friendshipId": foundFriendship.Id,
				},
			),
			Expand: "requester,target,relatedChat,friendship",
		})
		if err != nil {
			reason, rollbackErrors := tx.R(err)
			slog.Error(
				"error occured when getting key exchanges records 1",
				"main reason", reason,
				"rollback errors", rollbackErrors,
				"requester", requester.Username,
			)
			return ErrSomethingWentWrongAnswer
		}
		if len(kRecords) == 0 {
			reason, rollbackErrors := tx.R(errors.New("no key exchanges records found"))
			slog.Error(
				"error occured when getting key exchanges records 2",
				"main reason", reason,
				"rollback errors", rollbackErrors,
				"requester", requester.Username,
			)
			return ErrSomethingWentWrongAnswer
		}
		keyExchange := kRecords[0]
		relatedChat := keyExchange.E.RelatedChat

		// cleanup private key from the keyexchange - it's very long
		previousEPK := keyExchange.EncryptedPrivateKey
		err = tx.A(func() error {
			keyExchange.EncryptedPrivateKey = ""
			_, err = us.keyExchanges.SaveOne(&pocketbase.SaveOneInput[domain.KeyExchange]{
				Data: keyExchange,
			})
			return err
		}, func() error {
			keyExchange.EncryptedPrivateKey = previousEPK
			_, err = us.keyExchanges.SaveOne(&pocketbase.SaveOneInput[domain.KeyExchange]{
				Data: keyExchange,
			})
			return err
		})
		if err != nil {
			reason, rollbackErrors := tx.R(err)
			slog.Error(
				"error occured when deleting encrypted private key from the keyexchange",
				"main reason", reason,
				"rollback errors", rollbackErrors,
				"requester", requester.Username,
			)
			return ErrSomethingWentWrongAnswer
		}

		// add participant to the related chat
		var previousParticipantsIds []string
		copy(previousParticipantsIds, relatedChat.ParticipantsIds)
		tx.A(func() error {
			relatedChat.ParticipantsIds = append(
				relatedChat.ParticipantsIds,
				requester.Id,
			)
			_, err = us.chats.SaveOne(&pocketbase.SaveOneInput[domain.Chat]{
				Data: relatedChat,
			})
			return err
		}, func() error {
			relatedChat.ParticipantsIds = previousParticipantsIds
			_, err = us.chats.SaveOne(&pocketbase.SaveOneInput[domain.Chat]{
				Data: relatedChat,
			})
			return err
		})
		if err != nil {
			reason, rollbackErrors := tx.R(err)
			slog.Error(
				"error occured when saving new chat participant ids",
				"main reason", reason,
				"rollback errors", rollbackErrors,
				"requester", requester.Username,
			)
			return ErrSomethingWentWrongAnswer
		}
	} else {
		// user has cancelled the friendship - deleting the record
		err = tx.A(func() error {
			return us.friendships.DeleteOne(&pocketbase.DeleteOneInput{
				Id: foundFriendship.Id,
			})
		}, func() error {
			_, err := us.friendships.CreateOne(&pocketbase.CreateOneInput[domain.Friendship]{
				Data: foundFriendship,
			})
			return err
		})
		if err != nil {
			reason, rollbackErrors := tx.R(err)
			slog.Error(
				"error occured when deleting a friendship record with friends status",
				"main reason", reason,
				"rollback errors", rollbackErrors,
				"requester", requester.Username,
			)
			return ErrSomethingWentWrongAnswer
		}

		// inform opponent user about no longer friends
		if requester.Id == foundFriendship.InviteeId {
			notificationErr = us.expoClient.SendNotification(&expo.SendNotificationInput{
				Title: "Friendship request rejected",
				Body: fmt.Sprintf(
					"%s rejected the friendship.",
					foundFriendship.E.Invitee.Username,
				),
			}, foundFriendship.E.Requester.PushToken)
		} else {
			notificationErr = us.expoClient.SendNotification(&expo.SendNotificationInput{
				Title: "Friendship request rejected",
				Body: fmt.Sprintf(
					"%s rejected the friendship.",
					foundFriendship.E.Requester.Username,
				),
			}, foundFriendship.E.Invitee.PushToken)
		}
	}

	// if there was any notification error that is not related to the invalid token log it
	if notificationErr != nil && !strings.Contains(notificationErr.Error(), "Token should start with ExponentPushToken") {
		slog.Error(
			"notification error occurred",
			"err", notificationErr,
			"user invitee", requester.Username,
			"user requester", foundFriendship.E.Requester.Username,
			"expo token", foundFriendship.E.Requester.PushToken,
		)
	}
	return nil
}
