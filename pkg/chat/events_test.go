package chat

import (
	"errors"
	"strings"
	"testing"
)

func TestNewAddedToSpaceEvent(t *testing.T) {
	const body = `{"chat":{"user":{"displayName":"Ann","email":"ann@example.com"},"addedToSpacePayload":{"space":{"name":"spaces/1","type":"DM","spaceType":"DIRECT_MESSAGE"}}}}`
	got, err := NewAddedToSpaceEvent(strings.NewReader(body))
	if err != nil {
		t.Fatalf("NewAddedToSpaceEvent: %v", err)
	}
	if got.UserEmail != "ann@example.com" || got.UserDisplayName != "Ann" || got.SpaceType != TypeDirectMessage {
		t.Errorf("unexpected event: %+v", got)
	}
}

func TestNewRemovedFromSpaceEvent(t *testing.T) {
	const body = `{"chat":{"user":{"displayName":"Bob","email":"bob@example.com"},"removedFromSpacePayload":{"space":{"name":"spaces/2","type":"ROOM","spaceType":"SPACE"}}}}`
	got, err := NewRemovedFromSpaceEvent(strings.NewReader(body))
	if err != nil {
		t.Fatalf("NewRemovedFromSpaceEvent: %v", err)
	}
	if got.UserEmail != "bob@example.com" || got.UserDisplayName != "Bob" || got.SpaceType != TypeRoom {
		t.Errorf("unexpected event: %+v", got)
	}
}

func TestNewAppCommandEvent(t *testing.T) {
	const body = `{"chat":{"appCommandPayload":{"appCommandMetadata":{"appCommandId":10,"appCommandType":"SLASH_COMMAND"}}}}`
	got, err := NewAppCommandEvent(strings.NewReader(body))
	if err != nil {
		t.Fatalf("NewAppCommandEvent: %v", err)
	}
	if got.CommandId != MakePredictionsCommand {
		t.Errorf("CommandId = %v, want %v", got.CommandId, MakePredictionsCommand)
	}
}

func TestEventDecodeErrors(t *testing.T) {
	for name, fn := range map[string]func() error{
		"added":   func() error { _, err := NewAddedToSpaceEvent(strings.NewReader("{")); return err },
		"removed": func() error { _, err := NewRemovedFromSpaceEvent(strings.NewReader("{")); return err },
		"command": func() error { _, err := NewAppCommandEvent(strings.NewReader("{")); return err },
	} {
		if err := fn(); err == nil {
			t.Errorf("%s: expected error for malformed JSON", name)
		}
	}
}

func TestNewSpaceEvent(t *testing.T) {
	tests := map[string]struct {
		body string
		want SpaceEventKind
		err  error
	}{
		"added":   {`{"chat":{"user":{"email":"a@b.c"},"addedToSpacePayload":{"space":{"type":"DM"}}}}`, SpaceEventAdded, nil},
		"removed": {`{"chat":{"user":{"email":"a@b.c"},"removedFromSpacePayload":{"space":{"type":"DM"}}}}`, SpaceEventRemoved, nil},
		"other":   {`{"chat":{"appCommandPayload":{}}}`, 0, ErrUnknownSpaceEvent},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewSpaceEvent(strings.NewReader(tc.body))
			if !errors.Is(err, tc.err) {
				t.Fatalf("err = %v, want %v", err, tc.err)
			}
			if err == nil && (got.Kind != tc.want || got.UserEmail != "a@b.c" || got.SpaceType != TypeDirectMessage) {
				t.Errorf("unexpected event: %+v", got)
			}
		})
	}
	if _, err := NewSpaceEvent(strings.NewReader("{")); err == nil {
		t.Error("expected error for malformed JSON")
	}
}
