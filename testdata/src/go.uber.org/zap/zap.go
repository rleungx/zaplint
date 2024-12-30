package zap

import (
	"fmt"
	"time"
)

type Logger struct{}

func NewProduction() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}

type Field struct{}

func Binary(key string, val []byte) Field {
	return Field{}
}

func Bool(key string, val bool) Field {
	return Field{}
}

func Boolp(key string, val *bool) Field {
	return Field{}
}

func ByteString(key string, val []byte) Field {
	return Field{}
}

func Complex128(key string, val complex128) Field {
	return Field{}
}

func Complex128p(key string, val *complex128) Field {
	return Field{}
}

func Complex64(key string, val complex64) Field {
	return Field{}
}

func Complex64p(key string, val *complex64) Field {
	return Field{}
}

func Float64(key string, val float64) Field {
	return Field{}
}

func Float64p(key string, val *float64) Field {
	return Field{}
}

func Float32(key string, val float32) Field {
	return Field{}
}

func Float32p(key string, val *float32) Field {
	return Field{}
}

func Int(key string, val int) Field {
	return Field{}
}

func Intp(key string, val *int) Field {
	return Field{}
}

func Int64(key string, val int64) Field {
	return Field{}
}

func Int64p(key string, val *int64) Field {
	return Field{}
}

func Int32(key string, val int32) Field {
	return Field{}
}

func Int32p(key string, val *int32) Field {
	return Field{}
}

func Int8(key string, val int8) Field {
	return Field{}
}

func Int8p(key string, val *int8) Field {
	return Field{}
}

func String(key, val string) Field {
	return Field{}
}

func Stringp(key string, val *string) Field {
	return Field{}
}

func Uint(key string, val uint) Field {
	return Field{}
}

func Uintp(key string, val *uint) Field {
	return Field{}
}

func Uint64(key string, val uint64) Field {
	return Field{}
}

func Uint64p(key string, val *uint64) Field {
	return Field{}
}

func Uint32(key string, val uint32) Field {
	return Field{}
}

func Uint32p(key string, val *uint32) Field {
	return Field{}
}

func Uint16(key string, val uint16) Field {
	return Field{}
}

func Uint16p(key string, val *uint16) Field {
	return Field{}
}

func Uint8(key string, val uint8) Field {
	return Field{}
}

func Uint8p(key string, val *uint8) Field {
	return Field{}
}

func Uintptr(key string, val uintptr) Field {
	return Field{}
}

func Uintptrp(key string, val *uintptr) Field {
	return Field{}
}

func Reflect(key string, val interface{}) Field {
	return Field{}
}

func Stringer(key string, val fmt.Stringer) Field {
	return Field{}
}

func Time(key string, val time.Time) Field {
	return Field{}
}

func Timep(key string, val *time.Time) Field {
	return Field{}
}

func Duration(key string, val time.Duration) Field {
	return Field{}
}

func Durationp(key string, val *time.Duration) Field {
	return Field{}
}

func NamedError(key string, err error) Field {
	return Field{}
}

func Any(key string, val interface{}) Field {
	return Field{}
}
