package message

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type messageUseCaseImpl struct {
	logPro          port.LogProvider
	idPro           port.IDProvider
	txPro           port.TxProvider
	chatRepo        repository.ChatRepository
	bookingRoomRepo repository.BookingRoomRepository
	messageRepo     repository.MessageRepository
}

func NewMessageUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	chatRepo repository.ChatRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	messageRepo repository.MessageRepository,
) MessageUseCase {
	return &messageUseCaseImpl{
		logPro,
		idPro,
		txPro,
		chatRepo,
		bookingRoomRepo,
		messageRepo,
	}
}

func (u *messageUseCaseImpl) Create(ctx context.Context, schema string, req request.CreateMessageRequest) (*response.MessageResponse, error) {
	if req.ChatID == nil {
		bookingRoom, err := u.bookingRoomRepo.FindByIDWithBooking(ctx, schema, req.BookingRoomID)
		if err != nil {
			u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", req.BookingRoomID), port.Error(err))
			return nil, err
		}
		if bookingRoom == nil {
			return nil, exceptions.ErrBookingRoomNotFound
		}

		id, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate department id failed", port.String("schema", schema), port.Error(err))
			return nil, err
		}

		chat := model.NewChat(id, bookingRoom.ID, req.DepartmentID, bookingRoom.Booking.CheckOut)

		if err = u.chatRepo.Create(ctx, schema, chat); err != nil {
			if errors.Is(err, exceptions.ErrChatAlreadyExists) {
				return nil, err
			}
			u.logPro.Error("Create chat failed", port.String("schema", schema), port.Error(err))
			return nil, err
		}

		req.ChatID = &chat.ID
	}

	var messageReply *model.MessageReply
	if req.MessageID != nil {
		replyMsg, err := u.messageRepo.FindByID(ctx, *req.MessageID)
		if err != nil {
			u.logPro.Error("Find reply message by id failed", port.String("schema", schema), port.String("message_id", *req.MessageID), port.Error(err))
			return nil, err
		}
		if replyMsg == nil {
			return nil, exceptions.ErrMessageNotFound
		}

		messageReply = &model.MessageReply{
			MessageID:      replyMsg.ID,
			SenderName:     replyMsg.SenderName,
			Content:        replyMsg.Content,
			AttachmentType: replyMsg.AttachmentType,
		}
	}

	message := model.NewMessage(
		bson.NewObjectID().Hex(), schema, *req.ChatID,
		req.SenderType, req.SenderID, req.SenderName,
		req.AttachmentKey, req.AttachmentType, req.Content, nil,
		false, messageReply, time.Now())

	if err := u.messageRepo.Create(ctx, message); err != nil {
		u.logPro.Error("Create message failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	u.updateChatLastMessage(schema, *req.ChatID, message)

	return mapper.ToMessageResponse(message), nil
}

func (u *messageUseCaseImpl) UpdateRead(ctx context.Context, schema string, req request.UpdateMessageReadRequest) error {
	now := time.Now()

	switch req.ReaderType {
	case model.MessageActionTypeGuest:
		if err := u.messageRepo.SetBookingRoomReadAtAllByChatID(ctx, schema, req.ChatID, now); err != nil {
			u.logPro.Error("Update booking room read at messages failed", port.String("schema", schema), port.Error(err))
			return err
		}
	case model.MessageActionTypeStaff:
		messageUser := &model.MessageUser{
			UserID:       req.ReaderID,
			UserFullName: req.ReaderName,
			ReadAt:       now,
		}

		if err := u.messageRepo.PushMessageUserAllByChatID(ctx, schema, req.ChatID, messageUser); err != nil {
			u.logPro.Error("Push message user to messages failed", port.String("schema", schema), port.Error(err))
			return err
		}
	default:
		return exceptions.ErrBadRequest
	}

	return nil
}

func (u *messageUseCaseImpl) UpdateReaction(ctx context.Context, schema string, req request.UpdateMessageReactionRequest) (*response.UpdateMessageReactionResponse, error) {
	msg, err := u.messageRepo.FindByID(ctx, req.MessageID)
	if err != nil {
		u.logPro.Error("Find message by id failed", port.String("schema", schema), port.String("message_id", req.MessageID), port.Error(err))
		return nil, err
	}
	if msg == nil || msg.Schema != schema || msg.ChatID != req.ChatID {
		return nil, exceptions.ErrMessageNotFound
	}

	now := time.Now()
	reaction := &model.MessageReaction{
		SenderType: string(req.SenderType),
		SenderID:   req.SenderID,
		SenderName: req.SenderName,
		Emoji:      req.Emoji,
		CreatedAt:  now,
	}

	reactions := msg.MessageReactions
	if reactions == nil {
		reactions = make([]*model.MessageReaction, 0)
	}

	idx := findReactionIndexBySender(reactions, req.SenderID, req.SenderType)
	isReacted := true

	if idx >= 0 {
		if reactions[idx].Emoji == req.Emoji {
			reactions = append(reactions[:idx], reactions[idx+1:]...)
			isReacted = false
		} else {
			reactions[idx] = reaction
		}
	} else {
		reactions = append(reactions, reaction)
	}

	if err = u.messageRepo.SetMessageReactionsByID(ctx, schema, req.MessageID, reactions); err != nil {
		u.logPro.Error(
			"set message reactions failed",
			port.String("schema", schema),
			port.String("message_id", req.MessageID),
			port.Error(err),
		)
		return nil, err
	}

	return &response.UpdateMessageReactionResponse{
		MessageID:        req.MessageID,
		ChatID:           req.ChatID,
		SenderType:       req.SenderType,
		SenderID:         req.SenderID,
		SenderName:       req.SenderName,
		Emoji:            req.Emoji,
		IsReacted:        isReacted,
		MessageReactions: mapper.ToMessageReactionsResponse(reactions),
	}, nil
}

func findReactionIndexBySender(
	reactions []*model.MessageReaction,
	senderID int64,
	senderType model.MessageActionType,
) int {
	for i, reaction := range reactions {
		if reaction == nil {
			continue
		}
		if reaction.SenderID == senderID && reaction.SenderType == string(senderType) {
			return i
		}
	}

	return -1
}

func (u *messageUseCaseImpl) updateChatLastMessage(schema string, chatID int64, msg *model.Message) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		updateData := map[string]any{
			"last_message_at":              msg.CreatedAt,
			"last_message_content":         msg.Content,
			"last_message_attachment_type": msg.AttachmentType,
			"last_message_sender_type":     msg.SenderType,
			"last_message_sender_name":     msg.SenderName,
		}

		if err := u.chatRepo.Update(ctx, schema, chatID, updateData); err != nil {
			u.logPro.Error("Update chat failed", port.String("schema", schema), port.Error(err))
		}
	}()
}
