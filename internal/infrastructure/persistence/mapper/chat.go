package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	mongoEntity "github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mongo/entity"
	ormEntity "github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToChatEntity(chat *model.Chat) *ormEntity.ChatEntity {
	if chat == nil {
		return nil
	}

	return &ormEntity.ChatEntity{
		ID:                        chat.ID,
		CreatedAt:                 chat.CreatedAt,
		BookingRoomID:             chat.BookingRoomID,
		DepartmentID:              chat.DepartmentID,
		ExpiredAt:                 chat.ExpiredAt,
		LastMessageAt:             chat.LastMessageAt,
		LastMessageContent:        chat.LastMessageContent,
		LastMessageAttachmentType: (*string)(chat.LastMessageAttachmentType),
		LastMessageSenderName:     chat.LastMessageSenderName,
		LastMessageSenderType:     (*string)(chat.LastMessageSenderType),
	}
}

func ToChat(chat *ormEntity.ChatEntity) *model.Chat {
	if chat == nil {
		return nil
	}

	return &model.Chat{
		ID:                        chat.ID,
		CreatedAt:                 chat.CreatedAt,
		BookingRoomID:             chat.BookingRoomID,
		DepartmentID:              chat.DepartmentID,
		ExpiredAt:                 chat.ExpiredAt,
		LastMessageAt:             chat.LastMessageAt,
		LastMessageContent:        chat.LastMessageContent,
		LastMessageAttachmentType: (*model.MessageAttachmentType)(chat.LastMessageAttachmentType),
		LastMessageSenderName:     chat.LastMessageSenderName,
		LastMessageSenderType:     (*model.MessageActionType)(chat.LastMessageSenderType),
		BookingRoom:               ToBookingRoom(chat.BookingRoom),
		Department:                ToDepartment(chat.Department),
	}
}

func ToChats(chats []*ormEntity.ChatEntity) []*model.Chat {
	if len(chats) == 0 {
		return make([]*model.Chat, 0)
	}

	chatsRes := make([]*model.Chat, 0, len(chats))
	for _, chat := range chats {
		chatsRes = append(chatsRes, ToChat(chat))
	}

	return chatsRes
}

func ToMessageEntity(message *model.Message) *mongoEntity.MessageEntity {
	if message == nil {
		return nil
	}

	objID, err := bson.ObjectIDFromHex(message.ID)
	if err != nil {
		return nil
	}

	return &mongoEntity.MessageEntity{
		ID:                objID,
		Schema:            message.Schema,
		ChatID:            message.ChatID,
		SenderType:        string(message.SenderType),
		SenderID:          message.SenderID,
		SenderName:        message.SenderName,
		AttachmentKey:     message.AttachmentKey,
		AttachmentType:    (*string)(message.AttachmentType),
		Content:           message.Content,
		BookingRoomReadAt: message.BookingRoomReadAt,
		CreatedAt:         message.CreatedAt,
		IsDeleted:         message.IsDeleted,
		MessageReply:      ToMessageReplyEntity(message.MessageReply),
	}
}

func ToMessageReactionEntity(mr *model.MessageReaction) *mongoEntity.MessageReactionEntity {
	if mr == nil {
		return nil
	}

	return &mongoEntity.MessageReactionEntity{
		SenderType: mr.SenderType,
		SenderID:   mr.SenderID,
		SenderName: mr.SenderName,
		Emoji:      string(mr.Emoji),
		CreatedAt:  mr.CreatedAt,
	}
}

func ToMessageReactionsEntity(mrs []*model.MessageReaction) []*mongoEntity.MessageReactionEntity {
	if len(mrs) == 0 {
		return make([]*mongoEntity.MessageReactionEntity, 0)
	}

	mrsRes := make([]*mongoEntity.MessageReactionEntity, 0, len(mrs))
	for _, mr := range mrs {
		mrsRes = append(mrsRes, ToMessageReactionEntity(mr))
	}

	return mrsRes
}

func ToMessageReplyEntity(mr *model.MessageReply) *mongoEntity.MessageReplyEntity {
	if mr == nil {
		return nil
	}

	objID, err := bson.ObjectIDFromHex(mr.MessageID)
	if err != nil {
		return nil
	}

	return &mongoEntity.MessageReplyEntity{
		MessageID:      objID,
		SenderName:     mr.SenderName,
		Content:        mr.Content,
		AttachmentType: (*string)(mr.AttachmentType),
	}
}

func ToMessageUserEntity(mu *model.MessageUser) *mongoEntity.MessageUserEntity {
	if mu == nil {
		return nil
	}

	return &mongoEntity.MessageUserEntity{
		UserID:       mu.UserID,
		UserFullName: mu.UserFullName,
		ReadAt:       mu.ReadAt,
	}
}

func ToMessageReaction(mr *mongoEntity.MessageReactionEntity) *model.MessageReaction {
	if mr == nil {
		return nil
	}

	return &model.MessageReaction{
		SenderType: mr.SenderType,
		SenderID:   mr.SenderID,
		SenderName: mr.SenderName,
		Emoji:      model.MessageReactionEmojiType(mr.Emoji),
		CreatedAt:  mr.CreatedAt,
	}
}

func ToMessageReply(mr *mongoEntity.MessageReplyEntity) *model.MessageReply {
	if mr == nil {
		return nil
	}

	return &model.MessageReply{
		MessageID:      mr.MessageID.Hex(),
		SenderName:     mr.SenderName,
		Content:        mr.Content,
		AttachmentType: (*model.MessageAttachmentType)(mr.AttachmentType),
	}
}

func ToMessageUser(mu *mongoEntity.MessageUserEntity) *model.MessageUser {
	if mu == nil {
		return nil
	}

	return &model.MessageUser{
		UserID:       mu.UserID,
		UserFullName: mu.UserFullName,
		ReadAt:       mu.ReadAt,
	}
}

func ToMessageUsers(mus []*mongoEntity.MessageUserEntity) []*model.MessageUser {
	if len(mus) == 0 {
		return nil
	}

	musRes := make([]*model.MessageUser, 0, len(mus))
	for _, mu := range mus {
		musRes = append(musRes, ToMessageUser(mu))
	}

	return musRes
}

func ToMessageReactions(mrs []*mongoEntity.MessageReactionEntity) []*model.MessageReaction {
	if len(mrs) == 0 {
		return nil
	}

	mrsRes := make([]*model.MessageReaction, 0, len(mrs))
	for _, mr := range mrs {
		mrsRes = append(mrsRes, ToMessageReaction(mr))
	}

	return mrsRes
}

func ToMessage(message *mongoEntity.MessageEntity) *model.Message {
	if message == nil {
		return nil
	}

	return &model.Message{
		ID:                message.ID.Hex(),
		Schema:            message.Schema,
		ChatID:            message.ChatID,
		SenderType:        model.MessageActionType(message.SenderType),
		SenderID:          message.SenderID,
		SenderName:        message.SenderName,
		AttachmentKey:     message.AttachmentKey,
		AttachmentType:    (*model.MessageAttachmentType)(message.AttachmentType),
		Content:           message.Content,
		BookingRoomReadAt: message.BookingRoomReadAt,
		CreatedAt:         message.CreatedAt,
		IsDeleted:         message.IsDeleted,
		MessageReply:      ToMessageReply(message.MessageReply),
		MessageUsers:      ToMessageUsers(message.MessageUsers),
		MessageReactions:  ToMessageReactions(message.MessageReactions),
	}
}

func ToMessages(messages []*mongoEntity.MessageEntity) []*model.Message {
	if len(messages) == 0 {
		return make([]*model.Message, 0)
	}

	messagesRes := make([]*model.Message, 0, len(messages))
	for _, message := range messages {
		messagesRes = append(messagesRes, ToMessage(message))
	}

	return messagesRes
}
