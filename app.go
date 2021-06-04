package clapp

import (
	"context"
	"errors"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

var ErrConfigMustBeAPointer error = errors.New("cannot use non-pointer value for config")
var ErrConfigMustPointToAStruct error = errors.New("cannot point to non-struct value for config")
var ErrLogLevelMustBeInt error = errors.New("log level type in config struct must be int")
var ErrLogFormatMustBeString error = errors.New("log format type in config struct must be string")

type App struct {
	Config          interface{}
	ConfigPath      string
	ConfigMustExist bool
	InitialContext  context.Context
	Fs              afero.Fs
	Logger          zerolog.Logger
	RootCommand     Command
}

func Run(a App, e Executor) error {
	rval := reflect.ValueOf(a.Config)

	if rval.Kind() != reflect.Ptr {
		return ErrConfigMustBeAPointer
	}

	initCtx := a.InitialContext

	if initCtx == nil {
		initCtx = context.Background()
	}

	ctx := buildContext(initCtx, a.Fs, a.Logger, a.Config)
	cfgManager, err := newConfigManager(ctx, a.Config, a.RootCommand.Name, a.ConfigPath, a.ConfigMustExist)

	if err != nil {
		return err
	}

	if err := cfgManager.OverrideWithEnvVars(a.Config); err != nil {
		return ErrOverridingConfigWithEnvFailed{
			wrapped: err,
		}
	}

	return e.Run(a.RootCommand, ctx, a.Config)
}
