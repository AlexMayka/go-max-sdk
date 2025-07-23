package config

const (
	Scheme = "https"
	Host   = "botapi.max.ru"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	PATCH  HttpMethod = "PATCH"
)

type Endpoint string

const (
	GetBotInfo           Endpoint = "GetBotInfo"
	UpdateBotInfo        Endpoint = "UpdateBotInfo"
	GetChatsList         Endpoint = "GetChatsList"
	GetChatInfoByLink    Endpoint = "GetChatInfoByLink"
	GetChatInfoByID      Endpoint = "GetChatInfoByID"
	UpdateChart          Endpoint = "UpdateChart"
	DeleteChat           Endpoint = "DeleteChat"
	SendChatAction       Endpoint = "SendChatAction"
	GetChatPinMessage    Endpoint = "GetChatPinMessage"
	GetChatPinnedMessage Endpoint = "GetChatPinnedMessage"
	DeleteChatPinMessage Endpoint = "DeleteChatPinMessage"
	GetChatMyMembership  Endpoint = "GetChatMyMembership"
	LeaveChat            Endpoint = "LeaveChat"
	GetChatListAdmins    Endpoint = "GetChatListAdmins"
	SetChatAdmins        Endpoint = "SetChatAdmins"
	RemoveChatAdmin      Endpoint = "RemoveChatAdmin"
	GetChatMembers       Endpoint = "GetChatMembers"
	AddChatMembers       Endpoint = "AddChatMembers"
	RemoveChatMembers    Endpoint = "RemoveChatMembers"
	GetSubscription      Endpoint = "GetSubscription"
	UpdateSubscription   Endpoint = "UpdateSubscription"
	Unsubscribe          Endpoint = "Unsubscribe"
	GetSubscribeUpdate   Endpoint = "GetSubscribeUpdate"
	GetUploadUrl         Endpoint = "GetUploadUrl"
	GetListMsg           Endpoint = "GetListMsg"
	SendMsg              Endpoint = "SendMsg"
	EditMsg              Endpoint = "EditMsg"
	DeleteMsg            Endpoint = "DeleteMsg"
	GetMsgByID           Endpoint = "GetMsgByID"
	GetVideoInfo         Endpoint = "GetVideoInfo"
	AnswerCallback       Endpoint = "AnswerCallback"
)

type ContentType string

const (
	MultipartFormData ContentType = "multipart/form-data"
	ApplicationJson   ContentType = "application/json"
)
