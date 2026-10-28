package handler

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	fileUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/file"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileUC fileUC.FileUseCase
}

func NewFileHandler(fileUC fileUC.FileUseCase) *FileHandler {
	return &FileHandler{fileUC}
}

func (h *FileHandler) UploadPresignedURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.UploadPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	presignedURLs, err := h.fileUC.CreateUploadURLs(ctx, schema, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"presigned_urls": presignedURLs,
	})
}

func (h *FileHandler) ViewPresignedURLs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.ViewPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	presignedURLs, err := h.fileUC.CreateViewURLs(ctx, schema, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"presigned_urls": presignedURLs,
	})
}
