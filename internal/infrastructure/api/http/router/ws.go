package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupWSRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.WSHandler) {
	ws := rg.Group("/ws")
	{
		ws.GET("", authMid.OptionalAuth, hdl.ServeWS)
	}
}
