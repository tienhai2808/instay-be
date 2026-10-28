package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupPricingCalendarRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.PricingCalenderHandler) {
	pricingCalendar := rg.Group("/pricing-calendars", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		pricingCalendar.GET("", hdl.List)

		pricingCalendar.PUT("", hdl.UpdateBulk)
	}
}
