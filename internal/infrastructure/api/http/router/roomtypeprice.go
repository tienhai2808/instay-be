package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoomTypePriceRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.RoomTypePriceHandler) {
	roomTypePrice := rg.Group("/room-type-prices", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		roomTypePrice.POST("", hdl.Create)

		roomTypePrice.PUT("/:id", hdl.Update)

		roomTypePrice.DELETE("/:id", hdl.Delete)
	}
}
