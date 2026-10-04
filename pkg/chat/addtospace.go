package chat

import (
	"encoding/json"
	"io"
)

type AddedToSpace struct {
	cloudEvent      cloudEvent
	UserEmail       string
	UserDisplayName string
	SpaceType       SpaceType
}

func NewChatEvent(body io.Reader) (*AddedToSpace, error) {
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
