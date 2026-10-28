package utils

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"math/bits"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateSlug(t *testing.T) {
	got := GenerateSlug("Hello World!")
	if got != "hello-world" {
		t.Fatalf("expected slug 'hello-world', got %q", got)
	}
}

func TestCalculateStayDays(t *testing.T) {
	tests := []struct {
		name string
		from time.Time
		to   time.Time
		want int64
	}{
		{
			name: "same day rounds up to 1",
			from: time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC),
			want: 1,
		},
		{
			name: "under 24h rounds up to 1",
			from: time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 4, 20, 9, 59, 59, 0, time.UTC),
			want: 1,
		},
		{
			name: "exactly 24h is 1",
			from: time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC),
			want: 1,
		},
		{
			name: "over 24h rounds up",
			from: time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 4, 20, 10, 0, 1, 0, time.UTC),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateStayDays(tt.from, tt.to); got != tt.want {
				t.Fatalf("CalculateStayDays(%s, %s) = %d, want %d", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestNormalizeDateOnly(t *testing.T) {
	in := time.Date(2026, 4, 19, 23, 59, 58, 123_000_000, time.FixedZone("ICT", 7*60*60))
	got := NormalizeDateOnly(in)

	if got.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", got.Location())
	}
	if got.Year() != 2026 || got.Month() != time.April || got.Day() != 19 {
		t.Fatalf("expected 2026-04-19, got %s", got.Format("2006-01-02"))
	}
	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
		t.Fatalf("expected 00:00:00.000000000, got %s", got.Format(time.RFC3339Nano))
	}
}

func TestGenerateOTP(t *testing.T) {
	otp := GenerateOTP(6)
	if len(otp) != 6 {
		t.Fatalf("expected length 6, got %d", len(otp))
	}
	for _, ch := range otp {
		if ch < '0' || ch > '9' {
			t.Fatalf("expected numeric otp, got %q", otp)
		}
	}
}

func TestGenerateRandom(t *testing.T) {
	out := GenerateRandom(12)
	raw, err := base64.RawURLEncoding.DecodeString(out)
	if err != nil {
		t.Fatalf("expected base64 string, got error %v", err)
	}
	if len(raw) != 12 {
		t.Fatalf("expected decoded length 12, got %d", len(raw))
	}
}

func TestGenerateBase58ID(t *testing.T) {
	out := GenerateBase58ID(8)
	if out == "" {
		t.Fatalf("expected non-empty base58 id")
	}
}

func TestGenerateUUID(t *testing.T) {
	id := GenerateUUID()
	if _, err := uuid.Parse(id); err != nil {
		t.Fatalf("expected valid uuid, got %q", id)
	}
}

func TestSHA256Hash(t *testing.T) {
	got := SHA256Hash("hello")
	const expected = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	hashed, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if err := VerifyPassword("secret123", hashed); err != nil {
		t.Fatalf("verify should pass: %v", err)
	}
	if err := VerifyPassword("wrong", hashed); err == nil {
		t.Fatalf("verify should fail for wrong password")
	}
}

func TestExtractRootDomain(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"localhost:8080", "localhost"},
		{"127.0.0.1:80", "127.0.0.1"},
		{"example.com", ".example.com"},
		{"a.b.example.com", ".example.com"},
	}
	for _, c := range cases {
		if got := ExtractRootDomain(c.in); got != c.want {
			t.Fatalf("ExtractRootDomain(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRoutingAndQueueNames(t *testing.T) {
	if got := GetRoutingKey("a", "b"); got != "a.b" {
		t.Fatalf("expected routing key a.b, got %q", got)
	}
	if got := GetQueueName("x", "y"); got != "x.y" {
		t.Fatalf("expected queue name x.y, got %q", got)
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  int64
	}{
		{name: "int", value: int(12), want: 12},
		{name: "int64", value: int64(13), want: 13},
		{name: "int32", value: int32(14), want: 14},
		{name: "uint", value: uint(15), want: 15},
		{name: "uint64", value: uint64(16), want: 16},
		{name: "float64", value: float64(17.9), want: 17},
		{name: "float32", value: float32(18.9), want: 18},
		{name: "json number", value: json.Number("19"), want: 19},
		{name: "json number invalid", value: json.Number("19.5"), want: 0},
		{name: "string", value: "20", want: 20},
		{name: "string invalid", value: "abc", want: 0},
		{name: "default", value: true, want: 0},
		{name: "uint64 overflow", value: uint64(math.MaxInt64) + 1, want: 0},
		{name: "float64 overflow", value: math.MaxFloat64, want: 0},
		{name: "float64 underflow", value: -math.MaxFloat64, want: 0},
		{name: "float32 overflow", value: float32(math.MaxFloat32), want: 0},
		{name: "float32 underflow", value: -float32(math.MaxFloat32), want: 0},
	}

	if bits.UintSize == 64 {
		tests = append(tests, struct {
			name  string
			value any
			want  int64
		}{
			name:  "uint overflow",
			value: uint(math.MaxInt64) + 1,
			want:  0,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInt64(tt.value)
			if got != tt.want {
				t.Fatalf("ParseInt64(%v) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	if got := ParseString("hello"); got != "hello" {
		t.Fatalf("expected hello, got %q", got)
	}

	if got := ParseString(json.Number("123")); got != "123" {
		t.Fatalf("expected 123, got %q", got)
	}

	if got := ParseString(123); got != "" {
		t.Fatalf("expected empty string for unsupported type, got %q", got)
	}
}
