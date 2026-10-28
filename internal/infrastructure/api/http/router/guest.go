package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupGuestRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.GuestHandler) {
	guest := rg.Group("/guests", authMid.IsAuthentication, authMid.HasPermission(model.PermGuestManagement))
	{
		guest.POST("", hdl.Create)

		guest.POST("/create-verification", hdl.CreateVerification)

		guest.PUT("/:id", hdl.Update)

		guest.DELETE("/:id", hdl.Delete)

		guest.GET("", hdl.List)

		guest.POST("/export", hdl.Export)
	}

	webhook := rg.Group("/guests")
	{
		webhook.POST("/verify", hdl.Verify)
	}
}
