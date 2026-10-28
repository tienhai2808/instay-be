package recalculate

import "context"

type RecalculateUseCase interface {
	RecalculateDailyRoomTypeRate(ctx context.Context, schema string, roomTypeID int64) error
}
