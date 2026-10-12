package chat

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	MakePredictionsCommand = 10.0
)

type AppCommand struct {
	cloudEvent cloudEvent[appCommandEvent]
	CommandId  float64
	// UserEmail is the email of the user who invoked the command.
	UserEmail string
}

func NewAppCommandEvent(body io.Reader) (*AppCommand, error) {
	var event AppCommand
	err := json.NewDecoder(body).Decode(&event.cloudEvent)
	if err != nil {
		return nil, fmt.Errorf("decoding app command event: %w", err)
	}
	event.CommandId = event.cloudEvent.Chat.AppCommand.Metadata.AppCommandId
	event.UserEmail = event.cloudEvent.Chat.User.Email
	return &event, nil
}
