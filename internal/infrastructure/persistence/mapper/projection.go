package mapper

import (
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type BookingStatusSummaryProjection struct {
	All        int64 `gorm:"column:all_count"`
	CheckedIn  int64 `gorm:"column:checked_in_count"`
	Confirmed  int64 `gorm:"column:confirmed_count"`
	CheckedOut int64 `gorm:"column:checked_out_count"`
	Cancelled  int64 `gorm:"column:cancelled_count"`
}

type BookingDashboardOverviewProjection struct {
	TotalRevenue          float64 `gorm:"column:total_revenue"`
	SoldRoomNights        float64 `gorm:"column:sold_room_nights"`
	TotalBookingCount     int64   `gorm:"column:total_booking_count"`
	CancelledBookingCount int64   `gorm:"column:cancelled_booking_count"`
	AverageDaysStay       float64 `gorm:"column:average_days_stay"`
	TotalRooms            int64   `gorm:"column:total_rooms"`
}

type BookingDashboardRevenuePointProjection struct {
	BucketStart time.Time `gorm:"column:bucket_start"`
	Revenue     float64   `gorm:"column:revenue"`
}

type BookingDashboardRoomTypeRevenueProjection struct {
	RoomTypeID   int64   `gorm:"column:room_type_id"`
	RoomTypeName string  `gorm:"column:room_type_name"`
	Revenue      float64 `gorm:"column:revenue"`
}

type BookingDashboardSourceRevenueProjection struct {
	SourceID       int64   `gorm:"column:source_id"`
	SourceName     string  `gorm:"column:source_name"`
	Revenue        float64 `gorm:"column:revenue"`
	RevenuePercent float64 `gorm:"column:revenue_percent"`
}

type BookingTimelineProjection struct {
	RoomTypeID    int64      `gorm:"column:room_type_id"`
	RoomTypeName  string     `gorm:"column:room_type_name"`
	RoomID        int64      `gorm:"column:room_id"`
	RoomName      string     `gorm:"column:room_name"`
	BookingID     *int64     `gorm:"column:booking_id"`
	CheckIn       *time.Time `gorm:"column:check_in"`
	CheckOut      *time.Time `gorm:"column:check_out"`
	GuestName     *string    `gorm:"column:guest_name"`
	SourceName    *string    `gorm:"column:source_name"`
	Status        *string    `gorm:"column:status"`
	BookingNumber *string    `gorm:"column:booking_number"`
}

type RoomStatusSummaryProjection struct {
	All       int64 `gorm:"column:all_count"`
	Available int64 `gorm:"column:available_count"`
	Occupied  int64 `gorm:"column:occupied_count"`
	Booked    int64 `gorm:"column:booked_count"`
}

type RoomTimelineProjection struct {
	RoomTypeID    int64      `gorm:"column:room_type_id"`
	RoomTypeName  string     `gorm:"column:room_type_name"`
	RoomPrice     float64    `gorm:"column:room_price"`
	RoomID        int64      `gorm:"column:room_id"`
	RoomName      string     `gorm:"column:room_name"`
	RoomStatus    string     `gorm:"column:room_status"`
	BookingID     *int64     `gorm:"column:booking_id"`
	BookingNumber *string    `gorm:"column:booking_number"`
	BookingStatus *string    `gorm:"column:booking_status"`
	SourceName    *string    `gorm:"column:source_name"`
	GuestName     *string    `gorm:"column:guest_name"`
	CheckIn       *time.Time `gorm:"column:check_in"`
	CheckOut      *time.Time `gorm:"column:check_out"`
}

type SupportTicketDashboardTimeProjection struct {
	UnderOrEqual5Minutes    int64   `gorm:"column:under_or_equal_5_minutes"`
	UnderOrEqual10Minutes   int64   `gorm:"column:under_or_equal_10_minutes"`
	UnderOrEqual15Minutes   int64   `gorm:"column:under_or_equal_15_minutes"`
	UnderOrEqual30Minutes   int64   `gorm:"column:under_or_equal_30_minutes"`
	Over30Minutes           int64   `gorm:"column:over_30_minutes"`
	AverageResolutionSecond float64 `gorm:"column:average_resolution_seconds"`
}

type OrderServiceDashboardServiceTypeCountProjection struct {
	ServiceTypeID   int64   `gorm:"column:service_type_id"`
	ServiceTypeName string  `gorm:"column:service_type_name"`
	OrderCount      int64   `gorm:"column:order_count"`
	OrderPercent    float64 `gorm:"column:order_percent"`
}

type GuestDashboardNationalProjection struct {
	National   string `gorm:"column:national"`
	GuestCount int64  `gorm:"column:guest_count"`
}

type ReviewStarSummaryProjection struct {
	All       int64 `gorm:"column:all_count"`
	FiveStar  int64 `gorm:"column:five_star_count"`
	FourStar  int64 `gorm:"column:four_star_count"`
	ThreeStar int64 `gorm:"column:three_star_count"`
	TwoStar   int64 `gorm:"column:two_star_count"`
	OneStar   int64 `gorm:"column:one_star_count"`
}

type ReviewDashboardSummaryProjection struct {
	ReviewCount   int64   `gorm:"column:review_count"`
	AverageStar   float64 `gorm:"column:average_star"`
	PositiveCount int64   `gorm:"column:positive_count"`
	NegativeCount int64   `gorm:"column:negative_count"`
	NeutralCount  int64   `gorm:"column:neutral_count"`
}

type ReviewDashboardCategoryCountProjection struct {
	Category    string `gorm:"column:category"`
	ReviewCount int64  `gorm:"column:review_count"`
}

type ReviewDashboardCategorySentimentProjection struct {
	Category      string `gorm:"column:category"`
	PositiveCount int64  `gorm:"column:positive_count"`
	NegativeCount int64  `gorm:"column:negative_count"`
	NeutralCount  int64  `gorm:"column:neutral_count"`
}

type ReviewDashboardSentimentCountProjection struct {
	PositiveCount int64 `gorm:"column:positive_count"`
	NegativeCount int64 `gorm:"column:negative_count"`
	NeutralCount  int64 `gorm:"column:neutral_count"`
}

type ReviewDashboardSentimentTimeRangeProjection struct {
	MinDate *time.Time `gorm:"column:min_date"`
	MaxDate *time.Time `gorm:"column:max_date"`
}

type ReviewDashboardSentimentTimePointProjection struct {
	BucketStart   time.Time `gorm:"column:bucket_start"`
	PositiveCount int64     `gorm:"column:positive_count"`
	NegativeCount int64     `gorm:"column:negative_count"`
	NeutralCount  int64     `gorm:"column:neutral_count"`
}

func ToBookingStatusSummaryModel(summary BookingStatusSummaryProjection) *model.BookingStatusSummary {
	return &model.BookingStatusSummary{
		All:        summary.All,
		CheckedIn:  summary.CheckedIn,
		Confirmed:  summary.Confirmed,
		CheckedOut: summary.CheckedOut,
		Cancelled:  summary.Cancelled,
	}
}

func ToDashboardBookingOverviewModel(summary BookingDashboardOverviewProjection, occupancyRate, cancellationRate float64) *model.DashboardBookingOverview {
	return &model.DashboardBookingOverview{
		TotalRevenue:     summary.TotalRevenue,
		OccupancyRate:    occupancyRate,
		CancellationRate: cancellationRate,
		AverageDaysStay:  summary.AverageDaysStay,
	}
}

func ToDashboardBookingRevenuePoints(points []*BookingDashboardRevenuePointProjection) []*model.DashboardBookingRevenuePoint {
	if len(points) == 0 {
		return make([]*model.DashboardBookingRevenuePoint, 0)
	}

	chart := make([]*model.DashboardBookingRevenuePoint, 0, len(points))
	for _, point := range points {
		if point == nil {
			continue
		}

		chart = append(chart, &model.DashboardBookingRevenuePoint{
			Date:    utils.NormalizeDateOnly(point.BucketStart),
			Revenue: point.Revenue,
		})
	}

	return chart
}

func ToDashboardRoomTypeRevenueModel(rows []*BookingDashboardRoomTypeRevenueProjection) *model.DashboardRoomTypeRevenue {
	chart := make([]*model.DashboardRoomTypeRevenuePoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardRoomTypeRevenuePoint{
			RoomTypeID:   row.RoomTypeID,
			RoomTypeName: row.RoomTypeName,
			Revenue:      row.Revenue,
		})
	}

	return &model.DashboardRoomTypeRevenue{
		Chart: chart,
	}
}

func ToDashboardSourceRevenueModel(rows []*BookingDashboardSourceRevenueProjection) *model.DashboardSourceRevenue {
	chart := make([]*model.DashboardSourceRevenuePoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardSourceRevenuePoint{
			SourceID:       row.SourceID,
			SourceName:     row.SourceName,
			Revenue:        row.Revenue,
			RevenuePercent: row.RevenuePercent,
		})
	}

	return &model.DashboardSourceRevenue{
		Chart: chart,
	}
}

func ToBookingTimelineModel(startDate, endDate time.Time, rows []*BookingTimelineProjection) *model.BookingTimeline {
	timeline := &model.BookingTimeline{
		StartDate: startDate,
		EndDate:   endDate,
		Data:      make([]*model.BookingTimelineRoomType, 0),
	}

	if len(rows) == 0 {
		return timeline
	}

	roomTypeMap := make(map[int64]*model.BookingTimelineRoomType)
	roomMap := make(map[int64]map[int64]*model.BookingTimelineRoom)

	for _, row := range rows {
		if row == nil {
			continue
		}

		roomType, exists := roomTypeMap[row.RoomTypeID]
		if !exists {
			roomType = &model.BookingTimelineRoomType{
				RoomTypeID:   row.RoomTypeID,
				RoomTypeName: row.RoomTypeName,
				Rooms:        make([]*model.BookingTimelineRoom, 0),
			}
			roomTypeMap[row.RoomTypeID] = roomType
			roomMap[row.RoomTypeID] = make(map[int64]*model.BookingTimelineRoom)
			timeline.Data = append(timeline.Data, roomType)
		}

		room, exists := roomMap[row.RoomTypeID][row.RoomID]
		if !exists {
			room = &model.BookingTimelineRoom{
				RoomID:   row.RoomID,
				RoomName: row.RoomName,
				Events:   make([]*model.BookingTimelineEvent, 0),
			}
			roomMap[row.RoomTypeID][row.RoomID] = room
			roomType.Rooms = append(roomType.Rooms, room)
		}

		if row.BookingID == nil || row.CheckIn == nil || row.CheckOut == nil || row.Status == nil {
			continue
		}

		sourceName := ""
		if row.SourceName != nil {
			sourceName = *row.SourceName
		}

		guestName := ""
		if row.GuestName != nil {
			guestName = *row.GuestName
		}

		bookingNumber := ""
		if row.BookingNumber != nil {
			bookingNumber = *row.BookingNumber
		}

		room.Events = append(room.Events, &model.BookingTimelineEvent{
			BookingID:     *row.BookingID,
			CheckIn:       *row.CheckIn,
			CheckOut:      *row.CheckOut,
			Title:         fmt.Sprintf("%s - %s", sourceName, guestName),
			Status:        model.BookingStatus(*row.Status),
			BookingNumber: bookingNumber,
		})
	}

	for _, roomType := range timeline.Data {
		roomType.RoomsCount = int64(len(roomType.Rooms))
	}

	return timeline
}

func ToRoomStatusSummaryModel(summary RoomStatusSummaryProjection) *model.RoomStatusSummary {
	return &model.RoomStatusSummary{
		All:       summary.All,
		Available: summary.Available,
		Occupied:  summary.Occupied,
		Booked:    summary.Booked,
	}
}

func ToRoomTimelineModel(serverNow time.Time, rows []*RoomTimelineProjection) *model.RoomTimeline {
	timeline := &model.RoomTimeline{
		ServerNow: serverNow,
		Status:    &model.RoomTimelineStatusSummary{},
		Data:      make([]*model.RoomTimelineRoomType, 0),
	}

	if len(rows) == 0 {
		return timeline
	}

	roomTypeMap := make(map[int64]*model.RoomTimelineRoomType)

	for _, row := range rows {
		if row == nil {
			continue
		}

		timeline.Status.All++
		switch row.RoomStatus {
		case "occupied":
			timeline.Status.Occupied++
		case "upcoming_checkin":
			timeline.Status.UpcomingCheckIn++
		case "overdue_checkout":
			timeline.Status.OverdueCheckout++
		case "overdue_checkin":
			timeline.Status.OverdueCheckIn++
		default:
			timeline.Status.Available++
		}

		roomType, exists := roomTypeMap[row.RoomTypeID]
		if !exists {
			roomType = &model.RoomTimelineRoomType{
				RoomTypeID:   row.RoomTypeID,
				RoomTypeName: row.RoomTypeName,
				Rooms:        make([]*model.RoomTimelineRoom, 0),
			}
			roomTypeMap[row.RoomTypeID] = roomType
			timeline.Data = append(timeline.Data, roomType)
		}

		room := &model.RoomTimelineRoom{
			RoomID:   row.RoomID,
			RoomName: row.RoomName,
			Price:    row.RoomPrice,
		}

		if row.BookingID != nil && row.BookingNumber != nil && row.BookingStatus != nil && row.CheckIn != nil && row.CheckOut != nil {
			sourceName := ""
			if row.SourceName != nil {
				sourceName = *row.SourceName
			}

			guestName := ""
			if row.GuestName != nil {
				guestName = *row.GuestName
			}

			room.Booking = &model.RoomTimelineBooking{
				BookingID:     *row.BookingID,
				BookingNumber: *row.BookingNumber,
				Status:        model.BookingStatus(*row.BookingStatus),
				SourceName:    sourceName,
				GuestName:     guestName,
				CheckIn:       *row.CheckIn,
				CheckOut:      *row.CheckOut,
			}
		}

		roomType.Rooms = append(roomType.Rooms, room)
	}

	for _, roomType := range timeline.Data {
		roomType.RoomsCount = int64(len(roomType.Rooms))
	}

	return timeline
}

func ToDashboardSupportTicketTimeModel(result SupportTicketDashboardTimeProjection) *model.DashboardSupportTicketTime {
	return &model.DashboardSupportTicketTime{
		UnderOrEqual5Minutes:     result.UnderOrEqual5Minutes,
		UnderOrEqual10Minutes:    result.UnderOrEqual10Minutes,
		UnderOrEqual15Minutes:    result.UnderOrEqual15Minutes,
		UnderOrEqual30Minutes:    result.UnderOrEqual30Minutes,
		Over30Minutes:            result.Over30Minutes,
		AverageResolutionSeconds: result.AverageResolutionSecond,
	}
}

func ToDashboardServiceTypeCountModel(rows []*OrderServiceDashboardServiceTypeCountProjection) *model.DashboardServiceTypeCount {
	chart := make([]*model.DashboardServiceTypeCountPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardServiceTypeCountPoint{
			ServiceTypeID:   row.ServiceTypeID,
			ServiceTypeName: row.ServiceTypeName,
			OrderCount:      row.OrderCount,
			OrderPercent:    row.OrderPercent,
		})
	}

	return &model.DashboardServiceTypeCount{
		Chart: chart,
	}
}

func ToDashboardGuestNationalModel(rows []*GuestDashboardNationalProjection) *model.DashboardGuestNational {
	chart := make([]*model.DashboardGuestNationalPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardGuestNationalPoint{
			National:   row.National,
			GuestCount: row.GuestCount,
		})
	}

	return &model.DashboardGuestNational{
		Chart: chart,
	}
}

func ToReviewStarSummaryModel(summary ReviewStarSummaryProjection) *model.ReviewStarSummary {
	return &model.ReviewStarSummary{
		All:       summary.All,
		FiveStar:  summary.FiveStar,
		FourStar:  summary.FourStar,
		ThreeStar: summary.ThreeStar,
		TwoStar:   summary.TwoStar,
		OneStar:   summary.OneStar,
	}
}

func ToDashboardReviewSummaryModel(summary ReviewDashboardSummaryProjection) *model.DashboardReviewSummary {
	return &model.DashboardReviewSummary{
		ReviewCount:   summary.ReviewCount,
		AverageStar:   summary.AverageStar,
		PositiveCount: summary.PositiveCount,
		NegativeCount: summary.NegativeCount,
		NeutralCount:  summary.NeutralCount,
	}
}

func ToDashboardReviewCategoryCountModel(rows []*ReviewDashboardCategoryCountProjection) *model.DashboardReviewCategoryCount {
	chart := make([]*model.DashboardReviewCategoryCountPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardReviewCategoryCountPoint{
			Category:    row.Category,
			ReviewCount: row.ReviewCount,
		})
	}

	return &model.DashboardReviewCategoryCount{
		Chart: chart,
	}
}

func ToDashboardReviewCategorySentimentModel(rows []*ReviewDashboardCategorySentimentProjection) *model.DashboardReviewCategorySentiment {
	chart := make([]*model.DashboardReviewCategorySentimentPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardReviewCategorySentimentPoint{
			Category:      row.Category,
			PositiveCount: row.PositiveCount,
			NegativeCount: row.NegativeCount,
			NeutralCount:  row.NeutralCount,
		})
	}

	return &model.DashboardReviewCategorySentiment{
		Chart: chart,
	}
}

func ToDashboardReviewSentimentCountModel(data ReviewDashboardSentimentCountProjection) *model.DashboardReviewSentimentCount {
	return &model.DashboardReviewSentimentCount{
		PositiveCount: data.PositiveCount,
		NegativeCount: data.NegativeCount,
		NeutralCount:  data.NeutralCount,
	}
}

func ToDashboardReviewSentimentTimePoints(rows []*ReviewDashboardSentimentTimePointProjection) []*model.DashboardReviewSentimentTimePoint {
	chart := make([]*model.DashboardReviewSentimentTimePoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}

		chart = append(chart, &model.DashboardReviewSentimentTimePoint{
			Date:          utils.NormalizeDateOnly(row.BucketStart),
			PositiveCount: row.PositiveCount,
			NegativeCount: row.NegativeCount,
			NeutralCount:  row.NeutralCount,
		})
	}

	return chart
}
