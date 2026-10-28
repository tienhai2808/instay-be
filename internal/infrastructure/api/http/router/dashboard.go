package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupDashboardRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.DashboardHandler) {
	admin := rg.Group("/dashboards", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		admin.GET("/booking-overview", hdl.BookingOverview)

		admin.GET("/booking-revenue", hdl.BookingRevenue)

		admin.GET("/room-type-revenue", hdl.RoomTypeRevenue)

		admin.GET("/source-revenue", hdl.SourceRevenue)

		admin.GET("/hotel-operation", hdl.HotelOperation)

		admin.GET("/guest-national", hdl.GuestNational)

		admin.GET("/support-ticket-time", hdl.SupportTicketTime)

		admin.GET("/service-type-count", hdl.ServiceTypeCount)

		admin.GET("/review-summary", hdl.ReviewSummary)

		admin.GET("/review-category-count", hdl.ReviewCategoryCount)

		admin.GET("/review-category-sentiment", hdl.ReviewCategorySentiment)

		admin.GET("/review-sentiment-count", hdl.ReviewSentimentCount)

		admin.GET("/review-sentiment-time", hdl.ReviewSentimentTime)
	}
}
