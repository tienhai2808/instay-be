package port

import (
	"math"
	"time"
)

type LogFieldType int

const (
	LogFieldTypeAny LogFieldType = iota
	LogFieldTypeString
	LogFieldTypeInt
	LogFieldTypeInt64
	LogFieldTypeUint
	LogFieldTypeUint32
	LogFieldTypeUint64
	LogFieldTypeBool
	LogFieldTypeFloat64
	LogFieldTypeDuration
	LogFieldTypeTime
	LogFieldTypeError
)

type LogField struct {
	Key      string
	Type     LogFieldType
	Integer  int64
	Unsigned uint64
	String   string
	Value    any
}

type LogProvider interface {
	Debug(msg string, fields ...LogField)

	Info(msg string, fields ...LogField)

	Warn(msg string, fields ...LogField)

	Error(msg string, fields ...LogField)

	Panic(msg string, fields ...LogField)

	Fatal(msg string, fields ...LogField)

	With(fields ...LogField) LogProvider

	Named(name string) LogProvider

	Sync() error
}

func Any(key string, value any) LogField {
	return LogField{Key: key, Type: LogFieldTypeAny, Value: value}
}
func String(key, value string) LogField {
	return LogField{Key: key, Type: LogFieldTypeString, String: value}
}

func Int(key string, value int) LogField {
	return LogField{Key: key, Type: LogFieldTypeInt, Integer: int64(value)}
}

func Int64(key string, value int64) LogField {
	return LogField{Key: key, Type: LogFieldTypeInt64, Integer: value}
}

func Uint(key string, value uint) LogField {
	return LogField{Key: key, Type: LogFieldTypeUint, Unsigned: uint64(value)}
}

func Uint32(key string, value uint32) LogField {
	return LogField{Key: key, Type: LogFieldTypeUint32, Unsigned: uint64(value)}
}

func Uint64(key string, value uint64) LogField {
	return LogField{Key: key, Type: LogFieldTypeUint64, Unsigned: value}
}

func Float64(key string, value float64) LogField {
	return LogField{Key: key, Type: LogFieldTypeFloat64, Unsigned: math.Float64bits(value)}
}

func Duration(key string, value time.Duration) LogField {
	return LogField{Key: key, Type: LogFieldTypeDuration, Integer: int64(value)}
}

func Time(key string, value time.Time) LogField {
	return LogField{Key: key, Type: LogFieldTypeTime, Integer: value.UnixNano(), Value: value.Location()}
}

func Error(err error) LogField {
	return LogField{Key: "error", Type: LogFieldTypeError, Value: err}
}

func NamedError(key string, err error) LogField {
	return LogField{Key: key, Type: LogFieldTypeError, Value: err}
}

func Bool(key string, val bool) LogField {
	i := int64(0)
	if val {
		i = 1
	}
	return LogField{Key: key, Type: LogFieldTypeBool, Integer: i}
}
