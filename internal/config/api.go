package config

import (
	reqBots "github.com/AlexMayka/go-max-sdk/types/requests/bots"
	reqChats "github.com/AlexMayka/go-max-sdk/types/requests/chats"
	reqMsg "github.com/AlexMayka/go-max-sdk/types/requests/messages"
	reqSub "github.com/AlexMayka/go-max-sdk/types/requests/subscriptions"
	reqUp "github.com/AlexMayka/go-max-sdk/types/requests/uploads"

	resBots "github.com/AlexMayka/go-max-sdk/types/responses/bots"
	resChats "github.com/AlexMayka/go-max-sdk/types/responses/chats"
	resMsg "github.com/AlexMayka/go-max-sdk/types/responses/messages"
	resSub "github.com/AlexMayka/go-max-sdk/types/responses/subscriptions"
	resUp "github.com/AlexMayka/go-max-sdk/types/responses/uploads"

	"reflect"
)

type EndpointConfig struct {
	Method        HttpMethod
	ContentType   ContentType
	Path          string
	RequestModel  reflect.Type
	ResponseModel reflect.Type
}

var EndpointConfigs = map[Endpoint]*EndpointConfig{
	GetBotInfo: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/bots/me",
		RequestModel:  reflect.TypeOf((*reqBots.GetInfo)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resBots.GetInfo)(nil)).Elem(),
	},
	UpdateBotInfo: {
		Method:        PATCH,
		ContentType:   ApplicationJson,
		Path:          "/bots/me",
		RequestModel:  reflect.TypeOf((*reqBots.UpdateInfo)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resBots.UpdateInfo)(nil)).Elem(),
	},
	GetChatsList: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats",
		RequestModel:  reflect.TypeOf((*reqChats.GetList)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetList)(nil)).Elem(),
	},
	GetChatInfoByLink: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatLink}",
		RequestModel:  reflect.TypeOf((*reqChats.GetByLink)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetByLink)(nil)).Elem(),
	},
	GetChatInfoByID: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}",
		RequestModel:  reflect.TypeOf((*resChats.GetInfo)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetInfo)(nil)).Elem(),
	},
	UpdateChart: {
		Method:        PATCH,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}",
		RequestModel:  reflect.TypeOf((*reqChats.Update)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.Update)(nil)).Elem(),
	},
	DeleteChat: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}",
		RequestModel:  reflect.TypeOf((*reqChats.Delete)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.Delete)(nil)).Elem(),
	},
	SendChatAction: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/actions",
		RequestModel:  reflect.TypeOf((*reqChats.SendAction)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.SendAction)(nil)).Elem(),
	},
	GetChatPinMessage: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/pin",
		RequestModel:  reflect.TypeOf((*reqChats.GetPinnedMessage)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetPinnedMessage)(nil)).Elem(),
	},
	GetChatPinnedMessage: {
		Method:        PUT,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/pin",
		RequestModel:  reflect.TypeOf((*reqChats.PinMessage)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.PinMessage)(nil)).Elem(),
	},
	DeleteChatPinMessage: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/pin",
		RequestModel:  reflect.TypeOf((*reqChats.UnpinMessage)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.UnpinMessage)(nil)).Elem(),
	},
	GetChatMyMembership: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members/me",
		RequestModel:  reflect.TypeOf((*reqChats.GetMyMembership)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetMyMembership)(nil)).Elem(),
	},
	LeaveChat: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members/me",
		RequestModel:  reflect.TypeOf((*reqChats.LeaveChat)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.LeaveChat)(nil)).Elem(),
	},
	GetChatListAdmins: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members/admins",
		RequestModel:  reflect.TypeOf((*reqChats.GetAdmins)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetAdmins)(nil)).Elem(),
	},
	SetChatAdmins: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members/admins",
		RequestModel:  reflect.TypeOf((*reqChats.SetAdmins)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.SetAdmins)(nil)).Elem(),
	},
	RemoveChatAdmin: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members/admins/{userId}",
		RequestModel:  reflect.TypeOf((*reqChats.RemoveAdmin)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.RemoveAdmin)(nil)).Elem(),
	},
	GetChatMembers: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members",
		RequestModel:  reflect.TypeOf((*reqChats.GetMembers)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.GetMembers)(nil)).Elem(),
	},
	AddChatMembers: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members",
		RequestModel:  reflect.TypeOf((*reqChats.AddMembers)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.AddMembers)(nil)).Elem(),
	},
	RemoveChatMembers: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/chats/{chatId}/members",
		RequestModel:  reflect.TypeOf((*reqChats.RemoveMember)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resChats.RemoveMember)(nil)).Elem(),
	},
	GetSubscription: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/subscriptions",
		RequestModel:  reflect.TypeOf((*reqSub.GetList)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resSub.GetList)(nil)).Elem(),
	},
	UpdateSubscription: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/subscriptions",
		RequestModel:  reflect.TypeOf((*reqSub.Subscribe)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resSub.Subscribe)(nil)).Elem(),
	},
	Unsubscribe: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/subscriptions",
		RequestModel:  reflect.TypeOf((*reqSub.Unsubscribe)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resSub.Unsubscribe)(nil)).Elem(),
	},
	GetSubscribeUpdate: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/updates",
		RequestModel:  reflect.TypeOf((*resSub.GetUpdates)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resSub.GetUpdates)(nil)).Elem(),
	},
	GetUploadUrl: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/uploads",
		RequestModel:  reflect.TypeOf((*reqUp.GetUploadURL)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resUp.GetUploadURL)(nil)).Elem(),
	},
	GetListMsg: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/messages",
		RequestModel:  reflect.TypeOf((*reqMsg.GetList)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.GetList)(nil)).Elem(),
	},
	SendMsg: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/messages",
		RequestModel:  reflect.TypeOf((*reqMsg.Send)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.Send)(nil)).Elem(),
	},
	EditMsg: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/messages",
		RequestModel:  reflect.TypeOf((*reqMsg.Edit)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.Edit)(nil)).Elem(),
	},
	DeleteMsg: {
		Method:        DELETE,
		ContentType:   ApplicationJson,
		Path:          "/messages",
		RequestModel:  reflect.TypeOf((*reqMsg.Delete)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.Delete)(nil)).Elem(),
	},
	GetMsgByID: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/messages/{messageId}",
		RequestModel:  reflect.TypeOf((*reqMsg.GetByID)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.GetByID)(nil)).Elem(),
	},
	GetVideoInfo: {
		Method:        GET,
		ContentType:   ApplicationJson,
		Path:          "/videos/{videoToken}",
		RequestModel:  reflect.TypeOf((*reqMsg.GetVideoInfo)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.GetVideoInfo)(nil)).Elem(),
	},
	AnswerCallback: {
		Method:        POST,
		ContentType:   ApplicationJson,
		Path:          "/answers",
		RequestModel:  reflect.TypeOf((*reqMsg.AnswerCallback)(nil)).Elem(),
		ResponseModel: reflect.TypeOf((*resMsg.AnswerCallback)(nil)).Elem(),
	},
}
