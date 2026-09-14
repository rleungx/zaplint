package snake

import (
	"errors"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type object struct{}
type invalidKey string

func (object) MarshalLogObject(zapcore.ObjectEncoder) error { return nil }

func regression(logger *zap.Logger, unknown any, err error) {
	const good = "user_name"
	const bad = "UserName"
	zap.String(good, "v")
	zap.String(`user_name`, "v")
	zap.String("user\x5fname", "v")
	zap.String((bad), "v")                      // want "key 'UserName' should be in snake_case"
	zap.String("User"+"Name", "v")              // want "key 'UserName' should be in snake_case"
	zap.Int16(bad, 1)                           // want "key 'UserName' should be in snake_case"
	zap.Int16p(bad, nil)                        // want "key 'UserName' should be in snake_case"
	zap.Object(bad, object{})                   // want "key 'UserName' should be in snake_case"
	zap.Dict(bad)                               // want "key 'UserName' should be in snake_case"
	zap.Namespace(bad)                          // want "key 'UserName' should be in snake_case"
	zap.Stack(bad)                              // want "key 'UserName' should be in snake_case"
	zap.StackSkip(bad, 1)                       // want "key 'UserName' should be in snake_case"
	zap.Objects[object](bad, []object{{}})      // want "key 'UserName' should be in snake_case"
	zap.ObjectValues[object](bad, []object{{}}) // want "key 'UserName' should be in snake_case"

	sugar := logger.Sugar()
	sugar.Infow("Message", bad, 1)                        // want "key 'UserName' should be in snake_case"
	sugar.Logw(zap.InfoLevel, "Message", bad, 1)          // want "key 'UserName' should be in snake_case"
	sugar.With(bad, 1)                                    // want "key 'UserName' should be in snake_case"
	sugar.WithLazy(bad, 1)                                // want "key 'UserName' should be in snake_case"
	(*zap.SugaredLogger).Infow(sugar, "Message", bad, 1)  // want "key 'UserName' should be in snake_case"
	sugar.Infow("Message", zap.String(good, "v"), bad, 1) // want "key 'UserName' should be in snake_case"
	sugar.Infow("Message", 1, "NotAKey", bad, 1)          // want "key 'UserName' should be in snake_case"
	sugar.Infow("Message", unknown, "CouldBeValue", "CouldBeKey", 1)
	sugar.Infow("Message", err, "CouldBeValue", "CouldBeKey", 1)
	sugar.Infow("Message", bad) // Dangling keys are ignored by zap.
	sugar.Infow("Message", good, "NotAKey")
	sugar.Infow("Message", invalidKey("NotAKey"), 1, bad, 1) // want "key 'UserName' should be in snake_case"
	sugar.Infow("Message", errors.New("err"), "CouldBeValue", "CouldBeKey", 1)
	slog.Any(bad, 1)
}
