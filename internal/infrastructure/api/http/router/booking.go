package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupBookingRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.BookingHandler) {
	admin := rg.Group("/bookings", authMid.IsAuthentication, authMid.HasPermission(model.PermBookingManagement))
	{
		admin.POST("", hdl.Create)

		admin.PUT("/:id", hdl.Update)

		admin.PUT("/:id/status", hdl.UpdateStatus)

		admin.DELETE("/:id", hdl.Delete)

		admin.GET("", hdl.List)

		admin.GET("/statuses", hdl.StatusSummary)

		admin.GET("/timeline", hdl.Timeline)

		admin.GET("/:id", hdl.GetByID)
	}
}
