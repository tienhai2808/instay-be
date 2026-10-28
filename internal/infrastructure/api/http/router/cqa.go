package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupCQARoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, apiKeyMid *middleware.APIKeyMiddleware, hdl *handler.CQAHandler) {
	admin := rg.Group("/api-keys", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		admin.POST("", hdl.CreateAPIKey)
	}

	public := rg.Group("", apiKeyMid.IsValidAPIKey)
	{
		public.GET("/health", hdl.Health)

		public.GET("/conversations", hdl.ListConversations)

		public.GET("/messages", hdl.ListMessages)
	}
}
