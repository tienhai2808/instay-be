package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

const (
	CtxUserID        = "user_id"
	CtxAccessToken   = "access_token"
	CtxRefreshToken  = "refresh_token"
	CtxAccessTTL     = "access_ttl"
	CtxRole          = "role"
	CtxDepartmentID  = "department_id"
	CtxBookingRoomID = "booking_room_id"
	CtxClientID      = "client_id"
	CtxClientType    = "client_type"
)

type AuthMiddleware struct {
	cfg      config.TokenConfig
	logPro   port.LogProvider
	tokenPro port.TokenProvider
	cachePro port.CacheProvider
}

func NewAuthMiddleware(
	cfg config.TokenConfig,
	logPro port.LogProvider,
	tokenPro port.TokenProvider,
	cachePro port.CacheProvider,
) *AuthMiddleware {
	return &AuthMiddleware{
		cfg,
		logPro,
		tokenPro,
		cachePro,
	}
}

func (m *AuthMiddleware) IsAuthentication(c *gin.Context) {
	schema := c.GetString(CtxSchema)

	accessToken, err := c.Cookie(m.cfg.AccessName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	userID, role, tokenVersion, departmentID, ttl, err := m.tokenPro.ParseStaffToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	blackListKey := fmt.Sprintf("%s:black_list:%s", schema, accessToken)
	str, err := m.cachePro.GetString(ctx, blackListKey)
	if err != nil {
		m.logPro.Error("Get black list failed", port.String("schema", schema), port.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.APIResponse{
			Code:    exceptions.ErrInternalServerError.Code,
			Message: exceptions.ErrInternalServerError.Message,
		})
		return
	}

	if str != "" {
		c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
			Code:    exceptions.ErrInvalidUser.Code,
			Message: exceptions.ErrInvalidUser.Message,
		})
		return
	}

	userVersionKey := fmt.Sprintf("%s:user_version:%d", schema, userID)
	currentTokenVersion, err := m.cachePro.GetInt(ctx, userVersionKey)
	if err != nil {
		m.logPro.Error("Get token version failed", port.String("schema", schema), port.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.APIResponse{
			Code:    exceptions.ErrInternalServerError.Code,
			Message: exceptions.ErrInternalServerError.Message,
		})
		return
	}

	if tokenVersion != currentTokenVersion {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	c.Set(CtxUserID, userID)
	c.Set(CtxRole, string(role))
	c.Set(CtxDepartmentID, departmentID)
	c.Set(CtxAccessTTL, ttl)

	c.Next()
}

func (m *AuthMiddleware) AttachLogoutTokens(c *gin.Context) {
	accessToken, _ := c.Cookie(m.cfg.AccessName)
	refreshToken, err := c.Cookie(m.cfg.RefreshName)
	if err != nil || refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
			Code:    exceptions.ErrInvalidUser.Code,
			Message: exceptions.ErrInvalidUser.Message,
		})
		return
	}

	c.Set(CtxAccessToken, accessToken)
	c.Set(CtxRefreshToken, refreshToken)

	c.Next()
}

func (m *AuthMiddleware) HasRole(allowedRole model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetString(CtxRole)
		role := model.UserRole(roleStr)

		if roleStr == "" || !role.IsValid() || role != allowedRole {
			c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
				Code:    exceptions.ErrForbidden.Code,
				Message: exceptions.ErrForbidden.Message,
			})
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) HasPermission(allowedPermission model.PermissionName) gin.HandlerFunc {
	return func(c *gin.Context) {
		schema := c.GetString(CtxSchema)

		departmentIDAny, exists := c.Get(CtxDepartmentID)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
				Code:    exceptions.ErrForbidden.Code,
				Message: exceptions.ErrForbidden.Message,
			})
			return
		}

		departmentID, ok := departmentIDAny.(*int64)
		if !ok || departmentID == nil {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		cacheKey := fmt.Sprintf("%s:dept_perms:%d", schema, *departmentID)
		bytes, err := m.cachePro.GetObject(ctx, cacheKey)
		if err != nil {
			m.logPro.Error("Get permissions by department id in cache failed", port.String("schema", schema), port.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.APIResponse{
				Code:    exceptions.ErrInternalServerError.Code,
				Message: exceptions.ErrInternalServerError.Message,
			})
			return
		}
		if bytes == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
				Code:    exceptions.ErrForbidden.Code,
				Message: exceptions.ErrForbidden.Message,
			})
			return
		}

		var perms []model.PermissionName
		if err = json.Unmarshal(bytes, &perms); err != nil {
			m.logPro.Error("Json unmarshal permissions data failed", port.String("schema", schema), port.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.APIResponse{
				Code:    exceptions.ErrInternalServerError.Code,
				Message: exceptions.ErrInternalServerError.Message,
			})
			return
		}

		if !slices.Contains(perms, allowedPermission) {
			c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
				Code:    exceptions.ErrForbidden.Code,
				Message: exceptions.ErrForbidden.Message,
			})
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) IsGuest(c *gin.Context) {
	guestToken, err := c.Cookie(m.cfg.GuestName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	bookingRoomID, err := m.tokenPro.ParseGuestToken(guestToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	c.Set(CtxBookingRoomID, bookingRoomID)

	c.Next()
}

func (m *AuthMiddleware) OptionalAuth(c *gin.Context) {
	schema := c.GetString(CtxSchema)

	if accessToken, err := c.Cookie(m.cfg.AccessName); err == nil {
		if userID, _, tokenVersion, departmentID, _, err := m.tokenPro.ParseStaffToken(accessToken); err == nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()

			blackListKey := fmt.Sprintf("%s:black_list:%s", schema, accessToken)
			str, _ := m.cachePro.GetString(ctx, blackListKey)

			userVersionKey := fmt.Sprintf("%s:user_version:%d", schema, userID)
			currentTokenVersion, _ := m.cachePro.GetInt(ctx, userVersionKey)

			if str == "" && tokenVersion == currentTokenVersion {
				c.Set(CtxClientID, userID)
				c.Set(CtxClientType, constants.ClientTypeStaff)
				c.Set(CtxDepartmentID, departmentID)
				c.Next()
				return
			}
		}
	}

	if guestToken, err := c.Cookie(m.cfg.GuestName); err == nil {
		if bookingRoomID, err := m.tokenPro.ParseGuestToken(guestToken); err == nil {
			c.Set(CtxClientID, bookingRoomID)
			c.Set(CtxClientType, constants.ClientTypeGuest)
			c.Next()
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
		Code:    exceptions.ErrUnauth.Code,
		Message: exceptions.ErrUnauth.Message,
	})
}
