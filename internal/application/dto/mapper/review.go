package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBasicReviewResponse(rv *model.Review) *response.BasicReviewResponse {
	if rv == nil {
		return nil
	}

	return &response.BasicReviewResponse{
		ID:      rv.ID,
		Star:    rv.Star,
		Content: rv.Content,
	}
}

func ToReviewResponse(rv *model.Review) *response.ReviewResponse {
	if rv == nil {
		return nil
	}

	return &response.ReviewResponse{
		ID:          rv.ID,
		Star:        rv.Star,
		Content:     rv.Content,
		CreatedAt:   rv.CreatedAt,
		BookingRoom: ToSimpleBookingRoomResponse(rv.BookingRoom),
	}
}

func ToReviewsResponse(rvs []*model.Review) []*response.ReviewResponse {
	if len(rvs) == 0 {
		return make([]*response.ReviewResponse, 0)
	}

	reviewsRes := make([]*response.ReviewResponse, 0, len(rvs))
	for _, rv := range rvs {
		reviewsRes = append(reviewsRes, ToReviewResponse(rv))
	}

	return reviewsRes
}

func ToSimpleReviewsResponse(rvs []*model.Review) []*response.SimpleReviewResponse {
	if len(rvs) == 0 {
		return make([]*response.SimpleReviewResponse, 0)
	}

	reviewsRes := make([]*response.SimpleReviewResponse, 0, len(rvs))
	for _, rv := range rvs {
		reviewsRes = append(reviewsRes, ToSimpleReviewResponse(rv))
	}

	return reviewsRes
}

func ToReviewQuery(query request.ReviewPaginationQuery) model.ReviewQuery {
	return model.ReviewQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Sort:     query.Sort,
		Order:    query.Order,
		FromDate: query.FromDate,
		ToDate:   query.ToDate,
		Star:     query.Star,
	}
}
