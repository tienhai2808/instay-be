package dashboard

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"golang.org/x/sync/errgroup"
)

var dashboardReviewCategories = []string{
	"price",
	"service",
	"employee",
	"room",
	"utility",
	"amenities",
}

type dashboardUseCaseImpl struct {
	logPro            port.LogProvider
	bookingRepo       repository.BookingRepository
	reviewRepo        repository.ReviewRepository
	guestRepo         repository.GuestRepository
	supportTicketRepo repository.SupportTicketRepository
	orderServiceRepo  repository.OrderServiceRepository
	messageRepo       repository.MessageRepository
	chatRepo          repository.ChatRepository
	userRepo          repository.UserRepository
}

func NewDashboardUseCase(
	logPro port.LogProvider,
	bookingRepo repository.BookingRepository,
	reviewRepo repository.ReviewRepository,
	guestRepo repository.GuestRepository,
	supportTicketRepo repository.SupportTicketRepository,
	orderServiceRepo repository.OrderServiceRepository,
	messageRepo repository.MessageRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
) DashboardUseCase {
	return &dashboardUseCaseImpl{
		logPro:            logPro,
		bookingRepo:       bookingRepo,
		reviewRepo:        reviewRepo,
		guestRepo:         guestRepo,
		supportTicketRepo: supportTicketRepo,
		orderServiceRepo:  orderServiceRepo,
		messageRepo:       messageRepo,
		chatRepo:          chatRepo,
		userRepo:          userRepo,
	}
}

func (u *dashboardUseCaseImpl) BookingOverview(ctx context.Context, schema string, userID int64, query request.DashboardOverviewQuery) (*response.DashboardBookingOverviewResponse, error) {
	fromDate, toDate, err := normalizeDashboardDateRange(query)
	if err != nil {
		return nil, err
	}
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	currentOverview, err := u.bookingRepo.FindDashboardOverview(ctx, schema, model.DashboardOverviewQuery{
		FromDate:  fromDate,
		ToDate:    toDate,
		BranchIDs: branchIDs,
	})
	if err != nil {
		u.logPro.Error("Find dashboard booking overview failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	previousFromDate, previousToDate := previousDashboardDateRange(fromDate, toDate)

	previousOverview, err := u.bookingRepo.FindDashboardOverview(ctx, schema, model.DashboardOverviewQuery{
		FromDate:  previousFromDate,
		ToDate:    previousToDate,
		BranchIDs: branchIDs,
	})
	if err != nil {
		u.logPro.Error("Find previous dashboard booking overview failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	if currentOverview == nil {
		currentOverview = &model.DashboardBookingOverview{}
	}
	if previousOverview == nil {
		previousOverview = &model.DashboardBookingOverview{}
	}

	return &response.DashboardBookingOverviewResponse{
		TotalRevenue:                  round2(currentOverview.TotalRevenue),
		TotalRevenueChangePercent:     float32(round2(calculateChangePercent(currentOverview.TotalRevenue, previousOverview.TotalRevenue))),
		OccupancyRate:                 float32(round2(currentOverview.OccupancyRate)),
		OccupancyRateChangePercent:    float32(round2(calculateChangePercent(currentOverview.OccupancyRate, previousOverview.OccupancyRate))),
		CancellationRate:              float32(round2(currentOverview.CancellationRate)),
		CancellationRateChangePercent: float32(round2(calculateChangePercent(currentOverview.CancellationRate, previousOverview.CancellationRate))),
		AverageDaysStay:               float32(round2(currentOverview.AverageDaysStay)),
		AverageDaysStayChangePercent:  float32(round2(calculateChangePercent(currentOverview.AverageDaysStay, previousOverview.AverageDaysStay))),
	}, nil
}

func (u *dashboardUseCaseImpl) BookingRevenue(ctx context.Context, schema string, userID int64, query request.DashboardRevenueQuery) (*response.DashboardBookingRevenueResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	revenue, err := u.bookingRepo.FindDashboardRevenue(ctx, schema, model.DashboardRevenueQuery{
		RoomTypeID: query.RoomTypeID,
		BranchIDs:  branchIDs,
	})
	if err != nil {
		u.logPro.Error("Find dashboard booking revenue failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if revenue == nil {
		return &response.DashboardBookingRevenueResponse{
			Granularity: string(model.DashboardRevenueGranularityDay),
			Chart:       make([]*response.DashboardBookingRevenuePointResponse, 0),
		}, nil
	}

	chart := make([]*response.DashboardBookingRevenuePointResponse, 0, len(revenue.Chart))
	for _, point := range revenue.Chart {
		if point == nil {
			continue
		}

		chart = append(chart, &response.DashboardBookingRevenuePointResponse{
			Date:    point.Date,
			Revenue: round2(point.Revenue),
		})
	}

	return &response.DashboardBookingRevenueResponse{
		Granularity: string(revenue.Granularity),
		FromDate:    revenue.FromDate,
		ToDate:      revenue.ToDate,
		Chart:       chart,
	}, nil
}

func (u *dashboardUseCaseImpl) RoomTypeRevenue(ctx context.Context, schema string, userID int64, query request.DashboardOverviewQuery) (*response.DashboardRoomTypeRevenueResponse, error) {
	fromDate, toDate, err := normalizeDashboardDateRange(query)
	if err != nil {
		return nil, err
	}
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	revenue, err := u.bookingRepo.FindDashboardRoomTypeRevenue(ctx, schema, model.DashboardOverviewQuery{
		FromDate:  fromDate,
		ToDate:    toDate,
		BranchIDs: branchIDs,
	})
	if err != nil {
		u.logPro.Error("Find dashboard room type revenue failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if revenue == nil {
		return &response.DashboardRoomTypeRevenueResponse{
			Chart: make([]*response.DashboardRoomTypeRevenuePointResponse, 0),
		}, nil
	}

	chart := make([]*response.DashboardRoomTypeRevenuePointResponse, 0, len(revenue.Chart))
	for _, point := range revenue.Chart {
		if point == nil {
			continue
		}

		chart = append(chart, &response.DashboardRoomTypeRevenuePointResponse{
			RoomTypeID:   point.RoomTypeID,
			RoomTypeName: point.RoomTypeName,
			Revenue:      round2(point.Revenue),
		})
	}

	return &response.DashboardRoomTypeRevenueResponse{
		Chart: chart,
	}, nil
}

func (u *dashboardUseCaseImpl) SourceRevenue(ctx context.Context, schema string, userID int64, query request.DashboardOverviewQuery) (*response.DashboardSourceRevenueResponse, error) {
	fromDate, toDate, err := normalizeDashboardDateRange(query)
	if err != nil {
		return nil, err
	}
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	revenue, err := u.bookingRepo.FindDashboardSourceRevenue(ctx, schema, model.DashboardOverviewQuery{
		FromDate:  fromDate,
		ToDate:    toDate,
		BranchIDs: branchIDs,
	})
	if err != nil {
		u.logPro.Error("Find dashboard source revenue failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if revenue == nil {
		return &response.DashboardSourceRevenueResponse{
			Chart: make([]*response.DashboardSourceRevenuePointResponse, 0),
		}, nil
	}

	chart := make([]*response.DashboardSourceRevenuePointResponse, 0, len(revenue.Chart))
	for _, point := range revenue.Chart {
		if point == nil {
			continue
		}

		chart = append(chart, &response.DashboardSourceRevenuePointResponse{
			SourceID:       point.SourceID,
			SourceName:     point.SourceName,
			Revenue:        round2(point.Revenue),
			RevenuePercent: float32(round2(point.RevenuePercent)),
		})
	}

	return &response.DashboardSourceRevenueResponse{
		Chart: chart,
	}, nil
}

func (u *dashboardUseCaseImpl) HotelOperation(ctx context.Context, schema string, userID int64) (*response.DashboardHotelOperationResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}
	todayStart := utils.NormalizeDateOnly(time.Now().UTC())
	tomorrowStart := todayStart.AddDate(0, 0, 1)

	var (
		todayCheckIn      int64
		todayCheckOut     int64
		pendingRequest    int64
		requestInProgress int64
		unreadMessage     int64
	)

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		todayCheckIn, todayCheckOut, err = u.bookingRepo.CountDashboardTodayCheckInOut(groupCtx, schema, todayStart, tomorrowStart, branchIDs)
		if err != nil {
			u.logPro.Error("Count dashboard today check in/out failed", port.String("schema", schema), port.Error(err))
		}

		return err
	})

	group.Go(func() error {
		var err error
		pendingRequest, requestInProgress, err = u.supportTicketRepo.CountDashboardPendingAndInProgress(groupCtx, schema, branchIDs)
		if err != nil {
			u.logPro.Error("Count dashboard support ticket operation failed", port.String("schema", schema), port.Error(err))
		}

		return err
	})

	group.Go(func() error {
		var err error
		if len(branchIDs) == 0 {
			unreadMessage, err = u.messageRepo.CountDashboardUnreadGuestMessages(groupCtx, schema)
		} else {
			chatIDs, chatErr := u.chatRepo.FindIDsByBranchIDs(groupCtx, schema, branchIDs)
			if chatErr != nil {
				err = chatErr
			} else if len(chatIDs) == 0 {
				unreadMessage = 0
			} else {
				unreadMessage, err = u.messageRepo.CountDashboardUnreadGuestMessages(groupCtx, schema, chatIDs)
			}
		}
		if err != nil {
			u.logPro.Error("Count dashboard unread guest messages failed", port.String("schema", schema), port.Error(err))
		}

		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &response.DashboardHotelOperationResponse{
		TodayCheckIn:      todayCheckIn,
		TodayCheckOut:     todayCheckOut,
		PendingRequest:    pendingRequest,
		RequestInProgress: requestInProgress,
		UnreadMessage:     unreadMessage,
	}, nil
}

func (u *dashboardUseCaseImpl) GuestNational(ctx context.Context, schema string, userID int64) (*response.DashboardGuestNationalResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	guestNational, err := u.guestRepo.FindDashboardGuestNational(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard guest national failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if guestNational == nil {
		return &response.DashboardGuestNationalResponse{
			Chart: make([]*response.DashboardGuestNationalPointResponse, 0),
		}, nil
	}

	chart := make([]*response.DashboardGuestNationalPointResponse, 0, len(guestNational.Chart))
	for _, point := range guestNational.Chart {
		if point == nil {
			continue
		}

		chart = append(chart, &response.DashboardGuestNationalPointResponse{
			National:   point.National,
			GuestCount: point.GuestCount,
		})
	}

	return &response.DashboardGuestNationalResponse{
		Chart: chart,
	}, nil
}

func (u *dashboardUseCaseImpl) SupportTicketTime(ctx context.Context, schema string, userID int64) (*response.DashboardSupportTicketTimeResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	supportTicketTime, err := u.supportTicketRepo.FindDashboardSupportTicketTime(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard support ticket time failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if supportTicketTime == nil {
		supportTicketTime = &model.DashboardSupportTicketTime{}
	}

	totalTickets := supportTicketTime.UnderOrEqual5Minutes +
		supportTicketTime.UnderOrEqual10Minutes +
		supportTicketTime.UnderOrEqual15Minutes +
		supportTicketTime.UnderOrEqual30Minutes +
		supportTicketTime.Over30Minutes

	chart := []*response.DashboardSupportTicketTimePointResponse{
		{
			Tag:           "<=5m",
			TicketCount:   supportTicketTime.UnderOrEqual5Minutes,
			TicketPercent: float32(round2(calculatePartPercent(supportTicketTime.UnderOrEqual5Minutes, totalTickets))),
		},
		{
			Tag:           "<=10m",
			TicketCount:   supportTicketTime.UnderOrEqual10Minutes,
			TicketPercent: float32(round2(calculatePartPercent(supportTicketTime.UnderOrEqual10Minutes, totalTickets))),
		},
		{
			Tag:           "<=15m",
			TicketCount:   supportTicketTime.UnderOrEqual15Minutes,
			TicketPercent: float32(round2(calculatePartPercent(supportTicketTime.UnderOrEqual15Minutes, totalTickets))),
		},
		{
			Tag:           "<=30m",
			TicketCount:   supportTicketTime.UnderOrEqual30Minutes,
			TicketPercent: float32(round2(calculatePartPercent(supportTicketTime.UnderOrEqual30Minutes, totalTickets))),
		},
		{
			Tag:           ">30m",
			TicketCount:   supportTicketTime.Over30Minutes,
			TicketPercent: float32(round2(calculatePartPercent(supportTicketTime.Over30Minutes, totalTickets))),
		},
	}

	return &response.DashboardSupportTicketTimeResponse{
		AverageResolutionMinutes: round2(supportTicketTime.AverageResolutionSeconds / 60),
		Chart:                    chart,
	}, nil
}

func (u *dashboardUseCaseImpl) ServiceTypeCount(ctx context.Context, schema string, userID int64) (*response.DashboardServiceTypeCountResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	serviceTypeCount, err := u.orderServiceRepo.FindDashboardServiceTypeCount(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard service type count failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if serviceTypeCount == nil {
		return &response.DashboardServiceTypeCountResponse{
			Chart: make([]*response.DashboardServiceTypeCountPointResponse, 0),
		}, nil
	}

	chart := make([]*response.DashboardServiceTypeCountPointResponse, 0, len(serviceTypeCount.Chart))
	for _, point := range serviceTypeCount.Chart {
		if point == nil {
			continue
		}

		chart = append(chart, &response.DashboardServiceTypeCountPointResponse{
			ServiceTypeID:   point.ServiceTypeID,
			ServiceTypeName: point.ServiceTypeName,
			OrderCount:      point.OrderCount,
			OrderPercent:    float32(round2(point.OrderPercent)),
		})
	}

	return &response.DashboardServiceTypeCountResponse{
		Chart: chart,
	}, nil
}

func (u *dashboardUseCaseImpl) ReviewSummary(ctx context.Context, schema string, userID int64) (*response.DashboardReviewSummaryResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	reviewSummary, err := u.reviewRepo.FindDashboardReviewSummary(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard review summary failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if reviewSummary == nil {
		reviewSummary = &model.DashboardReviewSummary{}
	}

	totalSentiment := reviewSummary.PositiveCount + reviewSummary.NegativeCount + reviewSummary.NeutralCount

	return &response.DashboardReviewSummaryResponse{
		ReviewCount:  reviewSummary.ReviewCount,
		AverageStar:  float32(round2(reviewSummary.AverageStar)),
		PositiveRate: float32(round2(calculatePartPercent(reviewSummary.PositiveCount, totalSentiment))),
		NegativeRate: float32(round2(calculatePartPercent(reviewSummary.NegativeCount, totalSentiment))),
	}, nil
}

func (u *dashboardUseCaseImpl) ReviewCategoryCount(ctx context.Context, schema string, userID int64) (*response.DashboardReviewCategoryCountResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	categoryCount, err := u.reviewRepo.FindDashboardReviewCategoryCount(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard review category count failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if categoryCount == nil {
		categoryCount = &model.DashboardReviewCategoryCount{
			Chart: make([]*model.DashboardReviewCategoryCountPoint, 0),
		}
	}

	categoryMap := make(map[string]int64, len(dashboardReviewCategories))
	for _, category := range dashboardReviewCategories {
		categoryMap[category] = 0
	}
	for _, point := range categoryCount.Chart {
		if point == nil {
			continue
		}

		category := strings.ToLower(strings.TrimSpace(point.Category))
		if _, ok := categoryMap[category]; !ok {
			continue
		}

		categoryMap[category] = point.ReviewCount
	}

	chart := make([]*response.DashboardReviewCategoryCountPointResponse, 0, len(dashboardReviewCategories))
	for _, category := range dashboardReviewCategories {
		chart = append(chart, &response.DashboardReviewCategoryCountPointResponse{
			Category:    category,
			ReviewCount: categoryMap[category],
		})
	}

	return &response.DashboardReviewCategoryCountResponse{
		Chart: chart,
	}, nil
}

func (u *dashboardUseCaseImpl) ReviewCategorySentiment(ctx context.Context, schema string, userID int64) (*response.DashboardReviewCategorySentimentResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	categorySentiment, err := u.reviewRepo.FindDashboardReviewCategorySentiment(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard review category sentiment failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if categorySentiment == nil {
		categorySentiment = &model.DashboardReviewCategorySentiment{
			Chart: make([]*model.DashboardReviewCategorySentimentPoint, 0),
		}
	}

	type sentimentSummary struct {
		Positive int64
		Negative int64
		Neutral  int64
	}

	categoryMap := make(map[string]sentimentSummary, len(dashboardReviewCategories))
	for _, category := range dashboardReviewCategories {
		categoryMap[category] = sentimentSummary{}
	}
	for _, point := range categorySentiment.Chart {
		if point == nil {
			continue
		}

		category := strings.ToLower(strings.TrimSpace(point.Category))
		if _, ok := categoryMap[category]; !ok {
			continue
		}

		categoryMap[category] = sentimentSummary{
			Positive: point.PositiveCount,
			Negative: point.NegativeCount,
			Neutral:  point.NeutralCount,
		}
	}

	chart := make([]*response.DashboardReviewCategorySentimentPointResponse, 0, len(dashboardReviewCategories))
	for _, category := range dashboardReviewCategories {
		summary := categoryMap[category]
		total := summary.Positive + summary.Negative + summary.Neutral

		chart = append(chart, &response.DashboardReviewCategorySentimentPointResponse{
			Category: category,
			Positive: buildReviewCategorySentimentBreakdown(summary.Positive, total),
			Negative: buildReviewCategorySentimentBreakdown(summary.Negative, total),
			Neutral:  buildReviewCategorySentimentBreakdown(summary.Neutral, total),
		})
	}

	return &response.DashboardReviewCategorySentimentResponse{
		Chart: chart,
	}, nil
}

func (u *dashboardUseCaseImpl) ReviewSentimentCount(ctx context.Context, schema string, userID int64) (*response.DashboardReviewSentimentCountResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	sentimentCount, err := u.reviewRepo.FindDashboardReviewSentimentCount(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard review sentiment count failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if sentimentCount == nil {
		sentimentCount = &model.DashboardReviewSentimentCount{}
	}

	totalSentiment := sentimentCount.PositiveCount + sentimentCount.NegativeCount + sentimentCount.NeutralCount

	return &response.DashboardReviewSentimentCountResponse{
		Positive: buildReviewSentimentBreakdown(sentimentCount.PositiveCount, totalSentiment),
		Negative: buildReviewSentimentBreakdown(sentimentCount.NegativeCount, totalSentiment),
		Neutral:  buildReviewSentimentBreakdown(sentimentCount.NeutralCount, totalSentiment),
	}, nil
}

func (u *dashboardUseCaseImpl) ReviewSentimentTime(ctx context.Context, schema string, userID int64) (*response.DashboardReviewSentimentTimeResponse, error) {
	branchIDs, err := u.dashboardBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	sentimentTime, err := u.reviewRepo.FindDashboardReviewSentimentTime(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Find dashboard review sentiment time failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if sentimentTime == nil {
		return &response.DashboardReviewSentimentTimeResponse{
			Granularity: string(model.DashboardRevenueGranularityDay),
			Chart:       make([]*response.DashboardReviewSentimentTimePointResponse, 0),
		}, nil
	}

	chart := make([]*response.DashboardReviewSentimentTimePointResponse, 0, len(sentimentTime.Chart))
	for _, point := range sentimentTime.Chart {
		if point == nil {
			continue
		}

		totalSentiment := point.PositiveCount + point.NegativeCount + point.NeutralCount
		chart = append(chart, &response.DashboardReviewSentimentTimePointResponse{
			Date:     point.Date,
			Positive: buildReviewSentimentBreakdown(point.PositiveCount, totalSentiment),
			Negative: buildReviewSentimentBreakdown(point.NegativeCount, totalSentiment),
			Neutral:  buildReviewSentimentBreakdown(point.NeutralCount, totalSentiment),
		})
	}

	return &response.DashboardReviewSentimentTimeResponse{
		Granularity: string(sentimentTime.Granularity),
		FromDate:    sentimentTime.FromDate,
		ToDate:      sentimentTime.ToDate,
		Chart:       chart,
	}, nil
}

func normalizeDashboardDateRange(query request.DashboardOverviewQuery) (time.Time, time.Time, error) {
	var fromDate time.Time
	var toDate time.Time

	switch {
	case query.FromDate == nil && query.ToDate == nil:
		now := time.Now().UTC()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		toDate = fromDate
	case query.FromDate == nil && query.ToDate != nil:
		toDate = utils.NormalizeDateOnly(*query.ToDate)
		fromDate = toDate
	case query.FromDate != nil && query.ToDate == nil:
		fromDate = utils.NormalizeDateOnly(*query.FromDate)
		toDate = fromDate
	default:
		fromDate = utils.NormalizeDateOnly(*query.FromDate)
		toDate = utils.NormalizeDateOnly(*query.ToDate)
	}

	if toDate.Before(fromDate) {
		return time.Time{}, time.Time{}, exceptions.ErrBadRequest.WithData(map[string]string{
			"field": "to_date",
			"tag":   "gtefield",
			"param": "from_date",
		})
	}

	return fromDate, toDate, nil
}

func (u *dashboardUseCaseImpl) dashboardBranchIDs(ctx context.Context, schema string, userID int64) ([]int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}
	if scope.IsGlobal() {
		return nil, nil
	}

	return scope.BranchIDs, nil
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func previousDashboardDateRange(fromDate, toDate time.Time) (time.Time, time.Time) {
	dayCount := int(toDate.Sub(fromDate)/(24*time.Hour)) + 1
	previousToDate := fromDate.AddDate(0, 0, -1)
	previousFromDate := previousToDate.AddDate(0, 0, -(dayCount - 1))

	return previousFromDate, previousToDate
}

func calculateChangePercent(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}

		return 100
	}

	return ((current - previous) / previous) * 100
}

func calculatePartPercent(part, total int64) float64 {
	if total <= 0 || part <= 0 {
		return 0
	}

	return (float64(part) / float64(total)) * 100
}

func buildReviewCategorySentimentBreakdown(count, total int64) *response.DashboardReviewCategorySentimentBreakdownResponse {
	return &response.DashboardReviewCategorySentimentBreakdownResponse{
		Count:   count,
		Percent: float32(round2(calculatePartPercent(count, total))),
	}
}

func buildReviewSentimentBreakdown(count, total int64) *response.DashboardReviewSentimentCountBreakdownResponse {
	return &response.DashboardReviewSentimentCountBreakdownResponse{
		Count:   count,
		Percent: float32(round2(calculatePartPercent(count, total))),
	}
}
