package chat

import (
	"encoding/json"
	"errors"
	"fmt"
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

func NewAddedToSpaceEvent(body io.Reader) (*AddedToSpace, error) {
	var event AddedToSpace
	err := json.NewDecoder(body).Decode(&event.cloudEvent)
	if err != nil {
		return nil, fmt.Errorf("decoding added-to-space event: %w", err)
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
		return nil, fmt.Errorf("decoding removed-from-space event: %w", err)
	}
	event.UserEmail = event.cloudEvent.Chat.User.Email
	event.UserDisplayName = event.cloudEvent.Chat.User.DisplayName
	event.SpaceType = event.cloudEvent.Chat.RemovedFromSpace.Space.Type
	return &event, nil
}

// SpaceEventKind identifies which space lifecycle event a payload carries.
type SpaceEventKind int

const (
	SpaceEventAdded SpaceEventKind = iota + 1
	SpaceEventRemoved
)

// SpaceEvent is an added-to-space or removed-from-space event, as determined
// from the payload.
type SpaceEvent struct {
	Kind            SpaceEventKind
	UserEmail       string
	UserDisplayName string
	SpaceType       SpaceType
}

// ErrUnknownSpaceEvent is returned when a payload is neither an
// added-to-space nor a removed-from-space event.
var ErrUnknownSpaceEvent = errors.New("payload is not a space event")

type spaceEventEnvelope struct {
	User             chatUser      `json:"user"`
	AddedToSpace     *spacePayload `json:"addedToSpacePayload"`
	RemovedFromSpace *spacePayload `json:"removedFromSpacePayload"`
}

// NewSpaceEvent decodes a space event and works out from the payload whether
// the app was added to or removed from the space.
func NewSpaceEvent(body io.Reader) (*SpaceEvent, error) {
	var ce cloudEvent[spaceEventEnvelope]
	if err := json.NewDecoder(body).Decode(&ce); err != nil {
		return nil, fmt.Errorf("decoding space event: %w", err)
	}
	event := &SpaceEvent{
		UserEmail:       ce.Chat.User.Email,
		UserDisplayName: ce.Chat.User.DisplayName,
	}
	switch {
	case ce.Chat.AddedToSpace != nil:
		event.Kind = SpaceEventAdded
		event.SpaceType = ce.Chat.AddedToSpace.Space.Type
	case ce.Chat.RemovedFromSpace != nil:
		event.Kind = SpaceEventRemoved
		event.SpaceType = ce.Chat.RemovedFromSpace.Space.Type
	default:
		return nil, ErrUnknownSpaceEvent
	}
	return event, nil
}
