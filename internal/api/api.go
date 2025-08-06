// Package api provides HTTP client implementation for MAX Bot API.
// It defines API endpoints, request/response types mapping, and HTTP configuration.
package api

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
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

// API server configuration - make them variables for testing
var (
	Scheme = "https"         // Protocol scheme for API requests
	Host   = "botapi.max.ru" // MAX Bot API host
)

// HttpMethod represents HTTP request method types
type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	PATCH  HttpMethod = "PATCH"
)

// API Endpoint constants
const (
	GetBotInfo           core.Endpoint = "GetBotInfo"
	UpdateBotInfo        core.Endpoint = "UpdateBotInfo"
	GetChatsList         core.Endpoint = "GetChatsList"
	GetChatInfoByLink    core.Endpoint = "GetChatInfoByLink"
	GetChatInfoByID      core.Endpoint = "GetChatInfoByID"
	UpdateChat           core.Endpoint = "UpdateChat"
	DeleteChat           core.Endpoint = "DeleteChat"
	SendChatAction       core.Endpoint = "SendChatAction"
	GetChatPinMessage    core.Endpoint = "GetChatPinMessage"
	GetChatPinnedMessage core.Endpoint = "GetChatPinnedMessage"
	DeleteChatPinMessage core.Endpoint = "DeleteChatPinMessage"
	GetChatMyMembership  core.Endpoint = "GetChatMyMembership"
	LeaveChat            core.Endpoint = "LeaveChat"
	GetChatListAdmins    core.Endpoint = "GetChatListAdmins"
	SetChatAdmins        core.Endpoint = "SetChatAdmins"
	RemoveChatAdmin      core.Endpoint = "RemoveChatAdmin"
	GetChatMembers       core.Endpoint = "GetChatMembers"
	AddChatMembers       core.Endpoint = "AddChatMembers"
	RemoveChatMembers    core.Endpoint = "RemoveChatMembers"
	GetSubscription      core.Endpoint = "GetSubscription"
	UpdateSubscription   core.Endpoint = "UpdateSubscription"
	Unsubscribe          core.Endpoint = "Unsubscribe"
	GetSubscribeUpdate   core.Endpoint = "GetSubscribeUpdate"
	GetUploadUrl         core.Endpoint = "GetUploadUrl"
	GetListMsg           core.Endpoint = "GetListMsg"
	SendMsg              core.Endpoint = "SendMsg"
	EditMsg              core.Endpoint = "EditMsg"
	DeleteMsg            core.Endpoint = "DeleteMsg"
	GetMsgByID           core.Endpoint = "GetMsgByID"
	GetVideoInfo         core.Endpoint = "GetVideoInfo"
	AnswerCallback       core.Endpoint = "AnswerCallback"
)

// ContentType represents HTTP Content-Type header values
type ContentType string

const (
	Context           ContentType = "Content-Type"
	MultipartFormData ContentType = "multipart/form-data"
	ApplicationJson   ContentType = "application/json"
)

// EndpointConfig defines the configuration for each API endpoint.
// It includes HTTP method, content type, URL path template, and request/response type information.
type EndpointConfig struct {
	Method        HttpMethod   // HTTP method (GET, POST, etc.)
	ContentType   ContentType  // Content-Type header value
	Path          string       // URL path template with {param} placeholders
	RequestModel  reflect.Type // Expected request struct type
	ResponseModel reflect.Type // Expected response struct type
}

// EndpointConfigs maps each API endpoint to its configuration.
// This is used by the client to determine how to make requests for each endpoint.
var EndpointConfigs = map[core.Endpoint]*EndpointConfig{
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
	UpdateChat: {
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
