package chat

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/external"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

const attachmentURLPrefix = "https://storage.googleapis.com/instay-dev-public-bucket/"

type chatUseCaseImpl struct {
	logPro              port.LogProvider
	chatRepo            repository.ChatRepository
	messageRepo         repository.MessageRepository
	bookingRoomRepo     repository.BookingRoomRepository
	userRepo            repository.UserRepository
	serviceRepo         repository.ServiceRepository
	supportCategoryRepo repository.SupportCategoryRepository
}

func NewChatUseCase(
	logPro port.LogProvider,
	chatRepo repository.ChatRepository,
	messageRepo repository.MessageRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	userRepo repository.UserRepository,
	serviceRepo repository.ServiceRepository,
	supportCategoryRepo repository.SupportCategoryRepository,
) ChatUseCase {
	return &chatUseCaseImpl{
		logPro:              logPro,
		chatRepo:            chatRepo,
		messageRepo:         messageRepo,
		bookingRoomRepo:     bookingRoomRepo,
		userRepo:            userRepo,
		serviceRepo:         serviceRepo,
		supportCategoryRepo: supportCategoryRepo,
	}
}

func (u *chatUseCaseImpl) ListForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.ChatPaginationQuery) ([]*response.BasicChatResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	chats, total, err := u.chatRepo.FindAllByBookingRoomIDWithDepartmentPaginated(ctx, schema, bookingRoomID, mapper.ToChatQuery(query))
	if err != nil {
		u.logPro.Error("Find all chat by booking room with department paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToBasicChatsResponse(chats), meta, nil
}

func (u *chatUseCaseImpl) ListMessagesByIDForGuest(ctx context.Context, schema string, chatID int64, query request.MessagePaginationQuery, includeChat bool) (*response.BasicChatResponse, []*response.MessageResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	var chatRes *response.BasicChatResponse
	if includeChat {
		chat, err := u.chatRepo.FindByIDWithDepartment(ctx, schema, chatID)
		if err != nil {
			u.logPro.Error("Find chat by id with department failed", port.String("schema", schema), port.Int64("chat_id", chatID), port.Error(err))
			return nil, nil, nil, err
		}
		if chat == nil {
			return nil, nil, nil, exceptions.ErrChatNotFound
		}

		chatRes = mapper.ToBasicChatResponse(chat)
	}

	msgs, total, err := u.messageRepo.FindAllByChatIDPaginated(ctx, schema, chatID, query.Page, query.Limit)
	if err != nil {
		u.logPro.Error("Find all messages by chat id paginated failed", port.String("schema", schema), port.Int64("chat_id", chatID), port.Error(err))
		return nil, nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return chatRes, mapper.ToMessagesResponse(msgs), meta, nil
}

func (u *chatUseCaseImpl) AvailabilityForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.ChatAvailabilityForGuestQuery) (*response.ChatAvailabilityForGuestResponse, error) {
	if (query.ServiceID == 0 && query.SupportCategoryID == 0) || (query.ServiceID != 0 && query.SupportCategoryID != 0) {
		return nil, exceptions.ErrBadRequest.WithData(map[string]any{
			"field": "service_id/support_category_id",
			"tag":   "xor",
			"param": "",
		})
	}

	var departmentID int64

	if query.ServiceID != 0 {
		service, err := u.serviceRepo.FindByIDWithServiceType(ctx, schema, query.ServiceID)
		if err != nil {
			u.logPro.Error("Find service by id failed", port.String("schema", schema), port.Int64("service_id", query.ServiceID), port.Error(err))
			return nil, err
		}
		if service == nil || service.ServiceType == nil {
			return nil, exceptions.ErrServiceNotFound
		}

		departmentID = service.ServiceType.DepartmentID
	}

	if query.SupportCategoryID != 0 {
		supportCategory, err := u.supportCategoryRepo.FindByID(ctx, schema, query.SupportCategoryID)
		if err != nil {
			u.logPro.Error("Find support category by id failed", port.String("schema", schema), port.Int64("support_category_id", query.SupportCategoryID), port.Error(err))
			return nil, err
		}
		if supportCategory == nil {
			return nil, exceptions.ErrSupportCategoryNotFound
		}

		departmentID = supportCategory.DepartmentID
	}

	chat, err := u.chatRepo.FindByBookingRoomIDAndDepartmentID(ctx, schema, bookingRoomID, departmentID)
	if err != nil {
		u.logPro.Error("Find chat by booking room id and department id failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	if chat == nil {
		return &response.ChatAvailabilityForGuestResponse{
			Available:    false,
			ChatID:       0,
			DepartmentID: departmentID,
		}, nil
	}

	return &response.ChatAvailabilityForGuestResponse{
		Available:    true,
		ChatID:       chat.ID,
		DepartmentID: 0,
	}, nil
}

func (u *chatUseCaseImpl) AvailabilityForAdmin(ctx context.Context, schema string, userID, departmentID int64, query request.ChatAvailabilityForAdminQuery) (*response.ChatAvailabilityForAdminResponse, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, query.BookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", query.BookingRoomID), port.Error(err))
		return nil, err
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return nil, err
	}

	chat, err := u.chatRepo.FindByBookingRoomIDAndDepartmentID(ctx, schema, query.BookingRoomID, departmentID)
	if err != nil {
		u.logPro.Error("Find latest chat by booking room id failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	if chat == nil {
		return &response.ChatAvailabilityForAdminResponse{
			Available: false,
			ChatID:    0,
		}, nil
	}

	return &response.ChatAvailabilityForAdminResponse{
		Available: true,
		ChatID:    chat.ID,
	}, nil
}

func (u *chatUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID int64, departmentID *int64, query request.ChatPaginationQuery) ([]*response.SimpleChatResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	chatQuery := mapper.ToChatQuery(query)
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, nil, err
	}
	if !scope.IsGlobal() {
		chatQuery.BranchIDs = scope.BranchIDs
	}

	chats, total, err := u.chatRepo.FindAllWithBookingRoomDetailsPaginated(ctx, schema, departmentID, chatQuery)
	if err != nil {
		u.logPro.Error("Find all chat with booking room details paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSimpleChatsResponse(chats), meta, nil
}

func (u *chatUseCaseImpl) ListMessagesByIDForAdmin(ctx context.Context, schema string, userID, chatID int64, query request.MessagePaginationQuery, includeChat bool) (*response.SimpleChatResponse, []*response.MessageResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	var chatRes *response.SimpleChatResponse
	chat, err := u.chatRepo.FindByIDWithBookingRoomDetails(ctx, schema, chatID)
	if err != nil {
		u.logPro.Error("Find chat by id with booking room details failed", port.String("schema", schema), port.Int64("chat_id", chatID), port.Error(err))
		return nil, nil, nil, err
	}
	if chat == nil {
		return nil, nil, nil, exceptions.ErrChatNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, chat.BookingRoom); err != nil {
		return nil, nil, nil, err
	}
	if includeChat {
		chatRes = mapper.ToSimpleChatResponse(chat)
	}

	msgs, total, err := u.messageRepo.FindAllByChatIDPaginated(ctx, schema, chatID, query.Page, query.Limit)
	if err != nil {
		u.logPro.Error("Find all messages by chat id paginated failed", port.String("schema", schema), port.Int64("chat_id", chatID), port.Error(err))
		return nil, nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return chatRes, mapper.ToMessagesResponse(msgs), meta, nil
}

func (u *chatUseCaseImpl) ListForExternal(ctx context.Context, schema string, query external.ExternalConversationsQuery) ([]*external.ExternalConversationResponse, string, error) {
	if query.Limit == 0 {
		query.Limit = 100
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	cursorAt, cursorID, err := decodeConversationCursor(query.Cursor)
	if err != nil {
		return nil, "", exceptions.ErrBadRequest
	}

	chats, hasNext, err := u.chatRepo.FindAllExternalBySinceAndCursor(ctx, schema, query.Since, cursorAt, cursorID, query.Limit)
	if err != nil {
		u.logPro.Error("Find external conversations failed", port.String("schema", schema), port.Error(err))
		return nil, "", err
	}

	convRes := make([]*external.ExternalConversationResponse, 0, len(chats))
	for _, chat := range chats {
		lastMessageAt := chat.CreatedAt
		if chat.LastMessageAt != nil {
			lastMessageAt = *chat.LastMessageAt
		}

		customerName, metadata := buildCustomerNameAndMetadata(chat)
		conv := &external.ExternalConversationResponse{
			ExternalID:     chat.ID,
			ExternalUserID: chat.BookingRoomID,
			CustomerName:   customerName,
			LastMessageAt:  lastMessageAt,
			Metadata:       metadata,
		}

		convRes = append(convRes, conv)
	}

	nextCursor := ""
	if hasNext && len(chats) > 0 {
		last := chats[len(chats)-1]
		lastAt := last.CreatedAt
		if last.LastMessageAt != nil {
			lastAt = *last.LastMessageAt
		}

		nextCursor, err = encodeConversationCursor(lastAt, last.ID)
		if err != nil {
			u.logPro.Error("Encode conversations cursor failed", port.String("schema", schema), port.Error(err))
			return nil, "", err
		}
	}

	return convRes, nextCursor, nil
}

func (u *chatUseCaseImpl) ListMessagesForExternal(ctx context.Context, schema string, query external.ExternalMessagesQuery) ([]*external.ExternalMessageResponse, string, error) {
	if query.Limit == 0 {
		query.Limit = 100
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	cursorAt, cursorID, err := decodeMessageCursor(query.Cursor)
	if err != nil {
		return nil, "", exceptions.ErrBadRequest
	}

	msgs, hasNext, err := u.messageRepo.FindAllByChatIDSinceCursor(ctx, schema, query.ConversationID, query.Since, cursorAt, cursorID, query.Limit)
	if err != nil {
		u.logPro.Error("Find external messages failed", port.String("schema", schema), port.Int64("chat_id", query.ConversationID), port.Error(err))
		return nil, "", err
	}

	msgRes := make([]*external.ExternalMessageResponse, 0, len(msgs))
	for _, msg := range msgs {
		contentType := "text"
		if msg.AttachmentType != nil && *msg.AttachmentType != "" {
			contentType = string(*msg.AttachmentType)
		}

		content := ""
		if msg.Content != nil {
			content = *msg.Content
		}

		attachments := make([]*external.ExternalAttachmentResponse, 0)
		if msg.AttachmentKey != nil && *msg.AttachmentKey != "" {
			key := strings.TrimPrefix(*msg.AttachmentKey, "/")
			attachmentType := "file"
			if msg.AttachmentType != nil && *msg.AttachmentType != "" {
				attachmentType = string(*msg.AttachmentType)
			}

			attachments = append(attachments, &external.ExternalAttachmentResponse{
				Type: attachmentType,
				URL:  attachmentURLPrefix + key,
				Name: path.Base(key),
			})
		}

		msgRes = append(msgRes, &external.ExternalMessageResponse{
			ExternalID:  msg.ID,
			SenderType:  mapSenderType(msg.SenderType),
			SenderName:  msg.SenderName,
			Content:     content,
			ContentType: contentType,
			Attachments: attachments,
			SentAt:      msg.CreatedAt,
		})
	}

	nextCursor := ""
	if hasNext && len(msgs) > 0 {
		last := msgs[len(msgs)-1]
		nextCursor, err = encodeMessageCursor(last.CreatedAt, last.ID)
		if err != nil {
			u.logPro.Error("Encode messages cursor failed", port.String("schema", schema), port.Error(err))
			return nil, "", err
		}
	}

	return msgRes, nextCursor, nil
}

func (u *chatUseCaseImpl) ensureBookingRoomInUserScope(ctx context.Context, schema string, userID int64, bookingRoom *model.BookingRoom) error {
	if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil {
		return exceptions.ErrBookingRoomNotFound
	}

	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	return scope.Ensure(bookingRoom.Room.Floor.BranchID)
}

type conversationCursor struct {
	LastMessageAt time.Time `json:"last_message_at"`
	ChatID        int64     `json:"chat_id"`
}

type messageCursor struct {
	SentAt    time.Time `json:"sent_at"`
	MessageID string    `json:"message_id"`
}

func encodeConversationCursor(lastMessageAt time.Time, chatID int64) (string, error) {
	raw, err := json.Marshal(conversationCursor{
		LastMessageAt: lastMessageAt,
		ChatID:        chatID,
	})
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func decodeConversationCursor(encoded string) (*time.Time, *int64, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return nil, nil, nil
	}

	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, nil, err
	}

	var c conversationCursor
	if err = json.Unmarshal(raw, &c); err != nil {
		return nil, nil, err
	}
	if c.ChatID == 0 || c.LastMessageAt.IsZero() {
		return nil, nil, exceptions.ErrBadRequest
	}

	return &c.LastMessageAt, &c.ChatID, nil
}

func encodeMessageCursor(sentAt time.Time, messageID string) (string, error) {
	raw, err := json.Marshal(messageCursor{
		SentAt:    sentAt,
		MessageID: messageID,
	})
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func decodeMessageCursor(encoded string) (*time.Time, string, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return nil, "", nil
	}

	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, "", err
	}

	var c messageCursor
	if err = json.Unmarshal(raw, &c); err != nil {
		return nil, "", err
	}
	if c.MessageID == "" || c.SentAt.IsZero() {
		return nil, "", exceptions.ErrBadRequest
	}

	return &c.SentAt, c.MessageID, nil
}

func buildCustomerNameAndMetadata(chat *model.Chat) (string, map[string]any) {
	if chat == nil || chat.BookingRoom == nil {
		return "", nil
	}

	bookingNumber := ""
	if chat.BookingRoom.Booking != nil {
		bookingNumber = strings.TrimSpace(chat.BookingRoom.Booking.BookingNumber)
	}

	roomName := ""
	if chat.BookingRoom.Room != nil {
		roomName = strings.TrimSpace(chat.BookingRoom.Room.Name)
	}

	parts := make([]string, 0, 2)
	if bookingNumber != "" {
		parts = append(parts, bookingNumber)
	}
	if roomName != "" {
		parts = append(parts, roomName)
	}

	metadata := map[string]any{}
	if bookingNumber != "" {
		metadata["booking_number"] = bookingNumber
	}
	if roomName != "" {
		metadata["room_name"] = roomName
	}

	customerName := strings.Join(parts, " - ")
	if len(metadata) == 0 {
		metadata = nil
	}

	return customerName, metadata
}

func mapSenderType(senderType model.MessageActionType) string {
	switch senderType {
	case model.MessageActionTypeGuest:
		return "customer"
	case model.MessageActionTypeStaff:
		return "agent"
	default:
		return string(senderType)
	}
}
