package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/external"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	chatUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/chat"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

const (
	defaultExternalLimit = 100
	maxExternalLimit     = 100
)

type CQAHandler struct {
	chatUC    chatUC.ChatUseCase
	apiKeyPro port.APIKeyProvider
}

func NewCQAHandler(chatUC chatUC.ChatUseCase, apiKeyPro port.APIKeyProvider) *CQAHandler {
	return &CQAHandler{
		chatUC:    chatUC,
		apiKeyPro: apiKeyPro,
	}
}

func (h *CQAHandler) CreateAPIKey(c *gin.Context) {
	schema := c.GetString(middleware.CtxSchema)
	if schema == "" {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	key, err := h.apiKeyPro.Generate(schema)
	if err != nil {
		_ = c.Error(exceptions.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"api_key":        key,
		"api_key_header": h.apiKeyPro.HeaderName(),
		"api_key_prefix": h.apiKeyPro.HeaderPrefix(),
	})
}

func (h *CQAHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *CQAHandler) ListConversations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	limit, parseErr := parseExternalLimit(c.Query("limit"))
	if parseErr != nil {
		_ = c.Error(parseErr.WithData(gin.H{"field": "limit"}))
		return
	}

	since, parseErr := parseRFC3339Time(c.Query("since"))
	if parseErr != nil {
		_ = c.Error(parseErr.WithData(gin.H{"field": "since"}))
		return
	}

	conversations, nextCursor, err := h.chatUC.ListForExternal(ctx, schema, external.ExternalConversationsQuery{
		Since:  since,
		Limit:  limit,
		Cursor: strings.TrimSpace(c.Query("cursor")),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
		"next_cursor":   nextCursor,
	})
}

func (h *CQAHandler) ListMessages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	conversationID, parseErr := parseConversationID(c.Query("conversation_id"))
	if parseErr != nil {
		_ = c.Error(parseErr.WithData(gin.H{"field": "conversation_id"}))
		return
	}

	limit, parseErr := parseExternalLimit(c.Query("limit"))
	if parseErr != nil {
		_ = c.Error(parseErr.WithData(gin.H{"field": "limit"}))
		return
	}

	since, parseErr := parseRFC3339Time(c.Query("since"))
	if parseErr != nil {
		_ = c.Error(parseErr.WithData(gin.H{"field": "since"}))
		return
	}

	messages, nextCursor, err := h.chatUC.ListMessagesForExternal(ctx, schema, external.ExternalMessagesQuery{
		ConversationID: conversationID,
		Since:          since,
		Limit:          limit,
		Cursor:         strings.TrimSpace(c.Query("cursor")),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":    messages,
		"next_cursor": nextCursor,
	})
}

func parseExternalLimit(raw string) (uint32, *exceptions.Exception) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultExternalLimit, nil
	}

	limit, err := strconv.ParseUint(raw, 10, 32)
	if err != nil || limit == 0 {
		return 0, exceptions.ErrBadRequest
	}

	if limit > maxExternalLimit {
		return maxExternalLimit, nil
	}

	return uint32(limit), nil
}

func parseRFC3339Time(raw string) (*time.Time, *exceptions.Exception) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, exceptions.ErrBadRequest
	}

	return &parsed, nil
}

func parseConversationID(raw string) (int64, *exceptions.Exception) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, exceptions.ErrBadRequest
	}

	raw = strings.TrimPrefix(raw, "conv_")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exceptions.ErrBadRequest
	}

	return id, nil
}
