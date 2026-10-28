package apikey

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

const (
	defaultHeaderName = "X-API-Key"
	keyPrefix         = "ik"
	keyVersion        = "v1"
	randomPartLength  = 24
)

var (
	errInvalidAPIKey = errors.New("invalid api key")
)

type providerImpl struct {
	secret []byte
	header string
	prefix string
}

func NewAPIKeyProvider(cfg config.TokenConfig) port.APIKeyProvider {
	secret := strings.TrimSpace(cfg.APIKeySecret)
	if secret == "" {
		secret = cfg.SecretKey
	}

	header := strings.TrimSpace(cfg.APIKeyHeader)
	if header == "" {
		header = defaultHeaderName
	}

	return &providerImpl{
		secret: []byte(secret),
		header: header,
		prefix: strings.TrimSpace(cfg.APIKeyPrefix),
	}
}

func (p *providerImpl) Generate(schema string) (string, error) {
	schema = strings.TrimSpace(schema)
	if schema == "" {
		return "", errInvalidAPIKey
	}

	nonce, err := generateNonce()
	if err != nil {
		return "", err
	}

	signature := p.sign(schema, nonce)
	return fmt.Sprintf("%s.%s.%s.%s.%s", keyPrefix, keyVersion, schema, nonce, signature), nil
}

func (p *providerImpl) Parse(key string) (string, error) {
	parts := strings.Split(strings.TrimSpace(key), ".")
	if len(parts) != 5 {
		return "", errInvalidAPIKey
	}
	if parts[0] != keyPrefix || parts[1] != keyVersion {
		return "", errInvalidAPIKey
	}

	schema := strings.TrimSpace(parts[2])
	nonce := strings.TrimSpace(parts[3])
	signature := strings.TrimSpace(parts[4])
	if schema == "" || nonce == "" || signature == "" {
		return "", errInvalidAPIKey
	}

	expectedSig := p.sign(schema, nonce)
	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return "", errInvalidAPIKey
	}

	return schema, nil
}

func (p *providerImpl) Extract(headerValue string) (string, error) {
	raw := strings.TrimSpace(headerValue)
	if raw == "" {
		return "", errInvalidAPIKey
	}

	if p.prefix == "" {
		return raw, nil
	}

	if !strings.HasPrefix(raw, p.prefix) {
		return "", errInvalidAPIKey
	}

	key := strings.TrimSpace(strings.TrimPrefix(raw, p.prefix))
	if key == "" {
		return "", errInvalidAPIKey
	}

	return key, nil
}

func (p *providerImpl) HeaderName() string {
	return p.header
}

func (p *providerImpl) HeaderPrefix() string {
	return p.prefix
}

func (p *providerImpl) sign(schema, nonce string) string {
	payload := fmt.Sprintf("%s.%s.%s", keyVersion, schema, nonce)
	mac := hmac.New(sha256.New, p.secret)
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func generateNonce() (string, error) {
	bytes := make([]byte, randomPartLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
