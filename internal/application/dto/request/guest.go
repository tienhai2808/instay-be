package request

import "time"

type GuestPaginationQuery struct {
	Page     uint32    `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32    `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort     string    `form:"sort" json:"sort"`
	Order    string    `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Search   string    `form:"search" json:"search"`
	FromDate time.Time `form:"from_date" time_format:"2006-01-02"`
	ToDate   time.Time `form:"to_date" time_format:"2006-01-02"`
}

type CreateGuestRequest struct {
	FullName       string    `json:"full_name" binding:"required,min=2"`
	DOB            time.Time `json:"dob" binding:"required"`
	National       string    `json:"national" binding:"required,min=2"`
	PassportNumber string    `json:"passport_number" binding:"required,max=20"`
	BookingRoomID  int64     `json:"booking_room_id,string" binding:"required"`
}

type CreateGuestVerificationRequest struct {
	BookingRoomID int64 `json:"booking_room_id,string" binding:"required"`
}

type VerifyGuestRequest struct {
	SessionID   string                      `json:"session_id" binding:"required"`
	Status      string                      `json:"status" binding:"required"`       // Trạng thái tổng: "Approved", "Declined"
	WebhookType string                      `json:"webhook_type" binding:"required"` // Luôn là "status.updated"
	Timestamp   int64                       `json:"timestamp"`
	VendorData  string                      `json:"vendor_data"` // Dữ liệu mapping của hệ thống bạn
	Decision    *VerifyGuestDecisionRequest `json:"decision"`
}

type VerifyGuestDecisionRequest struct {
	Status          string                             `json:"status"`
	VendorData      string                             `json:"vendor_data"`
	Features        []string                           `json:"features"` // Sẽ trả về ["ID_VERIFICATION"]
	IDVerifications []VerifyGuestIDVerificationRequest `json:"id_verifications"`
}

type VerifyGuestIDVerificationRequest struct {
	Status         string `json:"status"`                                   // "Approved" nếu đọc thành công
	DocumentType   string `json:"document_type" binding:"required,doctype"` // "Identity Card" hoặc "Passport"
	DocumentNumber string `json:"document_number"`                          // Số CCCD / Số Hộ chiếu
	PersonalNumber string `json:"personal_number"`

	// Thông tin cá nhân
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	DateOfBirth string `json:"date_of_birth"` // Format: YYYY-MM-DD
	Gender      string `json:"gender"`        // "M" hoặc "F"
	Nationality string `json:"nationality"`
	Address     string `json:"address"`

	// Ngày cấp / Hết hạn
	DateOfIssue    string `json:"date_of_issue"`
	ExpirationDate string `json:"expiration_date"`

	// Link ảnh đã được Didit cắt gọn gàng
	PortraitImage string `json:"portrait_image"` // Ảnh chân dung trên thẻ
	FrontImage    string `json:"front_image"`    // Ảnh mặt trước thẻ
	BackImage     string `json:"back_image"`     // Ảnh mặt sau thẻ
}

type UpdateGuestRequest struct {
	FullName       string    `json:"full_name" binding:"required,min=2"`
	DOB            time.Time `json:"dob" binding:"required"`
	National       string    `json:"national" binding:"required,min=2"`
	PassportNumber string    `json:"passport_number" binding:"required,max=20"`
	CheckIn        time.Time `json:"check_in" binding:"required"`
	CheckOut       time.Time `json:"check_out" binding:"required"`
	BookingRoomID  int64     `json:"booking_room_id,string" binding:"required"`
}

type ExportListGuestsRequest struct {
	Sort     *string    `binding:"omitempty,oneof=check_in check_out" json:"sort"`
	Order    *string    `binding:"omitempty,oneof=asc desc" json:"order"`
	FromDate *time.Time `time_format:"2006-01-02" json:"from_date" binding:"omitempty"`
	ToDate   *time.Time `time_format:"2006-01-02" json:"to_date" binding:"omitempty"`
}
