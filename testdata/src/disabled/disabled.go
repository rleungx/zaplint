package disabled

import "go.uber.org/zap"

func disabled(logger *zap.Logger) {
	logger.Info("lower", zap.Any("BadKey", 1))
}
