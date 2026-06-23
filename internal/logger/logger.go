package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

// MustInitialize инициализирует синглтон логера с необходимым уровнем логирования.
func MustInitialize(level string) {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	// устанавливаем синглтон
	Log = zl
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&lw, r)

		Log.Info("completed HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", time.Since(start)),
			zap.Int("status", lw.statusCode),
			zap.Int("size", lw.size),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}
