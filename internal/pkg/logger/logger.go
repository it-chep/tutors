package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey string

const contextLoggerKey loggerKey = "clk"

// Logger - кастомная обертка над zap.Logger
type Logger struct {
	zap *zap.Logger
}

// New ...
func New() *Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	return &Logger{
		zap: zap.New(core),
	}
}

// ContextWithLogger прокинуть логер в контекст
func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, l)
}

// fromContext взять логгер из контекста
func fromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(contextLoggerKey).(*Logger)
	if ok {
		return l
	}
	return nil
}

// Message записать в лог
func Message(ctx context.Context, format string) {
	l := fromContext(ctx)
	if l != nil {
		l.zap.Info(format)
	}
}

// Error записать ошибку
func Error(ctx context.Context, format string, err error) {
	l := fromContext(ctx)
	if l != nil {
		l.zap.Error(format, zap.Error(err))
	}
}
