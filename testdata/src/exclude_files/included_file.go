package exclude_files

import (
	"go.uber.org/zap"
)

func includedTests() {
	logger := zap.NewProduction()
	// This file should be included in analysis
	logger.Info("message", zap.String("userName", "test")) // want "key 'userName' should be in snake_case"
}
