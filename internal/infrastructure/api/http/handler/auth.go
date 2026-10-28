package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	authUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/auth"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg    *config.Config
	authUC authUC.AuthUseCase
}

func NewAuthHandler(
	cfg *config.Config,
	authUC authUC.AuthUseCase,
) *AuthHandler {
	return &AuthHandler{
		cfg,
		authUC,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	userRes, accessToken, refreshToken, err := h.authUC.Login(ctx, schema, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.storeTokenInCookie(c, accessToken, refreshToken, int(h.cfg.Token.AccessExpiresIn.Seconds()), int(h.cfg.Token.RefreshExpiresIn.Seconds()))

	response.APIResponse(c, http.StatusOK, constants.CodeLoginSuccess, "Login successfully", gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	accessToken := c.GetString(middleware.CtxAccessToken)
	refreshToken := c.GetString(middleware.CtxRefreshToken)

	if refreshToken == "" {
		_ = c.Error(exceptions.ErrInvalidUser)
		return
	}

	if err := h.authUC.Logout(ctx, schema, accessToken, refreshToken); err != nil {
		_ = c.Error(err)
		return
	}

	h.storeTokenInCookie(c, "", "", -1, -1)

	response.APIResponse(c, http.StatusOK, constants.CodeLogoutSuccess, "Logout successfully", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	refreshToken, err := c.Cookie(h.cfg.Token.RefreshName)
	if err != nil || refreshToken == "" {
		_ = c.Error(exceptions.ErrInvalidUser)
		return
	}

	newAccessToken, newRefreshToken, err := h.authUC.RefreshToken(ctx, schema, refreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.storeTokenInCookie(c, newAccessToken, newRefreshToken, int(h.cfg.Token.AccessExpiresIn.Seconds()), int(h.cfg.Token.RefreshExpiresIn.Seconds()))

	response.OKResponse(c, nil)
}

func (h *AuthHandler) UserInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	userRes, err := h.authUC.UserInfo(ctx, schema, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.authUC.ChangePassword(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	h.storeTokenInCookie(c, "", "", -1, -1)

	response.APIResponse(c, http.StatusOK, constants.CodeChangePasswordSuccess, "Change password successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	forgotPasswordToken, err := h.authUC.ForgotPassword(ctx, schema, req.Email)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeForgotPasswordSuccess, "Forgot password successfully", gin.H{
		"forgot_password_token": forgotPasswordToken,
	})
}

func (h *AuthHandler) VerifyForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.VerifyForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	resetPasswordToken, err := h.authUC.VerifyForgotPassword(ctx, schema, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeVerifyForgotPasswordSuccess, "Verify forgot password successfully", gin.H{
		"reset_password_token": resetPasswordToken,
	})
}

func (h *AuthHandler) ResendForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.ResendForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.authUC.ResendForgotPassword(ctx, schema, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.authUC.ResetPassword(ctx, schema, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeResetPasswordSuccess, "Reset password successfully", nil)
}

func (h *AuthHandler) UpdateInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.authUC.UpdateInfo(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateInfoSuccess, "Information updated successfully", nil)
}

func (h *AuthHandler) storeTokenInCookie(c *gin.Context, accessToken, refreshToken string, accessExpiresIn, refreshExpiresIn int) {
	isSecure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	domain := utils.ExtractRootDomain(c.Request.Host)

	c.SetSameSite(http.SameSiteLaxMode)

	if (accessToken != "" && accessExpiresIn > 0) || (refreshToken != "" && refreshExpiresIn > 0) {
		c.SetCookie(
			h.cfg.Token.GuestName,
			"",
			-1,
			"/",
			domain,
			isSecure,
			true,
		)
	}

	c.SetCookie(
		h.cfg.Token.AccessName,
		accessToken,
		accessExpiresIn,
		"/",
		domain,
		isSecure,
		true,
	)
	c.SetCookie(
		h.cfg.Token.RefreshName,
		refreshToken,
		refreshExpiresIn,
		"/",
		domain,
		isSecure,
		true,
	)
}
