package handler

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	customerUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/customer"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	httpRes "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerUC customerUC.CustomerUseCase
}

func NewCustomerHandler(customerUC customerUC.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{customerUC}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.customerUC.Create(ctx, schema, bookingRoomID, req); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusCreated, constants.CodeCreateCustomerSuccess, "Customer created successfully", nil)
}

func (h *CustomerHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.CustomerPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	customersRes, meta, err := h.customerUC.List(ctx, schema, userID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.OKResponse(c, gin.H{
		"customers": customersRes,
		"meta":      meta,
	})
}

func (h *CustomerHandler) Export(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.ExportListCustomersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	customersRes, err := h.customerUC.Export(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	csvData, err := buildCustomerExportCSV(customersRes)
	if err != nil {
		_ = c.Error(err)
		return
	}

	fileName := fmt.Sprintf("customers_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", csvData)
}

func buildCustomerExportCSV(customers []*response.CustomerResponse) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return nil, err
	}

	w := csv.NewWriter(&buf)

	headers := []string{"STT", "birth_year", "gender", "phone", "created_at"}
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for i, customer := range customers {
		record := []string{
			strconv.Itoa(i + 1),
			strconv.FormatUint(uint64(customer.BirthYear), 10),
			string(customer.Gender),
			customer.Phone,
			formatCustomerCreatedAt(customer.CreatedAt),
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

func formatCustomerCreatedAt(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02 15:04:05")
}
