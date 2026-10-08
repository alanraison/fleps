package chat

import chatmodel "github.com/alanraison/fleps/pkg/chat/model"

type dataActions struct {
	HostAppDataAction hostAppDataAction `json:"hostAppDataAction"`
}

type hostAppDataAction struct {
	ChatDataAction chatDataAction `json:"chatDataAction"`
}

type chatDataAction struct {
	CreateMessageAction *createMessageAction `json:"createMessageAction"`
	UpdateMessageAction *updateMessageAction `json:"updateMessageAction"`
}

type createMessageAction struct {
	Message message `json:"message"`
}

type updateMessageAction struct {
	Message message `json:"message"`
}

type message struct {
	Text string `json:"text"`
}

type renderAction struct {
	Action        action `json:"action"`
	HostAppAction any    `json:"hostAppAction"`
}

type action struct {
	Navigations      []navigationAction `json:"navigations"`
	Link             any                `json:"link"`
	Notification     any                `json:"notification"`
	ModifyOperations any                `json:"modifyOperations"`
}

type navigationAction interface {
	NavigationAction()
}

type popToRootNavigationAction struct {
	PopToRoot bool `json:"popToRoot"`
}

func (p popToRootNavigationAction) NavigationAction() {}

type popNavigationAction struct {
	Pop bool `json:"pop"`
}

func (p popNavigationAction) NavigationAction() {}

type popToCardNavigationAction struct {
	PopToCard string `json:"popToCard"`
}

func (p popToCardNavigationAction) NavigationAction() {}

type pushCardNavigationAction struct {
	PushCard chatmodel.Card `json:"pushCard"`
}

func (p pushCardNavigationAction) NavigationAction() {}

type updateCardNavigationAction struct {
	UpdateCard chatmodel.Card `json:"updateCard"`
}

func (u updateCardNavigationAction) NavigationAction() {}

type endNavigationNavigationAction struct {
	EndNavigation any `json:"endNavigation"`
}

func (e endNavigationNavigationAction) NavigationAction() {}
