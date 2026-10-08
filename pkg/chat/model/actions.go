package model

type DataActionsMarkup struct {
	HostAppDataAction *HostAppDataActionMarkup `json:"hostAppDataAction,omitempty"`
}

type HostAppDataActionMarkup struct {
	ChatDataAction *ChatDataActionMarkup `json:"chatDataAction,omitempty"`
}

type ChatDataActionMarkup struct {
	CreateMessageAction *CreateMessageActionMarkup `json:"createMessageAction,omitempty"`
	UpdateMessageAction *UpdateMessageActionMarkup `json:"updateMessageAction,omitempty"`
	DeleteMessageAction *DeleteMessageActionMarkup `json:"deleteMessageAction,omitempty"`
}

type CreateMessageActionMarkup struct {
	Message *Message `json:"message,omitempty"`
}

type UpdateMessageActionMarkup struct {
	Message *Message `json:"message,omitempty"`
}

type DeleteMessageActionMarkup struct {
	Name string `json:"name,omitempty"`
}

type Message struct {
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
}
