package logging

import (
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once

var log zerolog.Logger

func Get(cfg config.Config) zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel := zerolog.TraceLevel

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339Nano,
		}

		buildinfo, _ := debug.ReadBuildInfo()

		if cfg.LoggingConfig.Environment == "dev" {
			log = zerolog.New(output).
				Level(logLevel).
				With().
				Timestamp().
				Caller().
				Int("pid", os.Getpid()).
				Str("go_version", buildinfo.GoVersion).
				Logger()
		} else {
			log = zerolog.New(zerolog.MultiLevelWriter(os.Stderr)).
				Level(logLevel).
				With().
				Timestamp().
				Str("go_version", buildinfo.GoVersion).
				Logger()

			lokiHook := &LokiHook{
				Endpoint: cfg.TracingConfig.Endpoint,
			}

			log = log.Hook(lokiHook)

			log = log.Hook()
		}
	})

	return log
}

func LoggingMiddleware(cfg config.LoggingConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		defer func() {
			log.
				Info().
				Str("method", ctx.Request.Method).
				Str("url", ctx.Request.RequestURI).
				Str("user_agent", ctx.Request.UserAgent()).
				Int("status_code", ctx.Writer.Status()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("received request")
		}()

		ctx.Next()
	}
}
