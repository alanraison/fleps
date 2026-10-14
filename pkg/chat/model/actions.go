package model

import "encoding/json"

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
	Name    string       `json:"name,omitempty"`
	Text    string       `json:"text,omitempty"`
	CardsV2 []CardWithID `json:"cardsV2,omitempty"`
}

// CardWithID is a card in a message, identified by a unique ID.
type CardWithID struct {
	CardID string `json:"cardId"`
	Card   *Card  `json:"card"`
}

// RenderActions is a set of render instructions that tells a host application
// how to render a card or perform an action in response to a user interaction.
type RenderActions struct {
	Action        *RenderAction       `json:"action,omitempty"`
	HostAppAction HostAppActionMarkup `json:"hostAppAction,omitempty"`
	Schema        string              `json:"schema,omitempty"`
}

// HostAppActionMarkup holds actions handled by individual host apps, keyed by
// host (for example "gmailAction"). The contents are passed through untouched.
type HostAppActionMarkup map[string]json.RawMessage

// RenderAction is the action that add-ons can use to update the UI. In Google
// Chat, dialogs are opened, updated and closed through Navigations.
type RenderAction struct {
	Navigations      []Navigation  `json:"navigations,omitempty"`
	Link             *OpenLink     `json:"link,omitempty"`
	Notification     *Notification `json:"notification,omitempty"`
	LinkPreview      *LinkPreview  `json:"linkPreview,omitempty"`
	ModifyOperations []ModifyCard  `json:"modifyOperations,omitempty"`
}

// Navigation updates or navigates between cards in a card stack. Only one
// field may be set. In Chat, PushCard and UpdateCard open or update a dialog
// and EndNavigation closes it.
type Navigation struct {
	PopToRoot     bool           `json:"popToRoot,omitempty"`
	Pop           bool           `json:"pop,omitempty"`
	PopToCard     string         `json:"popToCard,omitempty"`
	PushCard      *Card          `json:"pushCard,omitempty"`
	UpdateCard    *Card          `json:"updateCard,omitempty"`
	EndNavigation *EndNavigation `json:"endNavigation,omitempty"`
}

// EndNavigationAction is the action taken to close a dialog.
type EndNavigationAction string

const (
	EndNavigationActionUnspecified EndNavigationAction = "ACTION_UNSPECIFIED"
	// EndNavigationCloseDialog closes a dialog.
	EndNavigationCloseDialog EndNavigationAction = "CLOSE_DIALOG"
	// EndNavigationCloseDialogAndExecute closes a dialog and refreshes the card
	// that opened it.
	EndNavigationCloseDialogAndExecute EndNavigationAction = "CLOSE_DIALOG_AND_EXECUTE"
)

// EndNavigation closes a dialog in Google Chat.
type EndNavigation struct {
	Action EndNavigationAction `json:"action,omitempty"`
}

// Notification displays a notification in the host application, or in Chat
// when a dialog is submitted and closed.
type Notification struct {
	Text string `json:"text,omitempty"`
}

// LinkPreview previews a third-party link with a card and smart chip.
type LinkPreview struct {
	PreviewCard      *Card  `json:"previewCard,omitempty"`
	Title            string `json:"title,omitempty"`
	LinkPreviewTitle string `json:"linkPreviewTitle,omitempty"`
}

// ModifyCard updates a card based on a user interaction. Only one operation
// may be set.
type ModifyCard struct {
	UpdateWidget *UpdateWidget `json:"updateWidget,omitempty"`
}

// UpdateWidget updates a widget in a card or dialog, for example to provide
// multiselect suggestions.
type UpdateWidget struct {
	SelectionInputWidgetSuggestions *SelectionInputWidgetSuggestions `json:"selectionInputWidgetSuggestions,omitempty"`
}

// SelectionInputWidgetSuggestions returns selection items from an external
// data source for a multiselect menu.
type SelectionInputWidgetSuggestions struct {
	Suggestions []SelectionItem `json:"suggestions,omitempty"`
}
