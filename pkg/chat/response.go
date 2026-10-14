package chat

import chatmodel "github.com/alanraison/fleps/pkg/chat/model"

func NewResponseMessageWithText(text string) chatmodel.DataActionsMarkup {
	return chatmodel.DataActionsMarkup{
		HostAppDataAction: &chatmodel.HostAppDataActionMarkup{
			ChatDataAction: &chatmodel.ChatDataActionMarkup{
				CreateMessageAction: &chatmodel.CreateMessageActionMarkup{
					Message: &chatmodel.Message{
						Text: text,
					},
				},
			},
		},
	}
}

// NewErrorNotification returns a response which makes Google Chat show text to
// the user as a notification. It should be sent with a 200 status, as Chat
// only displays the content of successful responses.
func NewErrorNotification(text string) chatmodel.RenderActions {
	return chatmodel.RenderActions{
		Action: &chatmodel.RenderAction{
			Notification: &chatmodel.Notification{Text: text},
		},
	}
}

// NewResponseMessageWithImage returns a response which posts a card containing
// the image at imageURL.
func NewResponseMessageWithImage(title, imageURL, altText string) chatmodel.DataActionsMarkup {
	return chatmodel.DataActionsMarkup{
		HostAppDataAction: &chatmodel.HostAppDataActionMarkup{
			ChatDataAction: &chatmodel.ChatDataActionMarkup{
				CreateMessageAction: &chatmodel.CreateMessageActionMarkup{
					Message: &chatmodel.Message{
						CardsV2: []chatmodel.CardWithID{{
							CardID: "results",
							Card: &chatmodel.Card{
								Header: &chatmodel.CardHeader{Title: title},
								Sections: []chatmodel.Section{{
									Widgets: []chatmodel.Widget{{
										Image: &chatmodel.Image{ImageURL: imageURL, AltText: altText},
									}},
								}},
							},
						}},
					},
				},
			},
		},
	}
}
