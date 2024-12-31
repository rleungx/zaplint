package camel

import (
	"go.uber.org/zap"
)

func tests() {
	logger, _ := zap.NewProduction()
	// Positive cases - should pass
	logger.Info("message", zap.String("userName", "test"))
	logger.Info("message", zap.Int("requestId", 123))
	logger.Info("message", zap.String("firstName", "John"))
	logger.Info("message", zap.Int64("totalCount", 1000))
	logger.Info("message", zap.Float64("averageScore", 92.5))
	logger.Info("message", zap.Duration("processTime", 30))
	logger.Info("message", zap.Bool("isValid", true))
	logger.Info("message", zap.String("apiVersion", "v1"))
	logger.Info("message", zap.Int("retryCount", 3))
	logger.Info("message", zap.String("errorMessage", "timeout"))
	logger.Info("message", zap.Int64("memoryUsage", 1024))
	logger.Info("message", zap.Float64("responseTimeMs", 150.5))

	// Negative cases - should trigger lint errors
	logger.Info("message", zap.String("user_name", "test"))   // want "key 'user_name' should be in camelCase"
	logger.Info("message", zap.Int("RequestID", 123))         // want "key 'RequestID' should be in camelCase"
	logger.Info("message", zap.Bool("is_valid", true))        // want "key 'is_valid' should be in camelCase"
	logger.Info("message", zap.String("api_version", "v2"))   // want "key 'api_version' should be in camelCase"
	logger.Info("message", zap.Int("RetryCount", 5))          // want "key 'RetryCount' should be in camelCase"
	logger.Info("message", zap.Float64("response_time", 200)) // want "key 'response_time' should be in camelCase"

	// Additional cases for other conventions
	logger.Info("message", zap.String("user-name", "test"))  // want "key 'user-name' should be in camelCase"
	logger.Info("message", zap.String("UserName", "test"))   // want "key 'UserName' should be in camelCase"
	logger.Info("message", zap.String("USER_NAME", "test"))  // want "key 'USER_NAME' should be in camelCase"
	logger.Info("message", zap.String("API-Version", "v3"))  // want "key 'API-Version' should be in camelCase"
	logger.Info("message", zap.String("HTTPResponse", "ok")) // want "key 'HTTPResponse' should be in camelCase"
	logger.Info("message", zap.Int("MAX_RETRY", 10))         // want "key 'MAX_RETRY' should be in camelCase"
}
