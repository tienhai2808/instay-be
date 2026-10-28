package repository

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error

	FindByID(ctx context.Context, id string) (*model.Message, error)

	FindAllByChatIDPaginated(ctx context.Context, schema string, chatID int64, page, limit uint32) ([]*model.Message, int64, error)

	FindAllByChatIDSinceCursor(ctx context.Context, schema string, chatID int64, since, cursorAt *time.Time, cursorID string, limit uint32) ([]*model.Message, bool, error)

	SetMessageReactionsByID(ctx context.Context, schema, id string, reactions []*model.MessageReaction) error

	SetBookingRoomReadAtAllByChatID(ctx context.Context, schema string, chatID int64, readAt time.Time) error

	PushMessageUserAllByChatID(ctx context.Context, schema string, chatID int64, messageUser *model.MessageUser) error

	CountDashboardUnreadGuestMessages(ctx context.Context, schema string, chatIDs ...[]int64) (int64, error)
}
