package middleware

import (
	"net/http"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/gin-gonic/gin"
)

type APIKeyMiddleware struct {
	apiKeyPro port.APIKeyProvider
}

func NewAPIKeyMiddleware(
	apiKeyPro port.APIKeyProvider,
) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		apiKeyPro: apiKeyPro,
	}
}

func (m *APIKeyMiddleware) IsValidAPIKey(c *gin.Context) {
	headerName := m.apiKeyPro.HeaderName()
	rawValue := c.GetHeader(headerName)

	key, err := m.apiKeyPro.Extract(rawValue)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	schema, err := m.apiKeyPro.Parse(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
			Code:    exceptions.ErrUnauth.Code,
			Message: exceptions.ErrUnauth.Message,
		})
		return
	}

	c.Set(CtxSchema, schema)
	c.Next()
}
