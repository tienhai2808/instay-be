package config

import "github.com/spf13/viper"

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("server.port", "SV_PORT")
	_ = viper.BindEnv("server.write_timeout", "SV_WRITE_TIMEOUT")
	_ = viper.BindEnv("server.read_timeout", "SV_READ_TIMEOUT")
	_ = viper.BindEnv("server.idle_timeout", "SV_IDLE_TIMEOUT")
	_ = viper.BindEnv("server.allow_methods", "SV_ALLOW_METHODS")
	_ = viper.BindEnv("server.allow_headers", "SV_ALLOW_HEADERS")
	_ = viper.BindEnv("server.expose_headers", "SV_EXPOSE_HEADERS")
	_ = viper.BindEnv("server.allow_credentials", "SV_ALLOW_CREDENTIALS")
	_ = viper.BindEnv("server.max_age", "SV_MAX_AGE")
	_ = viper.BindEnv("server.max_header_bytes", "SV_MAX_HEADER_BYTES")

	_ = viper.BindEnv("services.ai_host", "SVC_AI_HOST")
	_ = viper.BindEnv("services.ai_port", "SVC_AI_PORT")

	_ = viper.BindEnv("sql.host", "SQL_HOST")
	_ = viper.BindEnv("sql.port", "SQL_PORT")
	_ = viper.BindEnv("sql.user", "SQL_USER")
	_ = viper.BindEnv("sql.password", "SQL_PASSWORD")
	_ = viper.BindEnv("sql.ssl_mode", "SQL_SSL_MODE")
	_ = viper.BindEnv("sql.db_name", "SQL_DB_NAME")

	_ = viper.BindEnv("nosql.uri", "NOSQL_URI")
	_ = viper.BindEnv("nosql.max_pool_size", "NOSQL_MAX_POOL_SIZE")
	_ = viper.BindEnv("nosql.min_pool_size", "NOSQL_MIN_POOL_SIZE")
	_ = viper.BindEnv("nosql.max_conn_idle_time", "NOSQL_MAX_CONN_IDLE_TIME")
	_ = viper.BindEnv("nosql.db_name", "NOSQL_DB_NAME")

	_ = viper.BindEnv("log.level", "LOG_LEVEL")
	_ = viper.BindEnv("log.encoding", "LOG_ENCODING")

	_ = viper.BindEnv("cache.host", "CA_HOST")
	_ = viper.BindEnv("cache.port", "CA_PORT")
	_ = viper.BindEnv("cache.password", "CA_PASSWORD")
	_ = viper.BindEnv("cache.use_ssl", "CA_USE_SSL")

	_ = viper.BindEnv("token.access_name", "TK_ACCESS_NAME")
	_ = viper.BindEnv("token.refresh_name", "TK_REFRESH_NAME")
	_ = viper.BindEnv("token.guest_name", "TK_GUEST_NAME")
	_ = viper.BindEnv("token.access_expires_in", "TK_ACCESS_EXPIRES_IN")
	_ = viper.BindEnv("token.refresh_expires_in", "TK_REFRESH_EXPIRES_IN")
	_ = viper.BindEnv("token.secret_key", "TK_SECRET_KEY")
	_ = viper.BindEnv("token.api_key_secret", "TK_API_KEY_SECRET")
	_ = viper.BindEnv("token.api_key_header", "TK_API_KEY_HEADER")
	_ = viper.BindEnv("token.api_key_prefix", "TK_API_KEY_PREFIX")

	_ = viper.BindEnv("storage.bucket", "ST_BUCKET")

	_ = viper.BindEnv("mailer.host", "ML_HOST")
	_ = viper.BindEnv("mailer.port", "ML_PORT")
	_ = viper.BindEnv("mailer.user", "ML_USER")
	_ = viper.BindEnv("mailer.password", "ML_PASSWORD")

	_ = viper.BindEnv("mq.host", "MQ_HOST")
	_ = viper.BindEnv("mq.port", "MQ_PORT")
	_ = viper.BindEnv("mq.user", "MQ_USER")
	_ = viper.BindEnv("mq.password", "MQ_PASSWORD")
	_ = viper.BindEnv("mq.vhost", "MQ_VHOST")
	_ = viper.BindEnv("mq.use_ssl", "MQ_USE_SSL")

	_ = viper.BindEnv("kyc.access_token", "KYC_ACCESS_TOKEN")
	_ = viper.BindEnv("kyc.workflow_id", "KYC_WORKFLOW_ID")
	_ = viper.BindEnv("kyc.webhook_secret", "KYC_WEBHOOK_SECRET")

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
