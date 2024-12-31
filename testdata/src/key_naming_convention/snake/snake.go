package snake

import (
	"go.uber.org/zap"
)

func tests() {
	logger, _ := zap.NewProduction()
	// Positive cases - should pass
	logger.Info("message", zap.String("user_name", "test"))
	logger.Info("message", zap.Int("request_id", 123))
	logger.Info("message", zap.String("first_name", "John"))
	logger.Info("message", zap.Int64("total_count", 1000))
	logger.Info("message", zap.Float64("average_score", 92.5))
	logger.Info("message", zap.Duration("process_time", 30))
	logger.Info("message", zap.Bool("is_valid", true))
	logger.Info("message", zap.String("api_version", "v1"))
	logger.Info("message", zap.Int("retry_count", 3))
	logger.Info("message", zap.String("error_message", "timeout"))
	logger.Info("message", zap.Int64("memory_usage", 1024))
	logger.Info("message", zap.Float64("response_time_ms", 150.5))

	// Negative cases - should trigger lint errors
	logger.Info("message", zap.String("userName", "test"))   // want "key 'userName' should be in snake_case"
	logger.Info("message", zap.Int("RequestID", 123))        // want "key 'RequestID' should be in snake_case"
	logger.Info("message", zap.Bool("isValid", true))        // want "key 'isValid' should be in snake_case"
	logger.Info("message", zap.String("apiVersion", "v2"))   // want "key 'apiVersion' should be in snake_case"
	logger.Info("message", zap.Int("RetryCount", 5))         // want "key 'RetryCount' should be in snake_case"
	logger.Info("message", zap.Float64("responseTime", 200)) // want "key 'responseTime' should be in snake_case"

	// Additional cases for other conventions
	logger.Info("message", zap.String("user-name", "test"))  // want "key 'user-name' should be in snake_case"
	logger.Info("message", zap.String("UserName", "test"))   // want "key 'UserName' should be in snake_case"
	logger.Info("message", zap.String("USER_NAME", "test"))  // want "key 'USER_NAME' should be in snake_case"
	logger.Info("message", zap.String("API-Version", "v3"))  // want "key 'API-Version' should be in snake_case"
	logger.Info("message", zap.String("HTTPResponse", "ok")) // want "key 'HTTPResponse' should be in snake_case"
	logger.Info("message", zap.Int("MAX_RETRY", 10))         // want "key 'MAX_RETRY' should be in snake_case"
}
