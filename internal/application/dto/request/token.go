package request

type DeviceTokenRequest struct {
	DeviceToken string `json:"device_token" binding:"required"`
}
