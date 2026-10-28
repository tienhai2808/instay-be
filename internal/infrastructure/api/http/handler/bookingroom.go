package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	bookingRoomUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/bookingroom"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type BookingRoomHandler struct {
	cfg           *config.Config
	bookingRoomUC bookingRoomUC.BookingRoomUseCase
}

func NewBookingRoomHandler(
	cfg *config.Config,
	bookingRoomUC bookingRoomUC.BookingRoomUseCase,
) *BookingRoomHandler {
	return &BookingRoomHandler{
		cfg,
		bookingRoomUC,
	}
}

func (h *BookingRoomHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateBookingRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.bookingRoomUC.Create(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateBookingRoomSuccess, "Booking room created successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *BookingRoomHandler) Verify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.VerifyBookingRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	bookingRoomRes, guestToken, ttl, err := h.bookingRoomUC.Verify(ctx, schema, req.Code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.storeTokenInCookie(c, guestToken, int(ttl.Seconds()))

	response.APIResponse(c, http.StatusOK, constants.CodeVerifyBookingRoomSuccess, "Verify booking room successfully", gin.H{
		"booking_room": bookingRoomRes,
	})
}

func (h *BookingRoomHandler) Info(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	bookingRoomRes, err := h.bookingRoomUC.Info(ctx, schema, bookingRoomID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"booking_room": bookingRoomRes,
	})
}

func (h *BookingRoomHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	bookingRoomIDStr := c.Param("id")
	bookingRoomID, err := strconv.ParseInt(bookingRoomIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	view := c.Query("view")

	if view == "basic" {
		bookingRoomRes, err := h.bookingRoomUC.GetByIDBasic(ctx, schema, bookingRoomID, userID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response.OKResponse(c, gin.H{
			"booking_room": bookingRoomRes,
		})
		return
	}

	bookingRoomRes, err := h.bookingRoomUC.GetByIDFull(ctx, schema, bookingRoomID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"booking_room": bookingRoomRes,
	})
}

func (h *BookingRoomHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	bookingRoomIDStr := c.Param("id")
	bookingRoomID, err := strconv.ParseInt(bookingRoomIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateBookingRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.bookingRoomUC.Update(ctx, schema, bookingRoomID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateBookingSuccess, "Booking updated successfully", nil)
}

func (h *BookingRoomHandler) ListGuestsByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	bookingRoomIDStr := c.Param("id")
	bookingRoomID, err := strconv.ParseInt(bookingRoomIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	guestsRes, err := h.bookingRoomUC.ListGuestsByID(ctx, schema, bookingRoomID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"guests": guestsRes,
	})
}

func (h *BookingRoomHandler) storeTokenInCookie(c *gin.Context, guestToken string, expiresIn int) {
	isSecure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	domain := utils.ExtractRootDomain(c.Request.Host)

	c.SetSameSite(http.SameSiteLaxMode)

	if guestToken != "" && expiresIn > 0 {
		c.SetCookie(h.cfg.Token.AccessName, "", -1, "/", domain, isSecure, true)
		c.SetCookie(h.cfg.Token.RefreshName, "", -1, "/", domain, isSecure, true)
	}

	c.SetCookie(
		h.cfg.Token.GuestName,
		guestToken,
		expiresIn,
		"/",
		domain,
		isSecure,
		true,
	)
}
