package zap

type Logger struct{}

func NewProduction() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, fields ...Field) {}

type Field struct{}

func String(key, val string) Field {
	return Field{}
}

func Int(key string, val int) Field {
	return Field{}
}

func Bool(key string, val bool) Field {
	return Field{}
}

func Int64(key string, val int64) Field {
	return Field{}
}

func Float64(key string, val float64) Field {
	return Field{}
}

func Duration(key string, val int) Field {
	return Field{}
}
