package container

import (
	httpHdl "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	httpMid "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/realtime/ws"
)

func (c *Container) initAPI() {
	c.CtxHTTPMid = httpMid.NewContextMiddleware(c.LogPro)
	c.DomainHTTPMid = httpMid.NewDomainMiddleware(c.LogPro, c.CachePro, c.TenantRepo)
	c.AuthHTTPMid = httpMid.NewAuthMiddleware(c.cfg.Token, c.LogPro, c.tokenPro, c.CachePro)
	c.APIKeyHTTPMid = httpMid.NewAPIKeyMiddleware(c.apiKeyPro)

	c.FileHTTPHdl = httpHdl.NewFileHandler(c.fileUC)
	c.AuthHTTPHdl = httpHdl.NewAuthHandler(c.cfg, c.authUC)
	c.BranchHTTPHdl = httpHdl.NewBranchHandler(c.branchUC)
	c.UserHTTPHdl = httpHdl.NewUserHandler(c.userUC)
	c.DepartmentHTTPHdl = httpHdl.NewDepartmentHandler(c.departmentUC)
	c.ServiceTypeHTTPHdl = httpHdl.NewServiceTypeHandler(c.serviceTypeUC)
	c.ServiceHTTPHdl = httpHdl.NewServiceHandler(c.serviceUC)
	c.SupportCategoryHTTPHdl = httpHdl.NewSupportCategoryHandler(c.supportCategoryUC)
	c.SourceHTTPHdl = httpHdl.NewSourceHandler(c.sourceUC)
	c.RoomTypeHTTPHdl = httpHdl.NewRoomTypeHandler(c.roomTypeUC)
	c.FloorHTTPHdl = httpHdl.NewFloorHandler(c.floorUC)
	c.RoomHTTPHdl = httpHdl.NewRoomHandler(c.roomUC)
	c.BookingHTTPHdl = httpHdl.NewBookingHandler(c.bookingUC)
	c.PermissionHTTPHdl = httpHdl.NewPermissionHandler(c.permissionUC)
	c.PricingTagHTTPHdl = httpHdl.NewPricingTagHandler(c.pricingTagUC)
	c.PricingCalendarHTTPHdl = httpHdl.NewPricingCalendarHandler(c.pricingCalendarUC)
	c.RoomTypePriceHTTPHdl = httpHdl.NewRoomTypePriceHandler(c.roomTypePriceUC)
	c.BookingRoomHTTPHdl = httpHdl.NewBookingRoomHandler(c.cfg, c.bookingRoomUC)
	c.OrderServiceHTTPHdl = httpHdl.NewOrderServiceHandler(c.orderServiceUC)
	c.SupportTicketHTTPHdl = httpHdl.NewSupportTicketHandler(c.supportTicketUC)
	c.DeviceTokenHTTPHdl = httpHdl.NewDeviceTokenHandler(c.deviceTokenUC)
	c.DashboardHTTPHdl = httpHdl.NewDashboardHandler(c.dashboardUC)
	c.ReviewHTTPHdl = httpHdl.NewReviewHandler(c.reviewUC)
	c.GuestHTTPHdl = httpHdl.NewGuestHandler(c.cfg, c.guestUC)
	c.CustomerHTTPHdl = httpHdl.NewCustomerHandler(c.customerUC)
	c.EventHTTPHdl = httpHdl.NewEventHandler(c.eventUC)
	c.NotificationHTTPHdl = httpHdl.NewNotificationHandler(c.notificationUC)
	c.TenantHTTPHdl = httpHdl.NewTenantHandler(c.tenantUC)
	c.ChatHTTPHdl = httpHdl.NewChatHandler(c.chatUC)
	c.CQAHTTPHdl = httpHdl.NewCQAHandler(c.chatUC, c.apiKeyPro)
	c.HubWS = ws.NewHubWS(c.messageUC, c.CachePro, c.LogPro)
	c.WSHTTPHdl = httpHdl.NewWSHandler(c.HubWS, c.DomainHTTPMid)
}
