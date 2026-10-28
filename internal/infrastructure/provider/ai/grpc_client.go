package ai

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	aipb "github.com/InstaySystem/is_v2-be/proto/gen/ai"
)

type aiProviderImpl struct {
	client aipb.AIServiceClient
}

func NewAIProvider(client aipb.AIServiceClient) port.AIProvider {
	return &aiProviderImpl{client}
}

func (p *aiProviderImpl) Ping(ctx context.Context) error {
	if _, err := p.client.Ping(ctx, &aipb.PingRequest{}); err != nil {
		return err
	}

	return nil
}

func (p *aiProviderImpl) AnalyzeReview(ctx context.Context, reviewID int64, reviewContent string) (int64, []model.ReviewAspect, error) {
	res, err := p.client.AnalyzeReview(ctx, &aipb.AnalyzeReviewRequest{
		Text: reviewContent,
	})
	if err != nil {
		return 0, nil, err
	}

	aspects := make([]model.ReviewAspect, 0, len(res.Aspects))
	for _, aspect := range res.Aspects {
		aspects = append(aspects, model.ReviewAspect{
			Category:  aspect.Aspect,
			Sentiment: aspect.Sentiment,
			Score:     aspect.Score,
		})
	}

	return reviewID, aspects, nil
}

func (p *aiProviderImpl) PredictBooking(ctx context.Context, booking *model.Booking) (*port.PredictBookingResult, error) {
	res, err := p.client.PredictBooking(ctx, &aipb.PredictBookingRequest{
		Id:       booking.ID,
		Features: toBookingFeatures(booking),
	})
	if err != nil {
		return nil, err
	}

	return &port.PredictBookingResult{
		ID:                      res.Id,
		CancellationProbability: res.CancellationProbability,
		PredictionLabel:         res.PredictionLabel,
		WillCancel:              res.WillCancel,
		RiskLevel:               res.RiskLevel,
	}, nil
}

func toBookingFeatures(booking *model.Booking) *aipb.BookingFeatures {
	checkIn := normalizeDate(booking.CheckIn)
	checkOut := normalizeDate(booking.CheckOut)
	if !checkOut.After(checkIn) {
		checkOut = checkIn.AddDate(0, 0, 1)
	}

	leadTime := int32(0)
	if booking.CheckIn.After(booking.BookedOn) {
		leadDuration := booking.CheckIn.Sub(booking.BookedOn)
		leadDays := int32(math.Ceil(leadDuration.Hours() / 24))
		if leadDays > 0 {
			leadTime = leadDays
		}
	}

	_, week := checkIn.ISOWeek()
	weekendNights, weekNights := countStayNights(checkIn, checkOut)

	feature := &aipb.BookingFeatures{
		LeadTime:              ptr(leadTime),
		ArrivalDateYear:       ptr(safeInt32FromInt(checkIn.Year())),
		ArrivalDateMonth:      ptr(checkIn.Month().String()),
		ArrivalDateWeekNumber: ptr(safeInt32FromInt(week)),
		ArrivalDateDayOfMonth: ptr(safeInt32FromInt(checkIn.Day())),
		StaysInWeekendNights:  ptr(weekendNights),
		StaysInWeekNights:     ptr(weekNights),
		Adr:                   ptr(calculateADR(booking, checkIn, checkOut)),
	}

	if adults := parseAdults(booking.GuestNumber); adults != nil {
		feature.Adults = adults
	}

	if booking.MealPlan != nil && strings.TrimSpace(*booking.MealPlan) != "" {
		feature.Meal = ptr(strings.TrimSpace(*booking.MealPlan))
	}

	if booking.Source != nil && strings.TrimSpace(booking.Source.Name) != "" {
		sourceName := strings.TrimSpace(booking.Source.Name)
		feature.MarketSegment = ptr(sourceName)
		feature.DistributionChannel = ptr(sourceName)
	}

	if booking.RoomType != nil && strings.TrimSpace(booking.RoomType.Name) != "" {
		feature.ReservedRoomType = ptr(strings.TrimSpace(booking.RoomType.Name))
	}

	if booking.BookingPreferences != nil && strings.TrimSpace(*booking.BookingPreferences) != "" {
		feature.TotalOfSpecialRequests = ptr(int32(1))
	}

	return feature
}

func calculateADR(booking *model.Booking, checkIn, checkOut time.Time) float64 {
	stayDays := checkOut.Sub(checkIn).Hours() / 24
	if stayDays <= 0 {
		stayDays = 1
	}

	roomNumber := float64(booking.RoomNumber)
	if roomNumber < 1 {
		roomNumber = 1
	}

	return booking.TotalSellPrice / stayDays / roomNumber
}

func countStayNights(checkIn, checkOut time.Time) (int32, int32) {
	var weekendNights int32
	var weekNights int32

	for day := checkIn; day.Before(checkOut); day = day.AddDate(0, 0, 1) {
		switch day.Weekday() {
		case time.Saturday, time.Sunday:
			weekendNights++
		default:
			weekNights++
		}
	}

	return weekendNights, weekNights
}

func parseAdults(guestNumber string) *int32 {
	guestNumber = strings.TrimSpace(guestNumber)
	if guestNumber == "" {
		return nil
	}

	if value, err := strconv.ParseInt(guestNumber, 10, 32); err == nil {
		if value < 0 {
			value = 0
		}
		return ptr(int32(value))
	}

	chunks := strings.FieldsFunc(guestNumber, func(r rune) bool {
		return r < '0' || r > '9'
	})

	for _, chunk := range chunks {
		if chunk == "" {
			continue
		}

		value, err := strconv.ParseInt(chunk, 10, 32)
		if err != nil {
			continue
		}

		if value < 0 {
			value = 0
		}

		return ptr(int32(value))
	}

	return nil
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ptr[T any](v T) *T {
	return &v
}

func safeInt32FromInt(v int) int32 {
	const (
		maxInt32 = int(^uint32(0) >> 1)
		minInt32 = -maxInt32 - 1
	)

	if v > maxInt32 {
		return int32(maxInt32)
	}
	if v < minInt32 {
		return int32(minInt32)
	}

	return int32(v)
}
