package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role         model.UserRole `json:"role"`
	TokenVersion int            `json:"token_version"`
	DepartmentID *int64         `json:"department_id"`
}

type tokenProviderImpl struct {
	cfg config.TokenConfig
}

func NewTokenProvider(cfg config.TokenConfig) port.TokenProvider {
	return &tokenProviderImpl{cfg}
}

func (p *tokenProviderImpl) GenerateStaffToken(userID int64, role model.UserRole, tokenVersion int, departmentID *int64, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:         role,
		TokenVersion: tokenVersion,
		DepartmentID: departmentID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(p.cfg.SecretKey))
}

func (p *tokenProviderImpl) ParseStaffToken(tokenStr string) (int64, model.UserRole, int, *int64, time.Duration, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
		}
		return []byte(p.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, "", 0, nil, 0, err
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, "", 0, nil, 0, err
	}

	ttl := time.Until(claims.ExpiresAt.Time)

	return userID, claims.Role, claims.TokenVersion, claims.DepartmentID, ttl, nil
}

func (p *tokenProviderImpl) GenerateGuestToken(bookingRoomID int64, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(bookingRoomID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(p.cfg.SecretKey))
}

func (p *tokenProviderImpl) ParseGuestToken(tokenStr string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
		}
		return []byte(p.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	bookingRoomID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return bookingRoomID, nil
}
