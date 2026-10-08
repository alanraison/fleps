package chat

import (
	"encoding/json"
	"io"
)

const (
	MakePredictionsCommand = 10.0
)

type AppCommand struct {
	cloudEvent cloudEvent[appCommandEvent]
	CommandId  float64
}

func NewAppCommandEvent(body io.Reader) (*AppCommand, error) {
	var event AppCommand
	err := json.NewDecoder(body).Decode(&event.cloudEvent)
	if err != nil {
		return nil, err
	}
	event.CommandId = event.cloudEvent.Chat.AppCommand.Metadata.AppCommandId
	return &event, nil
}
