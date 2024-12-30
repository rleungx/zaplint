package exclude_files

import (
	"go.uber.org/zap"
)

func excludedTestsFoo() {
	logger := zap.NewProduction()
	// This file should be excluded from analysis
	logger.Info("message", zap.String("userName", "test"))
}
