package userservice

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gutter/internal/domain"

	"github.com/tomek7667/cyberman-go/expo"
	"github.com/tomek7667/cyberman-go/pocketbase"
)

type AnswerDto struct {
	FriendshipId string `json:"friendshipId" binding:"required"`
	Accept       bool   `json:"accept"`
}

var (
	ErrFriendshipNotFound   = errors.New("friend request not found")
	ErrCantAcceptOwnRequest = errors.New("you can't accept your own friend request")
)

func (us *Client) Answer(requester *domain.User, input *AnswerDto) error {
	// TODO: Cleanup encrypted private key after accept - it's so > 9k characters combined damn
	var err, notificationErr error
	friendships, err := us.friendships.GetFullList(&pocketbase.GetFullListInput[domain.Friendship]{
		Filter: pocketbase.BuildFilter(
			"id = {:friendshipId}",
			map[string]interface{}{
				"friendshipId": input.FriendshipId,
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
	if foundFriendship.InviteeId != requester.Id && foundFriendship.RequesterId != requester.Id {
		return ErrFriendshipNotFound
	}
	if foundFriendship.RequesterId == requester.Id && input.Accept {
		return ErrCantAcceptOwnRequest
	}
	if input.Accept {
		foundFriendship.Status = domain.FriendsStatus
		_, err = us.friendships.SaveOne(&pocketbase.SaveOneInput[domain.Friendship]{
			Data: foundFriendship,
		})
		notificationErr = us.expoClient.SendNotification(&expo.SendNotificationInput{
			Title: "Friendship request accepted",
			Body: fmt.Sprintf(
				"%s accepted your invitation.",
				foundFriendship.E.Invitee.Username,
			),
		}, foundFriendship.E.Requester.PushToken)
	} else {
		err = us.friendships.DeleteOne(&pocketbase.DeleteOneInput{
			Id: foundFriendship.Id,
		})
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
	if notificationErr != nil && !strings.Contains(notificationErr.Error(), "Token should start with ExponentPushToken") {
		slog.Error(
			"notification error occurred",
			"err", notificationErr,
			"user invitee", requester.Username,
			"user requester", foundFriendship.E.Requester.Username,
			"expo token", foundFriendship.E.Requester.PushToken,
		)
	}
	if err != nil {
		slog.Error(
			"saveOne / deleteOne failed",
			"err", err,
			"user requester", requester.Username,
			"found friendship", foundFriendship,
			"input", input,
		)
		return err
	}
	return nil
}
