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
