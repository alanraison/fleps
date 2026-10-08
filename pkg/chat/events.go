package chat

// cloudEvent is the envelope Google Chat posts to an HTTP add-on. T is the
// event-specific payload found under the "chat" key.
type cloudEvent[T any] struct {
	Chat T `json:"chat"`
}

type chatUser struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

// SpaceType is the (deprecated) Space.type field of a Chat space.
type SpaceType string

// SpaceKind is the Space.spaceType field of a Chat space.
type SpaceKind string

const (
	TypeDirectMessage SpaceType = "DM"
	TypeRoom          SpaceType = "ROOM"

	SpaceTypeDirectMessage SpaceKind = "DIRECT_MESSAGE"
	SpaceTypeSpace         SpaceKind = "SPACE"
)

type chatSpace struct {
	Name      string    `json:"name"`
	Type      SpaceType `json:"type"`
	SpaceType SpaceKind `json:"spaceType"`
}

type chatMessage struct {
	Sender        chatUser  `json:"sender"`
	Space         chatSpace `json:"space"`
	Text          string    `json:"text"`
	FormattedText string    `json:"formattedText"`
}

type spacePayload struct {
	Space chatSpace `json:"space"`
}

type addedToSpaceEvent struct {
	User         chatUser     `json:"user"`
	AddedToSpace spacePayload `json:"addedToSpacePayload"`
}

type removedFromSpaceEvent struct {
	User             chatUser     `json:"user"`
	RemovedFromSpace spacePayload `json:"removedFromSpacePayload"`
}

type appCommandMetadata struct {
	AppCommandId   float64 `json:"appCommandId"`
	AppCommandType string  `json:"appCommandType"`
}

type appCommandPayload struct {
	Space    chatSpace          `json:"space"`
	Message  chatMessage        `json:"message"`
	Metadata appCommandMetadata `json:"appCommandMetadata"`
}

type appCommandEvent struct {
	User       chatUser          `json:"user"`
	AppCommand appCommandPayload `json:"appCommandPayload"`
}
