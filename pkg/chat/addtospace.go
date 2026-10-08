package chat

import (
	"encoding/json"
	"io"
)

type AddedToSpace struct {
	cloudEvent      cloudEvent[addedToSpaceEvent]
	UserEmail       string
	UserDisplayName string
	SpaceType       SpaceType
}

type RemovedFromSpace struct {
	cloudEvent      cloudEvent[removedFromSpaceEvent]
	UserEmail       string
	UserDisplayName string
	SpaceType       SpaceType
}


func NewResponseMessageWithText(text string) dataActions {
	return dataActions{
		HostAppDataAction: hostAppDataAction{
			ChatDataAction: chatDataAction{
				CreateMessageAction: &createMessageAction{
					Message: message{
						Text: text,
					},
				},
			},
		},
	}
}

func NewAddedToSpaceEvent(body io.Reader) (*AddedToSpace, error) {
	var event AddedToSpace
	err := json.NewDecoder(body).Decode(&event.cloudEvent)
	if err != nil {
		return nil, err
	}
	event.UserEmail = event.cloudEvent.Chat.User.Email
	event.UserDisplayName = event.cloudEvent.Chat.User.DisplayName
	event.SpaceType = event.cloudEvent.Chat.AddedToSpace.Space.Type
	return &event, nil
}

func NewRemovedFromSpaceEvent(body io.Reader) (*RemovedFromSpace, error) {
	var event RemovedFromSpace
	err := json.NewDecoder(body).Decode(&event.cloudEvent)
	if err != nil {
		return nil, err
	}
	event.UserEmail = event.cloudEvent.Chat.User.Email
	event.UserDisplayName = event.cloudEvent.Chat.User.DisplayName
	event.SpaceType = event.cloudEvent.Chat.RemovedFromSpace.Space.Type
	return &event, nil
}
