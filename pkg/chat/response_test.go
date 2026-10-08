package chat

import (
	"encoding/json"
	"testing"
)

func TestNewResponseMessageWithText(t *testing.T) {
	got, err := json.Marshal(NewResponseMessageWithText("hello"))
	if err != nil {
		t.Fatalf("marshaling response: %v", err)
	}

	const want = `{"hostAppDataAction":{"chatDataAction":{"createMessageAction":{"message":{"text":"hello"}}}}}`
	if string(got) != want {
		t.Errorf("marshaled response = %s, want %s", got, want)
	}
}
