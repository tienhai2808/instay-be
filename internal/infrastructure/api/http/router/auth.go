package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupAuthRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", hdl.Login)

		auth.POST("/logout", authMid.AttachLogoutTokens, hdl.Logout)

		auth.POST("/refresh-token", hdl.RefreshToken)

		auth.GET("/userinfo", authMid.IsAuthentication, hdl.UserInfo)

		auth.PUT("/change-password", authMid.IsAuthentication, hdl.ChangePassword)

		auth.POST("/forgot-password", hdl.ForgotPassword)

		auth.POST("/forgot-password/verify", hdl.VerifyForgotPassword)

		auth.POST("/forgot-password/resend", hdl.ResendForgotPassword)

		auth.PUT("/reset-password", hdl.ResetPassword)

		auth.PUT("/update-info", authMid.IsAuthentication, hdl.UpdateInfo)
	}
}
