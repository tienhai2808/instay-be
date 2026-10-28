package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	chatUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/chat"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatUC chatUC.ChatUseCase
}

func NewChatHandler(chatUC chatUC.ChatUseCase) *ChatHandler {
	return &ChatHandler{chatUC}
}

func (h *ChatHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.ChatPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	chatsRes, meta, err := h.chatUC.ListForGuest(ctx, schema, bookingRoomID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"chats": chatsRes,
		"meta":  meta,
	})
}

func (h *ChatHandler) ListMessagesByIDForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	include := c.Query("include")
	includeChat := include == "chat"

	var query request.MessagePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	chatRes, msgsRes, meta, err := h.chatUC.ListMessagesByIDForGuest(ctx, schema, chatID, query, includeChat)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := gin.H{
		"messages": msgsRes,
		"meta":     meta,
	}

	if includeChat {
		res["chat"] = chatRes
	}

	response.OKResponse(c, res)
}

func (h *ChatHandler) AvailabilityForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.ChatAvailabilityForGuestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	availability, err := h.chatUC.AvailabilityForGuest(ctx, schema, bookingRoomID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, availability)
}

func (h *ChatHandler) AvailabilityForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	deptIDAny, _ := c.Get(middleware.CtxDepartmentID)
	deptID, _ := deptIDAny.(*int64)
	if deptID == nil {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	var query request.ChatAvailabilityForAdminQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	availability, err := h.chatUC.AvailabilityForAdmin(ctx, schema, userID, *deptID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, availability)
}

func (h *ChatHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	departmentIDAny, exists := c.Get(middleware.CtxDepartmentID)
	if !exists {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	deptID, ok := departmentIDAny.(*int64)
	if !ok {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	var query request.ChatPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	chatsRes, meta, err := h.chatUC.ListForAdmin(ctx, schema, userID, deptID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"chats": chatsRes,
		"meta":  meta,
	})
}

func (h *ChatHandler) ListMessagesByIDForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	include := c.Query("include")
	includeChat := include == "chat"

	var query request.MessagePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	chatRes, msgsRes, meta, err := h.chatUC.ListMessagesByIDForAdmin(ctx, schema, userID, chatID, query, includeChat)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := gin.H{
		"messages": msgsRes,
		"meta":     meta,
	}

	if includeChat {
		res["chat"] = chatRes
	}

	response.OKResponse(c, res)
}
