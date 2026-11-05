package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigEnvPrecedence(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("ENV_FILE", "")
	t.Setenv("SD_NAME", "tenant")
	t.Setenv("ST_ENDPOINT", "")
	for _, key := range []string{"SV_PORT", "SV_READ_TIMEOUT", "SV_ALLOW_METHODS", "AWS_ACCESS_KEY_ID"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}

	base := "SV_PORT=8080\nSV_READ_TIMEOUT=30s\nSV_ALLOW_METHODS=GET, POST\nST_ENDPOINT=https://example.com\nAWS_ACCESS_KEY_ID=test-key\n"
	if err := os.WriteFile(".env", []byte(base), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(".env.seeder", []byte("SV_PORT=8081\nSD_NAME=admin\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(".env.seeder")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 8081 || cfg.Server.ReadTimeout != 30*time.Second ||
		!reflect.DeepEqual(cfg.Server.AllowMethods, []string{"GET", "POST"}) {
		t.Fatalf("unexpected server settings: %+v", cfg.Server)
	}
	if cfg.Storage.Endpoint != "" || os.Getenv("SD_NAME") != "tenant" || os.Getenv("AWS_ACCESS_KEY_ID") != "test-key" {
		t.Fatal("dotenv did not preserve process values or export SDK settings")
	}
}

func TestLoadConfigInvalidDotenv(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("ENV_FILE", "")
	if _, err := LoadConfig(); err != nil {
		t.Fatalf("missing default .env should be allowed: %v", err)
	}

	t.Setenv("ENV_FILE", "custom.env")
	if _, err := LoadConfig(); !os.IsNotExist(err) {
		t.Fatalf("missing explicit env file should fail: %v", err)
	}
	if err := os.WriteFile("custom.env", []byte("SQL_PASSWORD='example-sensitive-value\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(); err == nil || strings.Contains(err.Error(), "example-sensitive-value") {
		t.Fatal("expected a syntax error without the dotenv value")
	}
}
