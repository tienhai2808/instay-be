package port

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type TokenProvider interface {
	GenerateStaffToken(userID int64, role model.UserRole, tokenVersion int, departmentID *int64, ttl time.Duration) (string, error)

	ParseStaffToken(tokenStr string) (int64, model.UserRole, int, *int64, time.Duration, error)

	GenerateGuestToken(bookingRoomID int64, ttl time.Duration) (string, error)

	ParseGuestToken(tokenStr string) (int64, error)
}
