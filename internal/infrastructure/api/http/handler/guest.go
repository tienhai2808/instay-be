package handler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	guestUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/guest"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	httpRes "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type GuestHandler struct {
	cfg     *config.Config
	guestUC guestUC.GuestUseCase
}

func NewGuestHandler(cfg *config.Config, guestUC guestUC.GuestUseCase) *GuestHandler {
	return &GuestHandler{cfg, guestUC}
}

func (h *GuestHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.guestUC.Create(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusCreated, constants.CodeCreateGuestSuccess, "Guest created successfully", nil)
}

func (h *GuestHandler) CreateVerification(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.CreateGuestVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	callbackURL := h.resolveCallbackURL(c)

	sessionID, verificationURL, err := h.guestUC.CreateVerification(ctx, schema, req.BookingRoomID, callbackURL)
	if err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.OKResponse(c, gin.H{
		"session_id":       sessionID,
		"verification_url": verificationURL,
	})
}

func (h *GuestHandler) Verify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.VerifyGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.verifyKYCSignature(c, req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.guestUC.Verify(ctx, req); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.OKResponse(c, nil)
}

func (h *GuestHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	guestIDStr := c.Param("id")
	guestID, err := strconv.ParseInt(guestIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err = h.guestUC.Update(ctx, schema, guestID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusOK, constants.CodeUpdateGuestSuccess, "Guest updated successfully", nil)
}

func (h *GuestHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	guestIDStr := c.Param("id")
	guestID, err := strconv.ParseInt(guestIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err = h.guestUC.Delete(ctx, schema, guestID, userID); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusOK, constants.CodeSuccess, "Guest deleted successfully", nil)
}

func (h *GuestHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.GuestPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	guestsRes, meta, err := h.guestUC.List(ctx, schema, userID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.OKResponse(c, gin.H{
		"guests": guestsRes,
		"meta":   meta,
	})
}

func (h *GuestHandler) Export(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.ExportListGuestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	guestsRes, err := h.guestUC.Export(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	csvData, err := buildGuestExportCSV(guestsRes)
	if err != nil {
		_ = c.Error(err)
		return
	}

	fileName := fmt.Sprintf("guests_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", csvData)
}

func buildGuestExportCSV(guests []*response.GuestResponse) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return nil, err
	}

	w := csv.NewWriter(&buf)

	headers := []string{"STT", "Họ và tên", "Ngày sinh", "QT", "Số HC", "Ngày đến", "Ngày đi DK", "Phòng"}
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for i, guest := range guests {
		roomName := ""
		if guest.BookingRoom != nil && guest.BookingRoom.Room != nil {
			roomName = guest.BookingRoom.Room.Name
		}

		record := []string{
			strconv.Itoa(i + 1),
			guest.FullName,
			formatCSVDate(guest.DOB),
			guest.National,
			guest.PassportNumber,
			formatCSVDate(guest.CheckIn),
			formatCSVDate(guest.CheckOut),
			roomName,
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func formatCSVDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func (h *GuestHandler) resolveCallbackURL(c *gin.Context) string {
	referer := strings.TrimSpace(c.GetHeader("Referer"))
	if referer != "" {
		return referer
	}

	origin := strings.TrimSpace(c.GetHeader("Origin"))
	if origin != "" {
		return origin
	}

	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host
}

func (h *GuestHandler) verifyKYCSignature(c *gin.Context, req request.VerifyGuestRequest) error {
	secret := strings.TrimSpace(h.cfg.KYC.WebhookSecret)
	if secret == "" {
		return nil
	}

	signature := strings.TrimSpace(c.GetHeader("X-Signature-Simple"))
	timestampHeader := strings.TrimSpace(c.GetHeader("X-Timestamp"))
	if signature == "" || timestampHeader == "" {
		return fmt.Errorf("missing signature header")
	}

	timestamp, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if timestamp < now-300 || timestamp > now+300 {
		return fmt.Errorf("expired signature")
	}

	payloadTimestamp := req.Timestamp
	if payloadTimestamp == 0 {
		payloadTimestamp = timestamp
	}

	canonicalString := fmt.Sprintf("%d:%s:%s:%s", payloadTimestamp, req.SessionID, req.Status, req.WebhookType)

	hmacHash := hmac.New(sha256.New, []byte(secret))
	_, _ = hmacHash.Write([]byte(canonicalString))
	expectedSignature := hex.EncodeToString(hmacHash.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
