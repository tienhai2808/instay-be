package utils

import (
	cRand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/bcrypt"
)

func CalculateStayDays(fromDate, toDate time.Time) int64 {
	stayDurationHours := toDate.Sub(fromDate).Hours()
	return max(int64(math.Ceil(stayDurationHours/24)), int64(1))
}

func GenerateSlug(str string) string {
	return slug.Make(str)
}

func NormalizeDateOnly(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func GenerateOTP(length uint8) string {
	const chars = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		n, _ := cRand.Int(cRand.Reader, big.NewInt(int64(len(chars))))
		otp[i] = chars[n.Int64()]
	}
	return string(otp)
}

func GenerateRandom(length uint8) string {
	b := make([]byte, length)
	_, _ = cRand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func GenerateBase58ID(size int) string {
	b := make([]byte, size)
	_, _ = cRand.Read(b)
	return base58.Encode(b)
}

func GenerateUUID() string {
	return uuid.NewString()
}

func SHA256Hash(str string) string {
	hashArray := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashArray[:])
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ExtractRootDomain(host string) string {
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	if host == "localhost" || !strings.Contains(host, ".") {
		return host
	}

	parts := strings.Split(host, ".")
	if len(parts) == 4 {
		isIP := true
		for _, part := range parts {
			for _, ch := range part {
				if ch < '0' || ch > '9' {
					isIP = false
					break
				}
			}
		}
		if isIP {
			return host
		}
	}

	if len(parts) >= 2 {
		return "." + parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	return host
}

func GetRoutingKey(schema, routingKeyBase string) string {
	return fmt.Sprintf("%s.%s", schema, routingKeyBase)
}

func GetQueueName(schema, queueName string) string {
	return fmt.Sprintf("%s.%s", schema, queueName)
}

func ParseInt64(value any) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case int32:
		return int64(v)
	case uint:
		if uint64(v) > math.MaxInt64 {
			return 0
		}
		return int64(v)
	case uint64:
		if v > math.MaxInt64 {
			return 0
		}
		return int64(v)
	case float64:
		if v > float64(math.MaxInt64) || v < float64(math.MinInt64) {
			return 0
		}
		return int64(v)
	case float32:
		v64 := float64(v)
		if v64 > float64(math.MaxInt64) || v64 < float64(math.MinInt64) {
			return 0
		}
		return int64(v)
	case json.Number:
		id, err := v.Int64()
		if err != nil {
			return 0
		}
		return id
	case string:
		id, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		return id
	default:
		return 0
	}
}

func ParseString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	default:
		return ""
	}
}
