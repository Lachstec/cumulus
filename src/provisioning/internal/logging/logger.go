package logging

import (
	"context"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel/log"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	olog "go.opentelemetry.io/otel/sdk/log"
)

var once sync.Once
var logger *logrus.Logger
var provider log.LoggerProvider
var exporter *otlploghttp.Exporter

// Get initializes and returns a logrus.Logger instance based on the provided config.
func Get(cfg config.Config) *logrus.Logger {
	once.Do(func() {
		ctx := context.Background()
		ex, _ := otlploghttp.New(ctx, otlploghttp.WithEndpoint(cfg.TracingConfig.Endpoint), otlploghttp.WithInsecure())
		exporter = ex
		processor := olog.NewSimpleProcessor(exporter)

		logger = logrus.New()
		buildinfo, _ := debug.ReadBuildInfo()
		provider = olog.NewLoggerProvider(
			olog.WithProcessor(processor),
		)

		// Set the log level (trace level in this case)
		logger.SetLevel(logrus.TraceLevel)

		if cfg.LoggingConfig.Environment == "dev" {
			// In development, log to stdout with a human-friendly formatter.
			logger.SetOutput(os.Stdout)
			logger.SetFormatter(&logrus.TextFormatter{
				TimestampFormat: time.RFC3339Nano,
				FullTimestamp:   true,
			})
			// Add common fields.
			logger = logger.WithFields(logrus.Fields{
				"pid":        os.Getpid(),
				"go_version": buildinfo.GoVersion,
			}).Logger
		} else {
			// In production, log to stderr with a JSON formatter.
			logger.SetOutput(os.Stderr)
			logger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.RFC3339Nano,
			})

			// Attach the OTEL hook for further processing.
			otelHook := otellogrus.NewHook("otel", otellogrus.WithLoggerProvider(provider))
			logger.AddHook(otelHook)
		}
	})

	return logger
}

// LoggingMiddleware returns a gin.HandlerFunc that logs HTTP requests using logrus.
func LoggingMiddleware(cfg config.LoggingConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// Process request.
		ctx.Next()

		// Log the request details.
		Get(config.Config{LoggingConfig: cfg}).WithContext(ctx).WithFields(logrus.Fields{
			"method":      ctx.Request.Method,
			"url":         ctx.Request.RequestURI,
			"user_agent":  ctx.Request.UserAgent(),
			"status_code": ctx.Writer.Status(),
			"elapsed_ms":  time.Since(start).Milliseconds(),
		}).Info("received request")
	}
}
