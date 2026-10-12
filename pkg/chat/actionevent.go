package chat

import (
	"encoding/json"
	"fmt"
	"io"
)

// SubmitPredictionsAction is the value of the actionName parameter sent when a
// player submits their predictions.
const SubmitPredictionsAction = "submitPredictions"

type commonEventObject struct {
	Parameters map[string]string    `json:"parameters"`
	FormInputs map[string]formInput `json:"formInputs"`
}

type formInput struct {
	StringInputs struct {
		Value []string `json:"value"`
	} `json:"stringInputs"`
}

type actionEventEnvelope struct {
	CommonEventObject commonEventObject `json:"commonEventObject"`
	Chat              struct {
		User chatUser `json:"user"`
	} `json:"chat"`
}

// ActionEvent is a card action invoked by a user, such as a button click.
type ActionEvent struct {
	// UserEmail is the email of the user who invoked the action.
	UserEmail  string
	Parameters map[string]string
	// FormInputs holds the first value of each form input in the card.
	FormInputs map[string]string
}

// NewActionEvent decodes a card action event.
func NewActionEvent(body io.Reader) (*ActionEvent, error) {
	var env actionEventEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		return nil, fmt.Errorf("decoding action event: %w", err)
	}
	event := &ActionEvent{
		UserEmail:  env.Chat.User.Email,
		Parameters: env.CommonEventObject.Parameters,
		FormInputs: make(map[string]string, len(env.CommonEventObject.FormInputs)),
	}
	for name, in := range env.CommonEventObject.FormInputs {
		if len(in.StringInputs.Value) > 0 {
			event.FormInputs[name] = in.StringInputs.Value[0]
		}
	}
	return event, nil
}
