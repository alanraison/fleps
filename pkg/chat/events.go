package chat

type cloudEvent struct {
	Chat chatEvent `json:"chat"`
}

type chatEvent struct {
	User         chatUser             `json:"user"`
	AddedToSpace *addedToSpacePayload `json:"addedToSpacePayload"`
	Message      *messagePayload      `json:"messagePayload"`
}

type chatUser struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type SpaceType string
type SpaceSpaceType string

const (
	TypeDirectMessage      SpaceType      = "DM"
	TypeRoom               SpaceType      = "ROOM"
	SpaceTypeDirectMessage SpaceSpaceType = "DIRECT_MESSAGE"
	SpaceTypeSpace         SpaceSpaceType = "SPACE"
)

type chatSpace struct {
	Name      string         `json:"name"`
	Type      SpaceType      `json:"type"`
	SpaceType SpaceSpaceType `json:"spaceType"`
}

type chatMessage struct {
	Sender        chatUser  `json:"sender"`
	Space         chatSpace `json:"space"`
	Text          string    `json:"text"`
	FormattedText string    `json:"formattedText"`
}

type addedToSpacePayload struct {
	Space chatSpace `json:"space"`
}

type messagePayload struct {
	Space   chatSpace   `json:"space"`
	Message chatMessage `json:"message"`
}
