package mapper

import (
	"encoding/json"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/datatypes"
)

func ToReviewEntity(rv *model.Review) *entity.ReviewEntity {
	if rv == nil {
		return nil
	}

	var aspectsJSON datatypes.JSON
	if len(rv.Aspects) > 0 {
		bytes, _ := json.Marshal(rv.Aspects)
		aspectsJSON = datatypes.JSON(bytes)
	}

	return &entity.ReviewEntity{
		BaseEntity: entity.BaseEntity{
			ID:        rv.ID,
			CreatedAt: rv.CreatedAt,
			UpdatedAt: rv.UpdatedAt,
		},
		Star:          rv.Star,
		Content:       rv.Content,
		BookingRoomID: rv.BookingRoomID,
		Aspects:       aspectsJSON,
	}
}

func ToReview(rv *entity.ReviewEntity) *model.Review {
	if rv == nil {
		return nil
	}

	var aspects []model.ReviewAspect
	if len(rv.Aspects) > 0 {
		_ = json.Unmarshal(rv.Aspects, &aspects)
	}

	return &model.Review{
		Base: model.Base{
			ID:        rv.ID,
			CreatedAt: rv.CreatedAt,
			UpdatedAt: rv.UpdatedAt,
		},
		Star:          rv.Star,
		Content:       rv.Content,
		BookingRoomID: rv.BookingRoomID,
		Aspects:       aspects,
		BookingRoom:   ToBookingRoom(rv.BookingRoom),
	}
}

func ToReviews(rvs []*entity.ReviewEntity) []*model.Review {
	if len(rvs) == 0 {
		return make([]*model.Review, 0)
	}

	reviews := make([]*model.Review, 0, len(rvs))
	for _, rv := range rvs {
		reviews = append(reviews, ToReview(rv))
	}

	return reviews
}
