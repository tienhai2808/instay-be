package response

import (
	"net/http"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/gin-gonic/gin"
)

const (
	ActionSuccess = "Operation successful"
)

func APIResponse(c *gin.Context, status, code int, message string, data any) {
	c.JSON(status, response.APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OKResponse(c *gin.Context, data any) {
	APIResponse(c, http.StatusOK, constants.CodeSuccess, ActionSuccess, data)
}
