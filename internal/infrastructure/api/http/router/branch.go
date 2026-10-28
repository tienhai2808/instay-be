package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupBranchRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.BranchHandler) {
	branch := rg.Group("/branches", authMid.IsAuthentication)
	{
		branch.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		branch.GET("", hdl.List)

		branch.PUT("/:id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		branch.DELETE("/:id", authMid.HasRole(model.RoleAdmin), hdl.Delete)
	}
}
