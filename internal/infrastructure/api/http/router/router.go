package router

import (
	"net/http"

	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter(r *gin.Engine) *Router {
	return &Router{r}
}

func (r *Router) Setup(cfg config.ServerConfig, ctn *container.Container) {
	v2 := r.Engine.Group("")

	v2.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	v2.GET("/metrics", ctn.CtxHTTPMid.MetricsHandler())

	r.setupFileRoutes(v2, ctn.AuthHTTPMid, ctn.FileHTTPHdl)

	r.setupAuthRoutes(v2, ctn.AuthHTTPMid, ctn.AuthHTTPHdl)

	r.setupBranchRoutes(v2, ctn.AuthHTTPMid, ctn.BranchHTTPHdl)

	r.setupUserRoutes(v2, ctn.AuthHTTPMid, ctn.UserHTTPHdl)

	r.setupDepartmentRoutes(v2, ctn.AuthHTTPMid, ctn.DepartmentHTTPHdl)

	r.setupServiceTypeRoutes(v2, ctn.AuthHTTPMid, ctn.ServiceTypeHTTPHdl)

	r.setupServiceRoutes(v2, ctn.AuthHTTPMid, ctn.ServiceHTTPHdl)

	r.setupSupportCategoryRoutes(v2, ctn.AuthHTTPMid, ctn.SupportCategoryHTTPHdl)

	r.setupSourceRoutes(v2, ctn.AuthHTTPMid, ctn.SourceHTTPHdl)

	r.setupRoomTypeRoutes(v2, ctn.AuthHTTPMid, ctn.RoomTypeHTTPHdl)

	r.setupFloorRoutes(v2, ctn.AuthHTTPMid, ctn.FloorHTTPHdl)

	r.setupRoomRoutes(v2, ctn.AuthHTTPMid, ctn.RoomHTTPHdl)

	r.setupBookingRoutes(v2, ctn.AuthHTTPMid, ctn.BookingHTTPHdl)

	r.setupPermissionRoutes(v2, ctn.AuthHTTPMid, ctn.PermissionHTTPHdl)

	r.setupPricingTagRoutes(v2, ctn.AuthHTTPMid, ctn.PricingTagHTTPHdl)

	r.setupPricingCalendarRoutes(v2, ctn.AuthHTTPMid, ctn.PricingCalendarHTTPHdl)

	r.setupRoomTypePriceRoutes(v2, ctn.AuthHTTPMid, ctn.RoomTypePriceHTTPHdl)

	r.setupBookingRoomRoutes(v2, ctn.AuthHTTPMid, ctn.BookingRoomHTTPHdl)

	r.setupOrderServiceRoutes(v2, ctn.AuthHTTPMid, ctn.OrderServiceHTTPHdl)

	r.setupSupportTicketRoutes(v2, ctn.AuthHTTPMid, ctn.SupportTicketHTTPHdl)

	r.setupDeviceTokenRoutes(v2, ctn.AuthHTTPMid, ctn.DeviceTokenHTTPHdl)

	r.setupDashboardRoutes(v2, ctn.AuthHTTPMid, ctn.DashboardHTTPHdl)

	r.setupReviewRoutes(v2, ctn.AuthHTTPMid, ctn.ReviewHTTPHdl)

	r.setupGuestRoutes(v2, ctn.AuthHTTPMid, ctn.GuestHTTPHdl)

	r.setupCustomerRoutes(v2, ctn.AuthHTTPMid, ctn.CustomerHTTPHdl)

	r.setupEventRoutes(v2, ctn.AuthHTTPMid, ctn.EventHTTPHdl)

	r.setupTenantRoutes(v2, ctn.AuthHTTPMid, ctn.TenantHTTPHdl)

	r.setupNotificationRoutes(v2, ctn.AuthHTTPMid, ctn.NotificationHTTPHdl)

	r.setupChatRoutes(v2, ctn.AuthHTTPMid, ctn.ChatHTTPHdl)

	r.setupCQARoutes(v2, ctn.AuthHTTPMid, ctn.APIKeyHTTPMid, ctn.CQAHTTPHdl)

	r.setupWSRoutes(v2, ctn.AuthHTTPMid, ctn.WSHTTPHdl)
}
