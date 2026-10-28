package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBasicChatResponse(chat *model.Chat) *response.BasicChatResponse {
	if chat == nil {
		return nil
	}

	return &response.BasicChatResponse{
		ID:                        chat.ID,
		Department:                ToBasicDepartmentResponse(chat.Department),
		LastMessageAt:             chat.LastMessageAt,
		LastMessageContent:        chat.LastMessageContent,
		LastMessageAttachmentType: chat.LastMessageAttachmentType,
		LastMessageSenderType:     chat.LastMessageSenderType,
		LastMessageSenderName:     chat.LastMessageSenderName,
		ExpiredAt:                 chat.ExpiredAt,
		CreatedAt:                 chat.CreatedAt,
	}
}

func ToBasicChatsResponse(chats []*model.Chat) []*response.BasicChatResponse {
	if len(chats) == 0 {
		return make([]*response.BasicChatResponse, 0)
	}

	chatsRes := make([]*response.BasicChatResponse, 0, len(chats))
	for _, chat := range chats {
		chatsRes = append(chatsRes, ToBasicChatResponse(chat))
	}

	return chatsRes
}

func ToSimpleChatResponse(chat *model.Chat) *response.SimpleChatResponse {
	if chat == nil {
		return nil
	}

	return &response.SimpleChatResponse{
		ID:                        chat.ID,
		BookingRoom:               ToSimpleBookingRoomResponse(chat.BookingRoom),
		LastMessageAt:             chat.LastMessageAt,
		LastMessageContent:        chat.LastMessageContent,
		LastMessageAttachmentType: chat.LastMessageAttachmentType,
		LastMessageSenderType:     chat.LastMessageSenderType,
		LastMessageSenderName:     chat.LastMessageSenderName,
		ExpiredAt:                 chat.ExpiredAt,
		CreatedAt:                 chat.CreatedAt,
	}
}

func ToSimpleChatsResponse(chats []*model.Chat) []*response.SimpleChatResponse {
	if len(chats) == 0 {
		return make([]*response.SimpleChatResponse, 0)
	}

	chatsRes := make([]*response.SimpleChatResponse, 0, len(chats))
	for _, chat := range chats {
		chatsRes = append(chatsRes, ToSimpleChatResponse(chat))
	}

	return chatsRes
}

func ToChatQuery(query request.ChatPaginationQuery) model.ChatQuery {
	return model.ChatQuery{
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
	}
}

func ToMessageResponse(msg *model.Message) *response.MessageResponse {
	if msg == nil {
		return nil
	}

	return &response.MessageResponse{
		ID:                msg.ID,
		ChatID:            msg.ChatID,
		SenderType:        msg.SenderType,
		SenderName:        msg.SenderName,
		SenderID:          msg.SenderID,
		Content:           msg.Content,
		AttachmentKey:     msg.AttachmentKey,
		AttachmentType:    msg.AttachmentType,
		BookingRoomReadAt: msg.BookingRoomReadAt,
		MessageUsers:      ToMessageUsersResponse(msg.MessageUsers),
		MessageReactions:  ToMessageReactionsResponse(msg.MessageReactions),
		MessageReply:      ToMessageRelyResponse(msg.MessageReply),
		IsDeleted:         msg.IsDeleted,
		CreatedAt:         msg.CreatedAt,
	}
}

func ToMessagesResponse(msgs []*model.Message) []*response.MessageResponse {
	if len(msgs) == 0 {
		return make([]*response.MessageResponse, 0)
	}

	msgsRes := make([]*response.MessageResponse, 0, len(msgs))
	for _, msg := range msgs {
		msgsRes = append(msgsRes, ToMessageResponse(msg))
	}

	return msgsRes
}

func ToMessageRelyResponse(mr *model.MessageReply) *response.MessageReplyResponse {
	if mr == nil {
		return nil
	}

	return &response.MessageReplyResponse{
		MessageID:      mr.MessageID,
		Content:        mr.Content,
		SenderName:     mr.SenderName,
		AttachmentType: mr.AttachmentType,
	}
}

func ToMessageReactionResponse(mr *model.MessageReaction) *response.MessageReactionResponse {
	if mr == nil {
		return nil
	}

	return &response.MessageReactionResponse{
		SenderType: mr.SenderType,
		SenderID:   mr.SenderID,
		SenderName: mr.SenderName,
		Emoji:      mr.Emoji,
		CreatedAt:  mr.CreatedAt,
	}
}

func ToMessageReactionsResponse(mrs []*model.MessageReaction) []*response.MessageReactionResponse {
	if len(mrs) == 0 {
		return make([]*response.MessageReactionResponse, 0)
	}

	mrsRes := make([]*response.MessageReactionResponse, 0, len(mrs))
	for _, mr := range mrs {
		mrsRes = append(mrsRes, ToMessageReactionResponse(mr))
	}

	return mrsRes
}

func ToMessageUserResponse(mu *model.MessageUser) *response.MessageUserResponse {
	if mu == nil {
		return nil
	}

	return &response.MessageUserResponse{
		UserID:       mu.UserID,
		UserFullName: mu.UserFullName,
		ReadAt:       mu.ReadAt,
	}
}

func ToMessageUsersResponse(mus []*model.MessageUser) []*response.MessageUserResponse {
	if len(mus) == 0 {
		return make([]*response.MessageUserResponse, 0)
	}

	musRes := make([]*response.MessageUserResponse, 0, len(mus))
	for _, mu := range mus {
		musRes = append(musRes, ToMessageUserResponse(mu))
	}

	return musRes
}
