package review

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"gorm.io/datatypes"
)

type reviewUseCaseImpl struct {
	logPro     port.LogProvider
	idPro      port.IDProvider
	aiPro      port.AIProvider
	reviewRepo repository.ReviewRepository
	userRepo   repository.UserRepository
}

func NewReviewUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	aiPro port.AIProvider,
	reviewRepo repository.ReviewRepository,
	userRepo repository.UserRepository,
) ReviewUseCase {
	return &reviewUseCaseImpl{
		logPro,
		idPro,
		aiPro,
		reviewRepo,
		userRepo,
	}
}

func (u *reviewUseCaseImpl) Create(ctx context.Context, schema string, bookingRoomID int64, req request.CreateReviewRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate review id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	review := model.NewReview(id, bookingRoomID, req.Star, req.Content)

	if err = u.reviewRepo.Create(ctx, schema, review); err != nil {
		if errors.Is(err, exceptions.ErrReviewAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrUnauth) {
			return err
		}
		u.logPro.Error("Create review failed", port.String("schema", schema), port.Error(err))
		return err
	}

	u.analyzeReview(schema, id, req.Content)

	return nil
}

func (u *reviewUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID int64, query request.ReviewPaginationQuery) ([]*response.ReviewResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	reviewQuery, err := u.scopedReviewQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	reviews, total, err := u.reviewRepo.FindAllWithBookingRoomDetailsPaginated(ctx, schema, reviewQuery)
	if err != nil {
		u.logPro.Error("Find all review paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToReviewsResponse(reviews), meta, nil
}

func (u *reviewUseCaseImpl) ListForGuest(ctx context.Context, schema string, query request.ReviewPaginationQuery) ([]*response.SimpleReviewResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	reviews, total, err := u.reviewRepo.FindAllPaginated(ctx, schema, mapper.ToReviewQuery(query))
	if err != nil {
		u.logPro.Error("Find all review paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSimpleReviewsResponse(reviews), meta, nil
}

func (u *reviewUseCaseImpl) GetByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) (*response.BasicReviewResponse, error) {
	review, err := u.reviewRepo.FindByBookingRoomID(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find review by booking room id failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicReviewResponse(review), nil
}

func (u *reviewUseCaseImpl) Summary(ctx context.Context, schema string, userID int64) (*response.ReviewSummaryResponse, error) {
	branchIDs, err := u.reviewBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}

	summary, err := u.reviewRepo.CountStarSummary(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Count review star summary failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return &response.ReviewSummaryResponse{
		All:              summary.All,
		Average:          calculateReviewAverage(summary),
		FiveStarPercent:  calculateReviewStarPercent(summary.FiveStar, summary.All),
		FourStarPercent:  calculateReviewStarPercent(summary.FourStar, summary.All),
		ThreeStarPercent: calculateReviewStarPercent(summary.ThreeStar, summary.All),
		TwoStarPercent:   calculateReviewStarPercent(summary.TwoStar, summary.All),
		OneStarPercent:   calculateReviewStarPercent(summary.OneStar, summary.All),
	}, nil
}

func (u *reviewUseCaseImpl) UpdateByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64, req request.UpdateReviewRequest) error {
	review, err := u.reviewRepo.FindByBookingRoomID(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find review by booking room id failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if review == nil {
		return exceptions.ErrReviewNotFound
	}

	updateData := map[string]any{
		"star":    req.Star,
		"content": req.Content,
	}

	if err := u.reviewRepo.UpdateByBookingRoomID(ctx, schema, bookingRoomID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrReviewNotFound) {
			return err
		}
		u.logPro.Error("Update review by booking room failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if review.Content != req.Content {
		u.analyzeReview(schema, review.ID, req.Content)
	}

	return nil
}

func (u *reviewUseCaseImpl) scopedReviewQuery(ctx context.Context, schema string, userID int64, query request.ReviewPaginationQuery) (model.ReviewQuery, error) {
	branchIDs, err := u.reviewBranchIDs(ctx, schema, userID)
	if err != nil {
		return model.ReviewQuery{}, err
	}

	reviewQuery := mapper.ToReviewQuery(query)
	reviewQuery.BranchIDs = branchIDs

	return reviewQuery, nil
}

func (u *reviewUseCaseImpl) reviewBranchIDs(ctx context.Context, schema string, userID int64) ([]int64, error) {
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

func (u *reviewUseCaseImpl) analyzeReview(schema string, reviewID int64, reviewContent string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		reviewID, aspects, err := u.aiPro.AnalyzeReview(ctx, reviewID, reviewContent)
		if err != nil {
			u.logPro.Error("Analyze review from gRPC failed", port.String("schema", schema), port.Error(err))
			return
		}

		aspectsJSON, err := json.Marshal(aspects)
		if err != nil {
			u.logPro.Error("Marshal review aspects failed", port.String("schema", schema), port.Error(err))
			return
		}

		if err := u.reviewRepo.Update(ctx, schema, reviewID, map[string]any{
			"aspects": datatypes.JSON(aspectsJSON),
		}); err != nil {
			u.logPro.Error("Update review aspects failed", port.String("schema", schema), port.Error(err))
		}
	}()
}

func calculateReviewStarPercent(starCount, all int64) float32 {
	if all == 0 {
		return 0
	}

	percent := (float64(starCount) / float64(all)) * 100
	percent = math.Round(percent*100) / 100

	return float32(percent)
}

func calculateReviewAverage(summary *model.ReviewStarSummary) float32 {
	if summary == nil || summary.All == 0 {
		return 0
	}

	totalStar := summary.FiveStar*5 + summary.FourStar*4 + summary.ThreeStar*3 + summary.TwoStar*2 + summary.OneStar
	avg := float64(totalStar) / float64(summary.All)
	avg = math.Round(avg*100) / 100

	return float32(avg)
}
