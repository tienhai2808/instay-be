package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupBookingRoomRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.BookingRoomHandler) {
	bookingRoom := rg.Group("/booking-rooms")
	{
		bookingRoom.POST("", authMid.IsAuthentication, authMid.HasPermission(model.PermBookingManagement), hdl.Create)

		bookingRoom.POST("/verify", hdl.Verify)

		bookingRoom.GET("/info", authMid.IsGuest, hdl.Info)

		bookingRoom.GET("/:id", authMid.OptionalAuth, hdl.GetByID)

		bookingRoom.GET("/:id/guests", authMid.IsAuthentication, authMid.HasPermission(model.PermBookingManagement), hdl.ListGuestsByID)

		bookingRoom.PUT("/:id", authMid.IsAuthentication, authMid.HasPermission(model.PermBookingManagement), hdl.Update)
	}
}
