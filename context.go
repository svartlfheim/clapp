package clapp

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type ContextKey string

const FsContextKey ContextKey = "FILESYSTEM"
const ConfigContextKey ContextKey = "APP_CONFIG"
const ConfigManagerContextKey ContextKey = "CONFIG_MANAGER"
const LogManagerContextKey ContextKey = "LOG_MANAGER"

func FsFromContext(ctx context.Context) afero.Fs {
	return ctx.Value(FsContextKey).(afero.Fs)
}

func ConfigFromContext(ctx context.Context) interface{} {
	return ctx.Value(ConfigContextKey)
}

func LogManagerFromContext(ctx context.Context) *LogManager {
	return ctx.Value(LogManagerContextKey).(*LogManager)
}

func LoggerFromContext(ctx context.Context) zerolog.Logger {
	return LogManagerFromContext(ctx).logger
}

func contextWithFs(ctx context.Context, fs afero.Fs) context.Context {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	return context.WithValue(
		ctx,
		FsContextKey,
		fs,
	)
}

func contextWithConfig(ctx context.Context, cfg interface{}) context.Context {
	return context.WithValue(
		ctx,
		ConfigContextKey,
		cfg,
	)
}

func contextWithLogger(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(
		ctx,
		LogManagerContextKey,
		&LogManager{
			logger: l,
		},
	)
}

func buildContext(ctx context.Context, fs afero.Fs, l zerolog.Logger, cfg interface{}) (c context.Context) {
	c = contextWithFs(ctx, fs)
	c = contextWithConfig(c, cfg)
	c = contextWithLogger(c, l)

	return
}
