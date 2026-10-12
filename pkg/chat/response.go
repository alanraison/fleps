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
