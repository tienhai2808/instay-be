package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupChatRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.ChatHandler) {
	guest := rg.Group("/chats", authMid.IsGuest)
	{
		guest.GET("", hdl.ListForGuest)

		guest.GET("/:id/messages", hdl.ListMessagesByIDForGuest)

		guest.GET("/availability", hdl.AvailabilityForGuest)
	}

	admin := rg.Group("/admin/chats", authMid.IsAuthentication)
	{
		admin.GET("", hdl.ListForAdmin)

		admin.GET("/:id/messages", hdl.ListMessagesByIDForAdmin)

		admin.GET("/availability", hdl.AvailabilityForAdmin)
	}
}
