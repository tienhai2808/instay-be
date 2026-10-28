package zap

import (
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"go.uber.org/zap"
)

var messageNormalizers = []struct {
	pattern     *regexp.Regexp
	replacement string
}{
	{regexp.MustCompile(`(?i)\bhttp\b`), "HTTP"},
	{regexp.MustCompile(`(?i)\bjson\b`), "JSON"},
	{regexp.MustCompile(`(?i)\bgrpc\b`), "gRPC"},
	{regexp.MustCompile(`(?i)\bimap\b`), "IMAP"},
	{regexp.MustCompile(`(?i)\binbox\b`), "INBOX"},
	{regexp.MustCompile(`(?i)\bws\b`), "WS"},
	{regexp.MustCompile(`(?i)\bapi\b`), "API"},
	{regexp.MustCompile(`(?i)\bttl\b`), "TTL"},
	{regexp.MustCompile(`(?i)\bkyc\b`), "KYC"},
	{regexp.MustCompile(`(?i)\brabbitmq\b`), "RabbitMQ"},
	{regexp.MustCompile(`(?i)\bredis\b`), "Redis"},
	{regexp.MustCompile(`(?i)\back\b`), "ACK"},
	{regexp.MustCompile(`(?i)\bnack\b`), "NACK"},
	{regexp.MustCompile(`(?i)\bids\b`), "IDs"},
	{regexp.MustCompile(`(?i)\bid\b`), "ID"},
}

type logProviderImpl struct {
	log *zap.Logger
}

func NewLogProvider(log *zap.Logger) port.LogProvider {
	return &logProviderImpl{log}
}

func (p *logProviderImpl) Debug(msg string, fields ...port.LogField) {
	p.log.Debug(normalizeMessage(msg), toZapFields(fields)...)
}

func (p *logProviderImpl) Info(msg string, fields ...port.LogField) {
	p.log.Info(normalizeMessage(msg), toZapFields(fields)...)
}

func (p *logProviderImpl) Warn(msg string, fields ...port.LogField) {
	p.log.Warn(normalizeMessage(msg), toZapFields(fields)...)
}

func (p *logProviderImpl) Error(msg string, fields ...port.LogField) {
	p.log.Error(normalizeMessage(msg), toZapFields(fields)...)
}

func (p *logProviderImpl) Panic(msg string, fields ...port.LogField) {
	p.log.Panic(normalizeMessage(msg), toZapFields(fields)...)
}

func (p *logProviderImpl) Fatal(msg string, fields ...port.LogField) {
	p.log.Fatal(normalizeMessage(msg), toZapFields(fields)...)
}

func (p *logProviderImpl) With(fields ...port.LogField) port.LogProvider {
	return &logProviderImpl{p.log.With(toZapFields(fields)...)}
}

func (p *logProviderImpl) Named(name string) port.LogProvider {
	return &logProviderImpl{p.log.Named(name)}
}

func (p *logProviderImpl) Sync() error {
	return p.log.Sync()
}

func normalizeMessage(msg string) string {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return msg
	}

	if msg == strings.ToUpper(msg) {
		msg = strings.ToLower(msg)
	}

	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])
	msg = string(runes)

	for _, normalizer := range messageNormalizers {
		msg = normalizer.pattern.ReplaceAllString(msg, normalizer.replacement)
	}

	return msg
}

func toZapFields(fields []port.LogField) []zap.Field {
	if len(fields) == 0 {
		return nil
	}
	out := make([]zap.Field, len(fields))
	for i, f := range fields {
		out[i] = toZapField(f)
	}
	return out
}

func toZapField(f port.LogField) zap.Field {
	switch f.Type {
	case port.LogFieldTypeString:
		return zap.String(f.Key, f.String)

	case port.LogFieldTypeInt, port.LogFieldTypeInt64:
		return zap.Int64(f.Key, f.Integer)

	case port.LogFieldTypeUint, port.LogFieldTypeUint32, port.LogFieldTypeUint64:
		return zap.Uint64(f.Key, f.Unsigned)

	case port.LogFieldTypeBool:
		return zap.Bool(f.Key, f.Integer == 1)

	case port.LogFieldTypeFloat64:
		return zap.Float64(f.Key, math.Float64frombits(f.Unsigned))

	case port.LogFieldTypeDuration:
		return zap.Duration(f.Key, time.Duration(f.Integer))

	case port.LogFieldTypeTime:
		if loc, ok := f.Value.(*time.Location); ok {
			return zap.Time(f.Key, time.Unix(0, f.Integer).In(loc))
		}
		if v, ok := f.Value.(time.Time); ok {
			return zap.Time(f.Key, v)
		}
		return zap.Time(f.Key, time.Unix(0, f.Integer))

	case port.LogFieldTypeError:
		err, ok := f.Value.(error)
		if !ok {
			return zap.Skip()
		}
		if f.Key == "error" {
			return zap.Error(err)
		}
		return zap.NamedError(f.Key, err)

	default:
		return zap.Any(f.Key, f.Value)
	}
}
