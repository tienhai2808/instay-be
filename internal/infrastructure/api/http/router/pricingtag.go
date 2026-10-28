package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupPricingTagRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.PricingTagHandler) {
	pricingTag := rg.Group("/pricing-tags", authMid.IsAuthentication)
	{
		pricingTag.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		pricingTag.GET("", authMid.HasRole(model.RoleAdmin), hdl.List)

		pricingTag.PUT("/:id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		pricingTag.DELETE("/:id", authMid.HasRole(model.RoleAdmin), hdl.Delete)
	}
}
