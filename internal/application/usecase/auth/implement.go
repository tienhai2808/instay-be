package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

const (
	maxVerifyForgotPasswordAttempts = 5
	maxResendForgotPasswordAttempts = 3
)

type authUseCaseImpl struct {
	cfg              config.TokenConfig
	txPro            port.TxProvider
	logPro           port.LogProvider
	tokenPro         port.TokenProvider
	cachePro         port.CacheProvider
	mqPro            port.MessageQueueProvider
	storagePro       port.StorageProvider
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthUseCase(
	cfg config.TokenConfig,
	txPro port.TxProvider,
	logPro port.LogProvider,
	tokenPro port.TokenProvider,
	cachePro port.CacheProvider,
	mqPro port.MessageQueueProvider,
	storPro port.StorageProvider,
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) AuthUseCase {
	return &authUseCaseImpl{
		cfg,
		txPro,
		logPro,
		tokenPro,
		cachePro,
		mqPro,
		storPro,
		userRepo,
		refreshTokenRepo,
	}
}

func (u *authUseCaseImpl) Login(ctx context.Context, schema string, req request.LoginRequest) (*response.UserResponse, string, string, error) {
	user, err := u.userRepo.FindByUsernameWithDepartmentDetails(ctx, schema, req.Username)
	if err != nil {
		u.logPro.Error("Find user by username failed",
			port.String("schema", schema),
			port.String("username", req.Username),
			port.Error(err))
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", exceptions.ErrLoginFailed
	}

	if err = utils.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, "", "", exceptions.ErrLoginFailed
	}

	cacheKey := fmt.Sprintf("%s:user_version:%d", schema, user.ID)
	tokenVersion, err := u.cachePro.GetInt(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get token version failed", port.String("schema", schema), port.Error(err))
		return nil, "", "", err
	}
	if tokenVersion == 0 {
		if err = u.cachePro.SetString(ctx, cacheKey, "1", 0); err != nil {
			u.logPro.Error("Save token version failed", port.String("schema", schema), port.Error(err))
			return nil, "", "", err
		}
		tokenVersion = 1
	}

	accessToken, err := u.tokenPro.GenerateStaffToken(user.ID, user.Role, tokenVersion, user.DepartmentID, u.cfg.AccessExpiresIn)
	if err != nil {
		u.logPro.Error("Generate access token failed", port.String("schema", schema), port.Error(err))
		return nil, "", "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		u.logPro.Error("Generate refresh token failed", port.String("schema", schema), port.Error(err))
		return nil, "", "", err
	}

	token := &model.RefreshToken{
		UserID:    user.ID,
		Token:     utils.SHA256Hash(refreshToken),
		ExpiresAt: time.Now().Add(u.cfg.RefreshExpiresIn),
	}

	if err := u.refreshTokenRepo.Create(ctx, schema, token); err != nil {
		u.logPro.Error("Create refresh token failed", port.String("schema", schema), port.Error(err))
		return nil, "", "", err
	}

	return mapper.ToUserResponse(user), accessToken, refreshToken, nil
}

func (u *authUseCaseImpl) Logout(ctx context.Context, schema, accessToken, refreshToken string) error {
	hashedToken := utils.SHA256Hash(refreshToken)

	if err := u.refreshTokenRepo.Delete(ctx, schema, hashedToken); err != nil {
		if errors.Is(err, exceptions.ErrInvalidUser) {
			return err
		}
		u.logPro.Error("Delete refresh token by token failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if accessToken != "" {
		_, _, _, _, accessTTL, err := u.tokenPro.ParseStaffToken(accessToken)
		if err == nil && accessTTL > 0 {
			cacheKey := fmt.Sprintf("%s:black_list:%s", schema, accessToken)
			if err := u.cachePro.SetString(ctx, cacheKey, "1", accessTTL); err != nil {
				u.logPro.Error("Save black list failed", port.String("schema", schema), port.Error(err))
			}
		}
	}

	return nil
}

func (u *authUseCaseImpl) RefreshToken(ctx context.Context, schema, refreshToken string) (string, string, error) {
	hashedToken := utils.SHA256Hash(refreshToken)

	token, err := u.refreshTokenRepo.FindByToken(ctx, schema, hashedToken)
	if err != nil {
		u.logPro.Error("Find refresh token by token failed", port.String("schema", schema), port.Error(err))
		return "", "", nil
	}
	if token == nil || token.ExpiresAt.Before(time.Now()) {
		return "", "", exceptions.ErrInvalidUser
	}

	user, err := u.userRepo.FindByID(ctx, schema, token.UserID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Error(err))
		return "", "", nil
	}
	if user == nil || !user.IsActive {
		return "", "", exceptions.ErrInvalidUser
	}

	cacheKey := fmt.Sprintf("%s:user_version:%d", schema, user.ID)
	tokenVersion, err := u.cachePro.GetInt(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get token version failed", port.String("schema", schema), port.Error(err))
		return "", "", err
	}
	if tokenVersion == 0 {
		return "", "", exceptions.ErrInvalidUser
	}

	newAccessToken, err := u.tokenPro.GenerateStaffToken(user.ID, user.Role, tokenVersion, user.DepartmentID, u.cfg.AccessExpiresIn)
	if err != nil {
		u.logPro.Error("Generate access token failed", port.String("schema", schema), port.Error(err))
		return "", "", err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		u.logPro.Error("Generate refresh token failed", port.String("schema", schema), port.Error(err))
		return "", "", err
	}

	updateData := map[string]any{
		"token":      utils.SHA256Hash(newRefreshToken),
		"expires_at": time.Now().Add(u.cfg.RefreshExpiresIn),
	}

	if err := u.refreshTokenRepo.Update(ctx, schema, hashedToken, updateData); err != nil {
		u.logPro.Error("Update refresh token failed", port.String("schema", schema), port.Error(err))
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (u *authUseCaseImpl) UserInfo(ctx context.Context, schema string, userID int64) (*response.UserResponse, error) {
	user, err := u.userRepo.FindByIDWithDepartmentDetails(ctx, schema, userID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, exceptions.ErrInvalidUser
	}

	return mapper.ToUserResponse(user), nil
}

func (u *authUseCaseImpl) ChangePassword(ctx context.Context, schema string, userID int64, req request.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(ctx, schema, userID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}
	if user == nil {
		return exceptions.ErrInvalidUser
	}

	if err = utils.VerifyPassword(req.OldPassword, user.Password); err != nil {
		return exceptions.ErrInvalidPassword
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		u.logPro.Error("Hash password failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateData := map[string]any{
			"password":      hashedPassword,
			"updated_by_id": user.ID,
		}

		if err = u.userRepo.Update(ctx, schema, userID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrUserNotFound) {
				return exceptions.ErrInvalidUser
			}
			u.logPro.Error("Update password failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
			return err
		}

		if err := u.refreshTokenRepo.DeleteAllByUserID(ctx, schema, userID); err != nil {
			u.logPro.Error("Delete all refresh token by user id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("%suser_version:%d", schema, user.ID)
	if err = u.cachePro.Increment(ctx, cacheKey); err != nil {
		u.logPro.Error("Increase token version failed", port.String("schema", schema), port.Error(err))
	}

	return nil
}

func (u *authUseCaseImpl) ForgotPassword(ctx context.Context, schema, email string) (string, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, schema, email)
	if err != nil {
		u.logPro.Error("Check user by email failed", port.String("schema", schema), port.String("email", email), port.Error(err))
		return "", err
	}
	if !exists {
		return "", exceptions.ErrEmailDoesNotExist
	}

	otp := utils.GenerateOTP(6)
	forgotPasswordToken := utils.GenerateUUID()

	forgData := data.ForgotPasswordData{
		Email:          email,
		Otp:            otp,
		VerifyAttempts: 0,
		ResendAttempts: 0,
	}

	bytes, err := json.Marshal(forgData)
	if err != nil {
		u.logPro.Error("Json marshal forgot password data failed", port.String("schema", schema), port.Error(err))
		return "", err
	}

	cacheKey := fmt.Sprintf("%s:forgot_password:%s", schema, forgotPasswordToken)
	if err = u.cachePro.SetObject(ctx, cacheKey, bytes, 5*time.Minute); err != nil {
		u.logPro.Error("Save forgot password data failed", port.String("schema", schema), port.Error(err))
		return "", err
	}

	emailMsg := data.AuthEmailMessage{
		Schema:  schema,
		To:      email,
		Subject: "Xác thực quên mật khẩu tại Instay",
		Otp:     otp,
	}

	go func(msg data.AuthEmailMessage) {
		body, err := json.Marshal(msg)
		if err != nil {
			u.logPro.Error("Json marshal failed", port.String("schema", schema), port.Error(err))
		}

		routingKey := utils.GetRoutingKey(schema, port.RoutingKeyAuthEmail)

		if err = u.mqPro.PublishMessage(port.ExchangeEmail, routingKey, body); err != nil {
			u.logPro.Error("Publish auth email message failed",
				port.String("schema", schema),
				port.String("email", email),
				port.Error(err))
		}
	}(emailMsg)

	return forgotPasswordToken, nil
}

func (u *authUseCaseImpl) VerifyForgotPassword(ctx context.Context, schema string, req request.VerifyForgotPasswordRequest) (string, error) {
	cacheKey := fmt.Sprintf("%s:forgot_password:%s", schema, req.ForgotPasswordToken)
	bytes, err := u.cachePro.GetObject(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get forgot password data failed", port.String("schema", schema), port.Error(err))
		return "", err
	}
	if bytes == nil {
		return "", exceptions.ErrInvalidToken
	}

	var forgData data.ForgotPasswordData
	if err = json.Unmarshal(bytes, &forgData); err != nil {
		u.logPro.Error("Json unmarshal forgot password data failed", port.String("schema", schema), port.Error(err))
		return "", nil
	}

	if forgData.VerifyAttempts >= maxVerifyForgotPasswordAttempts || forgData.ResendAttempts >= maxResendForgotPasswordAttempts {
		if err = u.cachePro.Del(ctx, cacheKey); err != nil {
			u.logPro.Error("Delete forgot password data failed", port.String("schema", schema), port.Error(err))
		}
		return "", exceptions.ErrTooManyAttempts
	}

	forgData.VerifyAttempts++
	bytes, err = json.Marshal(forgData)
	if err != nil {
		u.logPro.Error("Json marshal forgot password data failed", port.String("schema", schema), port.Error(err))
		return "", err
	}

	newTTL, err := u.cachePro.TTL(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get forgot password ttl failed", port.String("schema", schema), port.Error(err))
		return "", err
	}
	if newTTL <= 0 {
		return "", exceptions.ErrInvalidToken
	}

	if err = u.cachePro.SetObject(ctx, cacheKey, bytes, newTTL); err != nil {
		u.logPro.Error("Save forgot password data failed", port.String("schema", schema), port.Error(err))
		return "", err
	}

	if forgData.Otp != req.Otp {
		return "", exceptions.ErrInvalidOTP
	}

	resetPasswordToken := utils.GenerateUUID()
	key := fmt.Sprintf("%s:reset_password:%s", schema, resetPasswordToken)

	if err = u.cachePro.SetString(ctx, key, forgData.Email, 3*time.Minute); err != nil {
		u.logPro.Error("Save email reset password failed", port.String("schema", schema), port.Error(err))
		return "", err
	}

	if err = u.cachePro.Del(ctx, cacheKey); err != nil {
		u.logPro.Error("Delete forgot password data failed", port.String("schema", schema), port.Error(err))
	}

	return resetPasswordToken, nil
}

func (u *authUseCaseImpl) ResendForgotPassword(ctx context.Context, schema string, req request.ResendForgotPasswordRequest) error {
	cacheKey := fmt.Sprintf("%s:forgot_password:%s", schema, req.ForgotPasswordToken)
	bytes, err := u.cachePro.GetObject(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get forgot password data failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if bytes == nil {
		return exceptions.ErrInvalidToken
	}

	var forgData data.ForgotPasswordData
	if err = json.Unmarshal(bytes, &forgData); err != nil {
		u.logPro.Error("Json unmarshal forgot password data failed", port.String("schema", schema), port.Error(err))
		return nil
	}

	if forgData.VerifyAttempts >= maxVerifyForgotPasswordAttempts || forgData.ResendAttempts >= maxResendForgotPasswordAttempts {
		if err = u.cachePro.Del(ctx, cacheKey); err != nil {
			u.logPro.Error("Delete forgot password data failed", port.String("schema", schema), port.Error(err))
		}
		return exceptions.ErrTooManyAttempts
	}

	otp := utils.GenerateOTP(6)
	forgData.ResendAttempts++
	forgData.Otp = otp

	bytes, err = json.Marshal(forgData)
	if err != nil {
		u.logPro.Error("Json marshal forgot password data failed", port.String("schema", schema), port.Error(err))
		return err
	}

	newTTL, err := u.cachePro.TTL(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get forgot password ttl failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if newTTL <= 0 {
		return exceptions.ErrInvalidToken
	}

	if err = u.cachePro.SetObject(ctx, cacheKey, bytes, newTTL); err != nil {
		u.logPro.Error("Save forgot password data failed", port.String("schema", schema), port.Error(err))
		return err
	}

	emailMsg := data.AuthEmailMessage{
		Schema:  schema,
		To:      forgData.Email,
		Subject: "Xác thực quên mật khẩu tại Instay",
		Otp:     otp,
	}

	go func(msg data.AuthEmailMessage) {
		body, err := json.Marshal(msg)
		if err != nil {
			u.logPro.Error("Json marshal failed", port.String("schema", schema), port.Error(err))
		}

		routingKey := utils.GetRoutingKey(schema, port.RoutingKeyAuthEmail)

		if err = u.mqPro.PublishMessage(port.ExchangeEmail, routingKey, body); err != nil {
			u.logPro.Error("Publish auth email message failed",
				port.String("schema", schema),
				port.String("email", forgData.Email),
				port.Error(err))
		}
	}(emailMsg)

	return nil
}

func (u *authUseCaseImpl) ResetPassword(ctx context.Context, schema string, req request.ResetPasswordRequest) error {
	cacheKey := fmt.Sprintf("%s:reset_password:%s", schema, req.ResetPasswordToken)
	email, err := u.cachePro.GetString(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get email reset password failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if email == "" {
		return exceptions.ErrInvalidToken
	}

	user, err := u.userRepo.FindByEmail(ctx, schema, email)
	if err != nil {
		u.logPro.Error("Find user by email failed", port.String("schema", schema), port.String("email", email), port.Error(err))
		return err
	}
	if user == nil {
		return exceptions.ErrUnauth
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		u.logPro.Error("Hash password failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateData := map[string]any{
			"password":      hashedPassword,
			"updated_by_id": user.ID,
		}

		if err = u.userRepo.Update(ctx, schema, user.ID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrUserNotFound) {
				return exceptions.ErrUnauth
			}
			u.logPro.Error("Update password failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if err := u.refreshTokenRepo.DeleteAllByUserID(ctx, schema, user.ID); err != nil {
			u.logPro.Error("Delete all refresh token by user id failed", port.String("schema", schema), port.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err = u.cachePro.Del(ctx, cacheKey); err != nil {
		u.logPro.Error("Delete reset password failed", port.String("schema", schema), port.Error(err))
	}

	cacheKey = fmt.Sprintf("%s:user_version:%d", schema, user.ID)
	if err = u.cachePro.Increment(ctx, cacheKey); err != nil {
		u.logPro.Error("Increase token version failed", port.String("schema", schema), port.Error(err))
	}

	return nil
}

func (u *authUseCaseImpl) UpdateInfo(ctx context.Context, schema string, userID int64, req request.UpdateInfoRequest) error {
	user, err := u.userRepo.FindByID(ctx, schema, userID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}
	if user == nil {
		return exceptions.ErrInvalidUser
	}

	updateData := map[string]any{
		"email":         req.Email,
		"phone":         req.Phone,
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
		"avatar":        req.Avatar,
		"updated_by_id": userID,
	}

	if err := u.userRepo.Update(ctx, schema, userID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrUserNotFound) {
			return exceptions.ErrInvalidUser
		}
		if errors.Is(err, exceptions.ErrEmailAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrPhoneAlreadyExists) {
			return err
		}
		u.logPro.Error("Update user failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	oldAvatar := user.Avatar
	newAvatar := req.Avatar
	shouldDelete := false
	if oldAvatar != nil {
		if newAvatar == nil {
			shouldDelete = true
		} else if *oldAvatar != *newAvatar {
			shouldDelete = true
		}
	}

	if shouldDelete {
		deleteFilesMsg := data.DeleteFileMessage{
			Schema: schema,
			Keys:   []string{*oldAvatar},
		}

		go func(msg data.DeleteFileMessage) {
			body, err := json.Marshal(msg)
			if err != nil {
				u.logPro.Error("Json marshal failed", port.String("schema", schema), port.Error(err))
			}

			if err = u.mqPro.PublishMessage(port.ExchangeFile, port.RoutingKeyDeleteFile, body); err != nil {
				u.logPro.Error("Publish delete files message failed", port.String("schema", schema), port.Error(err))
			}
		}(deleteFilesMsg)
	}

	return nil
}

func generateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	rawToken := base64.RawURLEncoding.EncodeToString(randomBytes)

	return rawToken, nil
}
