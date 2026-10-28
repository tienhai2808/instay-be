package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupUserRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.UserHandler) {
	user := rg.Group("/users", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		user.POST("", hdl.Create)

		user.GET("/:id", hdl.GetByID)

		user.GET("", hdl.List)

		user.PUT("/:id", hdl.Update)

		user.PUT("/:id/password", hdl.UpdatePassword)

		user.DELETE("/:id", hdl.Delete)

		user.DELETE("", hdl.DeleteBatch)
	}
}
