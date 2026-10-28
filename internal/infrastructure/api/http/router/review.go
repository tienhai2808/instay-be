package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupReviewRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.ReviewHandler) {
	guest := rg.Group("/reviews", authMid.IsGuest)
	{
		guest.POST("", hdl.Create)

		guest.GET("", hdl.ListForGuest)

		guest.GET("/me", hdl.GetByBookingRoomID)

		guest.PUT("/me", hdl.Update)
	}

	admin := rg.Group("/admin/reviews", authMid.IsAuthentication)
	{
		admin.GET("", hdl.ListForAdmin)

		admin.GET("/summary", hdl.Summary)
	}
}
