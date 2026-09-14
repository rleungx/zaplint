package capitalized

import (
	"log/slog"
	"testing"

	"go.uber.org/zap"
)

type customMessage string

func (customMessage) String() string { return "Upper" }

func regression(logger *zap.Logger, t *testing.T) {
	const good = "Started"
	const bad = "failed"
	logger.Info(good)
	logger.Info(`Started`)
	logger.Info("\x53tarted")
	logger.Info((bad))                            // want "message 'failed' should be capitalized"
	logger.Info("fail" + "ed")                    // want "message 'failed' should be capitalized"
	logger.Info(`failed`)                         // want "message 'failed' should be capitalized"
	logger.Log(zap.InfoLevel, bad)                // want "message 'failed' should be capitalized"
	logger.Check(zap.InfoLevel, bad)              // want "message 'failed' should be capitalized"
	(*zap.Logger).Info(logger, bad)               // want "message 'failed' should be capitalized"
	(*zap.Logger).Log(logger, zap.InfoLevel, bad) // want "message 'failed' should be capitalized"
	(logger.Info)(bad)                            // want "message 'failed' should be capitalized"

	sugar := logger.Sugar()
	sugar.Infow(bad, "key", 1)                 // want "message 'failed' should be capitalized"
	sugar.Infof("failed: %d", 1)               // want "message 'failed: %d' should be capitalized"
	sugar.Infoln(bad, 1)                       // want "message 'failed' should be capitalized"
	sugar.Log(zap.InfoLevel, bad)              // want "message 'failed' should be capitalized"
	sugar.Logf(zap.InfoLevel, "failed: %d", 1) // want "message 'failed: %d' should be capitalized"
	sugar.Logw(zap.InfoLevel, bad, "key", 1)   // want "message 'failed' should be capitalized"
	sugar.Logln(zap.InfoLevel, bad, 1)         // want "message 'failed' should be capitalized"
	(*zap.SugaredLogger).Infow(sugar, bad)     // want "message 'failed' should be capitalized"
	sugar.Infof("%s", good)                    // Dynamic prefix: no capitalization claim.
	sugar.Info("", good)
	sugar.Infof("", good)
	sugar.Infow(good)
	sugar.Info(customMessage("lower"))

	// Similarly named methods from unrelated packages are outside zaplint.
	t.Error("failed")
	slog.Info("started")
}
