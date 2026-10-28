package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoomTypeRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.RoomTypeHandler) {
	roomTypeAuth := rg.Group("/room-types", authMid.IsAuthentication, authMid.HasPermission(model.PermBookingManagement))
	{
		roomTypeAuth.GET("/:id/rooms", hdl.ListRoomsByID)

		roomTypeAuth.POST("/price-calculation", hdl.PriceCalculation)
	}

	roomType := rg.Group("/room-types", authMid.IsAuthentication)
	{
		roomType.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		roomType.GET("", authMid.HasPermission(model.PermBookingManagement), hdl.List)

		roomType.PUT("/:id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		roomType.DELETE("/:id", authMid.HasRole(model.RoleAdmin), hdl.Delete)
	}
}
