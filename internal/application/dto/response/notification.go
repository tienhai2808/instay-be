package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type NotificationResponse struct {
	ID          int64                         `json:"id,string"`
	ContentType model.NotificationContentType `json:"content_type"`
	ContentKey  string                        `json:"content_key"`
	ContentVI   string                        `json:"content_vi"`
	ContentEN   string                        `json:"content_en"`
	IsRead      bool                          `json:"is_read"`
	CreatedAt   time.Time                     `json:"created_at"`
}
