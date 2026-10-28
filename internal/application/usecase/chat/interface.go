package chat

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/external"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type ChatUseCase interface {
	ListForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.ChatPaginationQuery) ([]*response.BasicChatResponse, *response.MetaResponse, error)

	ListMessagesByIDForGuest(ctx context.Context, schema string, chatID int64, query request.MessagePaginationQuery, includeChat bool) (*response.BasicChatResponse, []*response.MessageResponse, *response.MetaResponse, error)

	AvailabilityForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.ChatAvailabilityForGuestQuery) (*response.ChatAvailabilityForGuestResponse, error)

	AvailabilityForAdmin(ctx context.Context, schema string, userID, departmentID int64, query request.ChatAvailabilityForAdminQuery) (*response.ChatAvailabilityForAdminResponse, error)

	ListForAdmin(ctx context.Context, schema string, userID int64, departmentID *int64, query request.ChatPaginationQuery) ([]*response.SimpleChatResponse, *response.MetaResponse, error)

	ListMessagesByIDForAdmin(ctx context.Context, schema string, userID, chatID int64, query request.MessagePaginationQuery, includeChat bool) (*response.SimpleChatResponse, []*response.MessageResponse, *response.MetaResponse, error)

	ListForExternal(ctx context.Context, schema string, query external.ExternalConversationsQuery) ([]*external.ExternalConversationResponse, string, error)

	ListMessagesForExternal(ctx context.Context, schema string, query external.ExternalMessagesQuery) ([]*external.ExternalMessageResponse, string, error)
}
