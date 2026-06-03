package utils

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	logLevel := zerolog.InfoLevel

	if os.Getenv("GIN_MODE") == "debug" {
		logLevel = zerolog.TraceLevel
	}

	if os.Getenv("LOG_OUTPUT") == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.SetGlobalLevel(logLevel)
}

type LoggerConfiguration struct {
	Name      string
	SkipPaths []string
}

type LogContent struct {
	Name    string
	Path    string
	Latency time.Duration
	Method  string
	Status  int
	Message string
}

func LoggerMiddleware(config LoggerConfiguration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		start := time.Now()

		ctx.Next()

		if Contains(config.SkipPaths, path) {
			return
		}

		message := ctx.Errors.String()

		if message == "" {
			message = "Request"
		}

		logRequest(&LogContent{
			Name:    config.Name,
			Path:    path,
			Latency: time.Since(start),
			Method:  ctx.Request.Method,
			Status:  ctx.Writer.Status(),
			Message: message,
		})
	}
}

func logRequest(content *LogContent) {
	var logger *zerolog.Event

	switch {
	case content.Status >= 500:
		{
			logger = log.Error()
		}
	case content.Status >= 400 && content.Status < 500:
		{
			logger = log.Warn()
		}
	default:
		logger = log.Trace()
	}

	logger.
		Str("name", content.Name).
		Str("method", content.Method).
		Str("path", content.Path).
		Dur("time", content.Latency).
		Int("status", content.Status).
		Msg(content.Message)
}

type ContextLogger = func(ctx context.Context) *zerolog.Logger

func Logger(ctx context.Context) *zerolog.Logger {
	logger := log.Logger
	return &logger
}

func InitServiceLogger(loggerTag string) ContextLogger {
	logger := log.With().Str("TAG", loggerTag).Logger()

	return func(ctx context.Context) *zerolog.Logger {
		return &logger
	}
}
