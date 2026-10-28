package config

import "time"

type ServerConfig struct {
	Port             int           `mapstructure:"port"`
	HeaderTenant     string        `mapstructure:"header_tenant"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes   int           `mapstructure:"max_header_bytes"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

type ServicesConfig struct {
	AIHost string `mapstructure:"ai_host"`
	AIPort int    `mapstructure:"ai_port"`
}

type TokenConfig struct {
	AccessName       string        `mapstructure:"access_name"`
	RefreshName      string        `mapstructure:"refresh_name"`
	GuestName        string        `mapstructure:"guest_name"`
	SecretKey        string        `mapstructure:"secret_key"`
	APIKeySecret     string        `mapstructure:"api_key_secret"`
	APIKeyHeader     string        `mapstructure:"api_key_header"`
	APIKeyPrefix     string        `mapstructure:"api_key_prefix"`
	AccessExpiresIn  time.Duration `mapstructure:"access_expires_in"`
	RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

type SQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
	DBName   string `mapstructure:"db_name"`
}

type NoSQLConfig struct {
	URI             string        `mapstructure:"uri"`
	DBName          string        `mapstructure:"db_name"`
	MaxPoolSize     uint64        `mapstructure:"max_pool_size"`
	MinPoolSize     uint64        `mapstructure:"min_pool_size"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

type CacheConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}

type StorageConfig struct {
	Bucket string `mapstructure:"bucket"`
}

type KYCConfig struct {
	AccessToken   string `mapstructure:"access_token"`
	WorkFlowID    string `mapstructure:"workflow_id"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

type MQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}

type MailerConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Services ServicesConfig `mapstructure:"services"`
	Token    TokenConfig    `mapstructure:"token"`
	Log      LogConfig      `mapstructure:"log"`
	SQL      SQLConfig      `mapstructure:"sql"`
	NoSQL    NoSQLConfig    `mapstructure:"nosql"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Storage  StorageConfig  `mapstructure:"storage"`
	MQ       MQConfig       `mapstructure:"mq"`
	Mailer   MailerConfig   `mapstructure:"mailer"`
	KYC      KYCConfig      `mapstructure:"kyc"`
}
