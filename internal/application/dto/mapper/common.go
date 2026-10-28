package mapper

import (
	"math"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

func ToMetaResponse(total int64, page, limit uint32) *response.MetaResponse {
	if total < 0 {
		total = 0
	}
	totalPages := uint32(0)
	if total > 0 {
		t := uint64(total)
		l := uint64(limit)
		pages := t / l
		if t%l != 0 {
			pages++
		}
		if pages > math.MaxUint16 {
			pages = math.MaxUint16
		}
		totalPages = uint32(pages)
	}

	return &response.MetaResponse{
		Total:      uint64(total),
		Page:       page,
		Limit:      limit,
		TotalPages: uint16(totalPages),
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
	}
}
