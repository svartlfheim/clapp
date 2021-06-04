package clapp

import (
	"context"
	"errors"
	"os"
	"reflect"

	"github.com/rs/zerolog"
)

var ErrInvalidLogFormat error = errors.New("log format must be one of: console, json")

type LogManager struct {
	logger zerolog.Logger
}

func (lm *LogManager) ChangeLevel(l int) {
	lm.logger = lm.logger.Level(zerolog.Level(l))
}

func (lm *LogManager) ChangeOutput(f string) error {
	switch f {
	case "console":
		lm.logger = lm.logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "json":
		lm.logger = lm.logger.Output(os.Stderr)
	default:
		return ErrInvalidLogFormat
	}

	return nil
}

func updateLoggerFormatPreRun(ctx context.Context) error {
	cfg := ConfigFromContext(ctx)

	c := reflect.ValueOf(cfg)

	// We should always be using a pointer here
	if c.Kind() != reflect.Ptr {
		return ErrConfigMustBeAPointer
	}

	// Get underlying struct value that we're pointing to
	c = c.Elem()
	if c.Kind() != reflect.Struct {
		return ErrConfigMustPointToAStruct
	}

	v := c.FieldByName("LogFormat")

	if v.Kind() == reflect.Invalid {
		// The field didn't exist so we do nothing
		return nil
	}

	if v.Kind() != reflect.String {
		return ErrLogFormatMustBeString
	}

	lm := LogManagerFromContext(ctx)
	return lm.ChangeOutput(v.Interface().(string))
}

func updateLoggerLevelPreRun(ctx context.Context) error {
	cfg := ConfigFromContext(ctx)

	c := reflect.ValueOf(cfg)

	// We should always be using a pointer here
	if c.Kind() != reflect.Ptr {
		return ErrConfigMustBeAPointer
	}

	// Get underlying struct value that we're pointing to
	c = c.Elem()
	if c.Kind() != reflect.Struct {
		return ErrConfigMustPointToAStruct
	}

	v := c.FieldByName("LogLevel")

	if v.Kind() == reflect.Invalid {
		// The field didn't exist so we do nothing
		return nil
	}

	if v.Kind() != reflect.Int {
		return ErrLogLevelMustBeInt
	}

	lm := LogManagerFromContext(ctx)
	lm.ChangeLevel(v.Interface().(int))

	return nil
}

func UpdateLoggerConfigPreRun(ctx context.Context) error {
	if err := updateLoggerLevelPreRun(ctx); err != nil {
		return err
	}

	if err := updateLoggerFormatPreRun(ctx); err != nil {
		return err
	}

	return nil
}
