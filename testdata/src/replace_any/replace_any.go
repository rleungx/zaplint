package replace_any

import (
	"fmt"
	"math/cmplx"
	"time"

	zap "go.uber.org/zap"
)

func tests() {
	logger, _ := zap.NewProduction()

	// Negative cases - should trigger lint errors
	logger.Info("message", zap.Any("user_name", "test"))                           // want "replace zap.Any with zap.String"
	logger.Info("message", zap.Any("request_id", 123))                             // want "replace zap.Any with zap.Int"
	logger.Info("message", zap.Any("is_valid", true))                              // want "replace zap.Any with zap.Bool"
	logger.Info("message", zap.Any("average_score", 92.5))                         // want "replace zap.Any with zap.Float64"
	logger.Info("message", zap.Any("process_time", 30*time.Second))                // want "replace zap.Any with zap.Duration"
	logger.Info("message", zap.Any("timestamp", time.Now()))                       // want "replace zap.Any with zap.Time"
	logger.Info("message", zap.Any("is_valid_ptr", new(bool)))                     // want "replace zap.Any with zap.Boolp"
	logger.Info("message", zap.Any("complex128", cmplx.Sqrt(-5+12i)))              // want "replace zap.Any with zap.Complex128"
	logger.Info("message", zap.Any("complex128_ptr", new(complex128)))             // want "replace zap.Any with zap.Complex128p"
	logger.Info("message", zap.Any("complex64", complex64(cmplx.Sqrt(-5+12i))))    // want "replace zap.Any with zap.Complex64"
	logger.Info("message", zap.Any("complex64_ptr", new(complex64)))               // want "replace zap.Any with zap.Complex64p"
	logger.Info("message", zap.Any("average_score_ptr", new(float64)))             // want "replace zap.Any with zap.Float64p"
	logger.Info("message", zap.Any("average_score_32", float32(92.5)))             // want "replace zap.Any with zap.Float32"
	logger.Info("message", zap.Any("average_score_32_ptr", new(float32)))          // want "replace zap.Any with zap.Float32p"
	logger.Info("message", zap.Any("request_id_ptr", new(int)))                    // want "replace zap.Any with zap.Intp"
	logger.Info("message", zap.Any("total_count", int64(1000)))                    // want "replace zap.Any with zap.Int64"
	logger.Info("message", zap.Any("total_count_ptr", new(int64)))                 // want "replace zap.Any with zap.Int64p"
	logger.Info("message", zap.Any("total_count_32", int32(1000)))                 // want "replace zap.Any with zap.Int32"
	logger.Info("message", zap.Any("total_count_32_ptr", new(int32)))              // want "replace zap.Any with zap.Int32p"
	logger.Info("message", zap.Any("total_count_8", int8(100)))                    // want "replace zap.Any with zap.Int8"
	logger.Info("message", zap.Any("total_count_8_ptr", new(int8)))                // want "replace zap.Any with zap.Int8p"
	logger.Info("message", zap.Any("user_name_ptr", new(string)))                  // want "replace zap.Any with zap.Stringp"
	logger.Info("message", zap.Any("request_id", uint(123)))                       // want "replace zap.Any with zap.Uint"
	logger.Info("message", zap.Any("request_id_ptr", new(uint)))                   // want "replace zap.Any with zap.Uintp"
	logger.Info("message", zap.Any("total_count", uint64(1000)))                   // want "replace zap.Any with zap.Uint64"
	logger.Info("message", zap.Any("total_count_ptr", new(uint64)))                // want "replace zap.Any with zap.Uint64p"
	logger.Info("message", zap.Any("total_count_32", uint32(1000)))                // want "replace zap.Any with zap.Uint32"
	logger.Info("message", zap.Any("total_count_32_ptr", new(uint32)))             // want "replace zap.Any with zap.Uint32p"
	logger.Info("message", zap.Any("total_count_16", uint16(1000)))                // want "replace zap.Any with zap.Uint16"
	logger.Info("message", zap.Any("total_count_16_ptr", new(uint16)))             // want "replace zap.Any with zap.Uint16p"
	logger.Info("message", zap.Any("total_count_8", uint8(100)))                   // want "replace zap.Any with zap.Uint8"
	logger.Info("message", zap.Any("total_count_8_ptr", new(uint8)))               // want "replace zap.Any with zap.Uint8p"
	logger.Info("message", zap.Any("uintptr", uintptr(123)))                       // want "replace zap.Any with zap.Uintptr"
	logger.Info("message", zap.Any("uintptr_ptr", new(uintptr)))                   // want "replace zap.Any with zap.Uintptrp"
	logger.Info("message", zap.Any("reflect_data", struct{ Name string }{"test"})) // want "replace zap.Any with zap.Reflect"
	logger.Info("message", zap.Any("timestamp_ptr", new(time.Time)))               // want "replace zap.Any with zap.Timep"
	logger.Info("message", zap.Any("process_time_ptr", new(time.Duration)))        // want "replace zap.Any with zap.Durationp"
	logger.Info("message", zap.Any("named_error", fmt.Errorf("error")))            // want "replace zap.Any with zap.NamedError"

	// Array types
	logger.Info("message", zap.Any("int_array", [3]int{1, 2, 3}))                       // want "replace zap.Any with zap.Ints"
	logger.Info("message", zap.Any("str_array", [2]string{"a", "b"}))                   // want "replace zap.Any with zap.Strings"
	logger.Info("message", zap.Any("bool_array", [2]bool{true, false}))                 // want "replace zap.Any with zap.Bools"
	logger.Info("message", zap.Any("float64_array", [2]float64{1.1, 2.2}))              // want "replace zap.Any with zap.Float64s"
	logger.Info("message", zap.Any("complex128_array", [2]complex128{1 + 2i, 3 + 4i}))  // want "replace zap.Any with zap.Complex128s"
	logger.Info("message", zap.Any("uint_array", [2]uint{1, 2}))                        // want "replace zap.Any with zap.Uints"
	logger.Info("message", zap.Any("uintptr_array", [2]uintptr{1, 2}))                  // want "replace zap.Any with zap.Uintptrs"
	logger.Info("message", zap.Any("time_array", [2]time.Time{time.Now(), time.Now()})) // want "replace zap.Any with zap.Times"
	logger.Info("message", zap.Any("duration_array", [2]time.Duration{time.Second, time.Minute})) // want "replace zap.Any with zap.Durations"
	logger.Info("message", zap.Any("error_array", [2]error{fmt.Errorf("error1"), fmt.Errorf("error2")})) // want "replace zap.Any with zap.Errors"

	// Slice types
	logger.Info("message", zap.Any("int_slice", []int{1, 2, 3}))                       // want "replace zap.Any with zap.Ints"
	logger.Info("message", zap.Any("str_slice", []string{"a", "b"}))                   // want "replace zap.Any with zap.Strings"
	logger.Info("message", zap.Any("bool_slice", []bool{true, false}))                 // want "replace zap.Any with zap.Bools"
	logger.Info("message", zap.Any("float64_slice", []float64{1.1, 2.2}))              // want "replace zap.Any with zap.Float64s"
	logger.Info("message", zap.Any("complex128_slice", []complex128{1 + 2i, 3 + 4i}))  // want "replace zap.Any with zap.Complex128s"
	logger.Info("message", zap.Any("uint_slice", []uint{1, 2}))                        // want "replace zap.Any with zap.Uints"
	logger.Info("message", zap.Any("uintptr_slice", []uintptr{1, 2}))                  // want "replace zap.Any with zap.Uintptrs"
	logger.Info("message", zap.Any("time_slice", []time.Time{time.Now(), time.Now()})) // want "replace zap.Any with zap.Times"
	logger.Info("message", zap.Any("duration_slice", []time.Duration{time.Second, time.Minute})) // want "replace zap.Any with zap.Durations"
	logger.Info("message", zap.Any("error_slice", []error{fmt.Errorf("error1"), fmt.Errorf("error2")})) // want "replace zap.Any with zap.Errors"

}