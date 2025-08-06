package core

import (
	"context"
	"github.com/AlexMayka/go-max-sdk/types/models"
	reqMsg "github.com/AlexMayka/go-max-sdk/types/requests/messages"
)

// BotContext provides access to the current update and bot functionality within handlers
type BotContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	UserID    int64
	ChatID    int64
	FirstName string
	LastName  string
	Username  string

	Update *models.Update
	Client APIClient
	FSM    FSM
	Text   string
}

// Reply sends a text message back to the user/chat from the current update
func (ctx *BotContext) Reply(text string) {
	body := &models.NewMessageBody{Text: &text}
	_ = ctx.sendMessage(body)
}

// ReplyWithKeyboard sends a text message with inline keyboard back to the user/chat
func (ctx *BotContext) ReplyWithKeyboard(text string, keyboard *models.InlineKeyboard) error {
	attachments := []models.AttachmentRequest{
		{
			Type: models.AttachmentTypeInlineKeyboard,
			Payload: &models.InlineKeyboardAttachmentRequestPayload{
				Buttons: keyboard.Buttons,
			},
		},
	}

	body := &models.NewMessageBody{
		Text:        &text,
		Attachments: attachments,
	}

	return ctx.sendMessage(body)
}

// EditMessage updates the text of an existing message
func (ctx *BotContext) EditMessage(messageID, newText string) error {
	request := reqMsg.Edit{
		MessageID: messageID,
		NewMessageBody: models.NewMessageBody{
			Text: &newText,
		},
	}

	_, err := ctx.Client.Call(ctx.Ctx, EditMsg, request)
	return err
}

// DeleteMessage removes a message by its ID
func (ctx *BotContext) DeleteMessage(messageID string) error {
	request := reqMsg.Delete{
		MessageID: messageID,
	}

	_, err := ctx.Client.Call(ctx.Ctx, DeleteMsg, request)
	return err
}

// SetState sets the FSM state for the current user
func (ctx *BotContext) SetState(state string) bool {
	return ctx.FSM.SetStateUser(ctx.UserID, state)
}

// GetState gets the current FSM state for the user
func (ctx *BotContext) GetState() (string, bool) {
	return ctx.FSM.GetStateUser(ctx.UserID)
}

// SetData stores custom data for the current user in FSM
func (ctx *BotContext) SetData(key, value string) bool {
	return ctx.FSM.SetValue(ctx.UserID, key, value)
}

// GetData retrieves custom data for the current user from FSM
func (ctx *BotContext) GetData(key string) (string, bool) {
	return ctx.FSM.GetValue(ctx.UserID, key)
}

// sendMessage is a helper method to send message to the appropriate target (user or chat)
func (ctx *BotContext) sendMessage(body *models.NewMessageBody) error {
	request := reqMsg.Send{
		NewMessageBody: *body,
	}

	if ctx.ChatID != 0 {
		request.ChatID = &ctx.ChatID
	} else if ctx.UserID != 0 {
		request.UserID = &ctx.UserID
	}

	_, err := ctx.Client.Call(ctx.Ctx, SendMsg, request)
	return err
}
