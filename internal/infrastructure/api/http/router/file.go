package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupFileRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.FileHandler) {
	file := rg.Group("/files")
	{
		file.POST("/presigned-urls/uploads", authMid.OptionalAuth, hdl.UploadPresignedURLs)

		file.POST("/presigned-urls/views", authMid.OptionalAuth, hdl.ViewPresignedURLs)
	}
}
