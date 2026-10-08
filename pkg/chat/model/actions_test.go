package model

import (
	"encoding/json"
	"testing"
)

func TestChatDataActionMarkupJSON(t *testing.T) {
	action := DataActionsMarkup{
		HostAppDataAction: &HostAppDataActionMarkup{
			ChatDataAction: &ChatDataActionMarkup{
				CreateMessageAction: &CreateMessageActionMarkup{
					Message: &Message{Text: "hello"},
				},
				UpdateMessageAction: &UpdateMessageActionMarkup{
					Message: &Message{Name: "spaces/space/messages/message", Text: "updated"},
				},
				DeleteMessageAction: &DeleteMessageActionMarkup{
					Name: "spaces/space/messages/message",
				},
			},
		},
	}

	got, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshaling data actions: %v", err)
	}

	const want = `{"hostAppDataAction":{"chatDataAction":{"createMessageAction":{"message":{"text":"hello"}},"updateMessageAction":{"message":{"name":"spaces/space/messages/message","text":"updated"}},"deleteMessageAction":{"name":"spaces/space/messages/message"}}}}`
	if string(got) != want {
		t.Errorf("marshaled JSON = %s, want %s", got, want)
	}
}
