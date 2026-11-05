package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LoadConfig gives the process environment priority over overrides and the base dotenv file.
func LoadConfig(overrides ...string) (*Config, error) {
	for _, file := range overrides {
		if file != "" {
			if err := loadEnvFile(file, false); err != nil {
				return nil, err
			}
		}
	}

	file := os.Getenv("ENV_FILE")
	optional := file == ""
	if optional {
		file = ".env"
	}
	if err := loadEnvFile(file, optional); err != nil {
		return nil, err
	}

	var parseErrors []error
	cfg := &Config{
		Server: ServerConfig{
			Port:             parseEnv("SV_PORT", strconv.Atoi, &parseErrors),
			WriteTimeout:     parseEnv("SV_WRITE_TIMEOUT", time.ParseDuration, &parseErrors),
			ReadTimeout:      parseEnv("SV_READ_TIMEOUT", time.ParseDuration, &parseErrors),
			IdleTimeout:      parseEnv("SV_IDLE_TIMEOUT", time.ParseDuration, &parseErrors),
			MaxHeaderBytes:   parseEnv("SV_MAX_HEADER_BYTES", strconv.Atoi, &parseErrors),
			AllowMethods:     envList("SV_ALLOW_METHODS"),
			AllowHeaders:     envList("SV_ALLOW_HEADERS"),
			ExposeHeaders:    envList("SV_EXPOSE_HEADERS"),
			AllowCredentials: parseEnv("SV_ALLOW_CREDENTIALS", strconv.ParseBool, &parseErrors),
			MaxAge:           parseEnv("SV_MAX_AGE", time.ParseDuration, &parseErrors),
		},
		Services: ServicesConfig{
			AIHost: os.Getenv("SVC_AI_HOST"),
			AIPort: parseEnv("SVC_AI_PORT", strconv.Atoi, &parseErrors),
		},
		Token: TokenConfig{
			AccessName:       os.Getenv("TK_ACCESS_NAME"),
			RefreshName:      os.Getenv("TK_REFRESH_NAME"),
			GuestName:        os.Getenv("TK_GUEST_NAME"),
			SecretKey:        os.Getenv("TK_SECRET_KEY"),
			APIKeySecret:     os.Getenv("TK_API_KEY_SECRET"),
			APIKeyHeader:     os.Getenv("TK_API_KEY_HEADER"),
			APIKeyPrefix:     os.Getenv("TK_API_KEY_PREFIX"),
			AccessExpiresIn:  parseEnv("TK_ACCESS_EXPIRES_IN", time.ParseDuration, &parseErrors),
			RefreshExpiresIn: parseEnv("TK_REFRESH_EXPIRES_IN", time.ParseDuration, &parseErrors),
		},
		Log: LogConfig{
			Level:    os.Getenv("LOG_LEVEL"),
			Encoding: os.Getenv("LOG_ENCODING"),
		},
		SQL: SQLConfig{
			Host:     os.Getenv("SQL_HOST"),
			Port:     parseEnv("SQL_PORT", strconv.Atoi, &parseErrors),
			User:     os.Getenv("SQL_USER"),
			Password: os.Getenv("SQL_PASSWORD"),
			SSLMode:  os.Getenv("SQL_SSL_MODE"),
			DBName:   os.Getenv("SQL_DB_NAME"),
		},
		NoSQL: NoSQLConfig{
			URI:             os.Getenv("NOSQL_URI"),
			DBName:          os.Getenv("NOSQL_DB_NAME"),
			MaxPoolSize:     parseEnv("NOSQL_MAX_POOL_SIZE", parseUint64, &parseErrors),
			MinPoolSize:     parseEnv("NOSQL_MIN_POOL_SIZE", parseUint64, &parseErrors),
			MaxConnIdleTime: parseEnv("NOSQL_MAX_CONN_IDLE_TIME", time.ParseDuration, &parseErrors),
		},
		Cache: CacheConfig{
			Host:     os.Getenv("CA_HOST"),
			Port:     parseEnv("CA_PORT", strconv.Atoi, &parseErrors),
			Password: os.Getenv("CA_PASSWORD"),
			UseSSL:   parseEnv("CA_USE_SSL", strconv.ParseBool, &parseErrors),
		},
		Storage: StorageConfig{
			Bucket:        os.Getenv("ST_BUCKET"),
			Endpoint:      os.Getenv("ST_ENDPOINT"),
			PublicBaseURL: os.Getenv("ST_PUBLIC_BASE_URL"),
		},
		MQ: MQConfig{
			Host:     os.Getenv("MQ_HOST"),
			Port:     parseEnv("MQ_PORT", strconv.Atoi, &parseErrors),
			User:     os.Getenv("MQ_USER"),
			Password: os.Getenv("MQ_PASSWORD"),
			VHost:    os.Getenv("MQ_VHOST"),
			UseSSL:   parseEnv("MQ_USE_SSL", strconv.ParseBool, &parseErrors),
		},
		Mailer: MailerConfig{
			Host:     os.Getenv("ML_HOST"),
			Port:     parseEnv("ML_PORT", strconv.Atoi, &parseErrors),
			User:     os.Getenv("ML_USER"),
			Password: os.Getenv("ML_PASSWORD"),
		},
		KYC: KYCConfig{
			AccessToken:   os.Getenv("KYC_ACCESS_TOKEN"),
			WorkFlowID:    os.Getenv("KYC_WORKFLOW_ID"),
			WebhookSecret: os.Getenv("KYC_WEBHOOK_SECRET"),
		},
	}
	if err := errors.Join(parseErrors...); err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadEnvFile(path string, optional bool) error {
	err := godotenv.Load(path)
	if err == nil || (optional && errors.Is(err, os.ErrNotExist)) {
		return nil
	}
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return pathErr
	}
	// Parser errors can contain part of the file, including credentials.
	return fmt.Errorf("invalid dotenv syntax in %q", path)
}

func parseEnv[T any](key string, parse func(string) (T, error), errs *[]error) T {
	var zero T
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return zero
	}
	value, err := parse(raw)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("invalid value for %s", key))
		return zero
	}
	return value
}

func parseUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func envList(key string) []string {
	var list []string
	for part := range strings.SplitSeq(os.Getenv(key), ",") {
		if item := strings.TrimSpace(part); item != "" {
			list = append(list, item)
		}
	}
	return list
}
