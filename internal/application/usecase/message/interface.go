package message

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type MessageUseCase interface {
	Create(ctx context.Context, schema string, req request.CreateMessageRequest) (*response.MessageResponse, error)

	UpdateRead(ctx context.Context, schema string, req request.UpdateMessageReadRequest) error

	UpdateReaction(ctx context.Context, schema string, req request.UpdateMessageReactionRequest) (*response.UpdateMessageReactionResponse, error)
}
