package external

import "time"

type ExternalConversationsQuery struct {
	Since  *time.Time
	Limit  uint32
	Cursor string
}

type ExternalMessagesQuery struct {
	ConversationID int64
	Since          *time.Time
	Limit          uint32
	Cursor         string
}

type ExternalConversationResponse struct {
	ExternalID     int64          `json:"external_id,string"`
	ExternalUserID int64          `json:"external_user_id,string"`
	CustomerName   string         `json:"customer_name"`
	LastMessageAt  time.Time      `json:"last_message_at"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

type ExternalAttachmentResponse struct {
	Type string `json:"type"`
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
}

type ExternalMessageResponse struct {
	ExternalID  string                        `json:"external_id"`
	SenderType  string                        `json:"sender_type"`
	SenderName  string                        `json:"sender_name,omitempty"`
	Content     string                        `json:"content"`
	ContentType string                        `json:"content_type"`
	Attachments []*ExternalAttachmentResponse `json:"attachments"`
	SentAt      time.Time                     `json:"sent_at"`
	RawData     map[string]any                `json:"raw_data,omitempty"`
}
