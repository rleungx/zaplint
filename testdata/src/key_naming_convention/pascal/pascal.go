package pascal

import (
	"go.uber.org/zap"
)

func tests() {
	logger, _ := zap.NewProduction()
	// Positive cases - should pass
	logger.Info("message", zap.String("UserName", "test"))
	logger.Info("message", zap.Int("RequestId", 123))
	logger.Info("message", zap.String("FirstName", "John"))
	logger.Info("message", zap.Int64("TotalCount", 1000))
	logger.Info("message", zap.Float64("AverageScore", 92.5))
	logger.Info("message", zap.Duration("ProcessTime", 30))
	logger.Info("message", zap.Bool("IsValid", true))
	logger.Info("message", zap.String("ApiVersion", "v1"))
	logger.Info("message", zap.Int("RetryCount", 3))
	logger.Info("message", zap.String("ErrorMessage", "timeout"))
	logger.Info("message", zap.Int64("MemoryUsage", 1024))
	logger.Info("message", zap.Float64("ResponseTimeMs", 150.5))

	// Negative cases - should trigger lint errors
	logger.Info("message", zap.String("user_name", "test"))   // want "key 'user_name' should be in PascalCase"
	logger.Info("message", zap.Int("requestID", 123))         // want "key 'requestID' should be in PascalCase"
	logger.Info("message", zap.Bool("is_valid", true))        // want "key 'is_valid' should be in PascalCase"
	logger.Info("message", zap.String("api_version", "v2"))   // want "key 'api_version' should be in PascalCase"
	logger.Info("message", zap.Int("retryCount", 5))          // want "key 'retryCount' should be in PascalCase"
	logger.Info("message", zap.Float64("response_time", 200)) // want "key 'response_time' should be in PascalCase"

	// Additional cases for other conventions
	logger.Info("message", zap.String("user-name", "test"))  // want "key 'user-name' should be in PascalCase"
	logger.Info("message", zap.String("userName", "test"))   // want "key 'userName' should be in PascalCase"
	logger.Info("message", zap.String("USER_NAME", "test"))  // want "key 'USER_NAME' should be in PascalCase"
	logger.Info("message", zap.String("apiVersion", "v3"))   // want "key 'apiVersion' should be in PascalCase"
	logger.Info("message", zap.String("httpResponse", "ok")) // want "key 'httpResponse' should be in PascalCase"
	logger.Info("message", zap.Int("maxRetry", 10))          // want "key 'maxRetry' should be in PascalCase"
}
