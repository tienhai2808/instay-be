package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupDeviceTokenRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.DeviceTokenHandler) {
	deviceToken := rg.Group("/device-tokens")
	{
		deviceToken.POST("", authMid.OptionalAuth, hdl.Register)

		deviceToken.DELETE("", authMid.OptionalAuth, hdl.Delete)
	}
}
