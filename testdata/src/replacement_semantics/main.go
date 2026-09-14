package main

import (
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ids []int
type integer int
type byteAlias = byte
type timeAlias = time.Time

type masked []int

func (masked) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	enc.AppendString("masked")
	return nil
}

type stringer []int

func (stringer) String() string { return "custom" }

type objectError struct{}

func (objectError) Error() string { return "error" }
func (objectError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("custom", "object")
	return nil
}

func emit(label string, field zap.Field) {
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})
	buf, err := enc.EncodeEntry(zapcore.Entry{}, []zapcore.Field{field})
	if err != nil {
		panic(err)
	}
	defer buf.Free()
	fmt.Printf("%s: %s", label, buf.String())
}

func main() {
	emit("int", zap.Any("value", 1))                                            // want "replace zap.Any with zap.Int"
	emit("int16", zap.Any("value", int16(1)))                                   // want "replace zap.Any with zap.Int16"
	emit("int16 ptr", zap.Any("value", new(int16)))                             // want "replace zap.Any with zap.Int16p"
	emit("nil ptr", zap.Any("value", (*int16)(nil)))                            // want "replace zap.Any with zap.Int16p"
	emit("bool", zap.Any("value", true))                                        // want "replace zap.Any with zap.Bool"
	emit("string", zap.Any("value", "text"))                                    // want "replace zap.Any with zap.String"
	emit("bytes", zap.Any("value", []byte{65}))                                 // want "replace zap.Any with zap.Binary"
	emit("uint8s", zap.Any("value", []uint8{65}))                               // want "replace zap.Any with zap.Binary"
	emit("byte alias", zap.Any("value", []byteAlias{65}))                       // want "replace zap.Any with zap.Binary"
	emit("nil bytes", zap.Any("value", []byte(nil)))                            // want "replace zap.Any with zap.Binary"
	emit("empty bytes", zap.Any("value", []byte{}))                             // want "replace zap.Any with zap.Binary"
	emit("ints", zap.Any("value", []int{1, 2}))                                 // want "replace zap.Any with zap.Ints"
	emit("nil ints", zap.Any("value", []int(nil)))                              // want "replace zap.Any with zap.Ints"
	emit("runes", zap.Any("value", []rune{'a'}))                                // want "replace zap.Any with zap.Int32s"
	emit("duration", zap.Any("value", time.Second))                             // want "replace zap.Any with zap.Duration"
	emit("durations", zap.Any("value", []time.Duration{time.Second}))           // want "replace zap.Any with zap.Durations"
	emit("time", zap.Any("value", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))) // want "replace zap.Any with zap.Time"
	emit("nil time ptr", zap.Any("value", (*timeAlias)(nil)))                   // want "replace zap.Any with zap.Timep"
	emit("struct", zap.Any("value", struct{ Count int }{1}))                    // want "replace zap.Any with zap.Reflect"
	emit("errors", zap.Any("value", []error{nil, fmt.Errorf("failed")}))        // want "replace zap.Any with zap.Errors"

	// These values must retain Any: replacing by underlying type is not safe.
	emit("array", zap.Any("value", [3]int{1, 2, 3}))
	emit("empty array", zap.Any("value", [0]int{}))
	emit("nested bytes", zap.Any("value", [][]byte{{65}}))
	emit("nested uint8s", zap.Any("value", [][]uint8{{65}}))
	emit("named nil slice", zap.Any("value", ids(nil)))
	emit("named scalar", zap.Any("value", integer(1)))
	emit("masked", zap.Any("value", masked{1, 2}))
	emit("stringer", zap.Any("value", stringer{1, 2}))
	emit("embedded stringer", zap.Any("value", struct{ fmt.Stringer }{stringer{1}}))
	var err error
	emit("nil error", zap.Any("value", err))
	err = objectError{}
	emit("object error", zap.Any("value", err))
	emit("error", zap.Any("value", fmt.Errorf("failed")))
	var unknown any = masked{1}
	emit("interface", zap.Any("value", unknown))

	// This is not a zap field and must not receive a replacement suggestion.
	fmt.Println(slog.Any("value", 1))
}
