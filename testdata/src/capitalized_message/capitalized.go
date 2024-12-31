package capitalized

import "go.uber.org/zap"

func tests() {
	logger, _ := zap.NewProduction()
	// Positive cases - should pass
	logger.Info("Message should be capitalized")
	logger.Error("Error message should be capitalized")
	logger.Warn("Warning message should be capitalized")
	logger.Debug("Debug message should be capitalized")
	logger.DPanic("DPanic message should be capitalized")
	logger.Panic("Panic message should be capitalized")
	logger.Fatal("Fatal message should be capitalized")

	// Negative cases - should trigger lint errors
	logger.Info("message should be capitalized")          // want "message 'message should be capitalized' should be capitalized"
	logger.Error("error message should be capitalized")   // want "message 'error message should be capitalized' should be capitalized"
	logger.Warn("warning message should be capitalized")  // want "message 'warning message should be capitalized' should be capitalized"
	logger.Debug("debug message should be capitalized")   // want "message 'debug message should be capitalized' should be capitalized"
	logger.DPanic("dpanic message should be capitalized") // want "message 'dpanic message should be capitalized' should be capitalized"
	logger.Panic("panic message should be capitalized")   // want "message 'panic message should be capitalized' should be capitalized"
	logger.Fatal("fatal message should be capitalized")   // want "message 'fatal message should be capitalized' should be capitalized"
}
