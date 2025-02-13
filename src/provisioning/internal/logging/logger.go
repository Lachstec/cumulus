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

func Get(cfg config.LoggingConfig) zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel := zerolog.TraceLevel

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339Nano,
		}

		var gitRevision string

		buildinfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildinfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		if cfg.Environment == "dev" {
			log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}).
				Level(logLevel).
				With().
				Timestamp().
				Caller().
				Int("pid", os.Getpid()).
				Str("go_version", buildinfo.GoVersion).
				Logger()
		} else {
			log = zerolog.New(output).
				Level(logLevel).
				With().
				Timestamp().
				Str("git_revision", gitRevision).
				Str("go_version", buildinfo.GoVersion).
				Logger()
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
				Dur("elapsed_ms", time.Since(start)).
				Msg("received request")
		}()

		ctx.Next()
	}
}
