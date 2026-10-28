package port

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type AIProvider interface {
	Ping(ctx context.Context) error

	AnalyzeReview(ctx context.Context, reviewID int64, reviewContent string) (int64, []model.ReviewAspect, error)

	PredictBooking(ctx context.Context, booking *model.Booking) (*PredictBookingResult, error)
}

type PredictBookingResult struct {
	ID                      int64
	CancellationProbability float32
	PredictionLabel         int32
	WillCancel              bool
	RiskLevel               string
}
