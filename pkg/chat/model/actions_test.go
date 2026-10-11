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

func TestRenderActionsJSON(t *testing.T) {
	ra := RenderActions{
		Action: &RenderAction{
			Navigations:  []Navigation{{PushCard: &Card{}}, {PopToRoot: true}, {EndNavigation: &EndNavigation{Action: EndNavigationCloseDialog}}},
			Notification: &Notification{Text: "done"},
		},
	}
	got, err := json.Marshal(ra)
	if err != nil {
		t.Fatalf("marshaling render actions: %v", err)
	}
	const want = `{"action":{"navigations":[{"pushCard":{}},{"popToRoot":true},{"endNavigation":{"action":"CLOSE_DIALOG"}}],"notification":{"text":"done"}}}`
	if string(got) != want {
		t.Errorf("marshaled JSON = %s, want %s", got, want)
	}
}
