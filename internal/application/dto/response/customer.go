package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type CustomerResponse struct {
	ID        int64                `json:"id,string"`
	BirthYear uint32               `json:"birth_year"`
	Gender    model.CustomerGender `json:"gender"`
	Phone     string               `json:"phone"`
	CreatedAt time.Time            `json:"created_at"`
}
