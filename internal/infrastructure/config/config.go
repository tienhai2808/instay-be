package config

import "time"

type ServerConfig struct {
	Port int

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	IdleTimeout      time.Duration
	MaxHeaderBytes   int
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

type ServicesConfig struct {
	AIHost string
	AIPort int
}

type TokenConfig struct {
	AccessName       string
	RefreshName      string
	GuestName        string
	SecretKey        string
	APIKeySecret     string
	APIKeyHeader     string
	APIKeyPrefix     string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

type LogConfig struct {
	Level    string
	Encoding string
}

type SQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	SSLMode  string
	DBName   string
}

type NoSQLConfig struct {
	URI             string
	DBName          string
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
}

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	UseSSL   bool
}

type StorageConfig struct {
	Bucket        string
	Endpoint      string
	PublicBaseURL string
}

type KYCConfig struct {
	AccessToken   string
	WorkFlowID    string
	WebhookSecret string
}

type MQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
	UseSSL   bool
}

type MailerConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Token    TokenConfig
	Log      LogConfig
	SQL      SQLConfig
	NoSQL    NoSQLConfig
	Cache    CacheConfig
	Storage  StorageConfig
	MQ       MQConfig
	Mailer   MailerConfig
	KYC      KYCConfig
}
