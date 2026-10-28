package usecase

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type AuthUseCase interface {
	Login(ctx context.Context, schema string, req request.LoginRequest) (*response.UserResponse, string, string, error)

	Logout(ctx context.Context, schema, accessToken, refreshToken string) error

	RefreshToken(ctx context.Context, schema, refreshToken string) (string, string, error)

	UserInfo(ctx context.Context, schema string, userID int64) (*response.UserResponse, error)

	ChangePassword(ctx context.Context, schema string, userID int64, req request.ChangePasswordRequest) error

	ForgotPassword(ctx context.Context, schema, email string) (string, error)

	VerifyForgotPassword(ctx context.Context, schema string, req request.VerifyForgotPasswordRequest) (string, error)

	ResendForgotPassword(ctx context.Context, schema string, req request.ResendForgotPasswordRequest) error

	ResetPassword(ctx context.Context, schema string, req request.ResetPasswordRequest) error

	UpdateInfo(ctx context.Context, schema string, userID int64, req request.UpdateInfoRequest) error
}
