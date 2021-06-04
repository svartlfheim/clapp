package clapp

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogManager_ChangeLevel(t *testing.T) {
	lm := LogManager{
		logger: zerolog.Logger{},
	}

	// We didn't set it so it should default to 0 (DebugLevel)
	assert.Equal(t, lm.logger.GetLevel(), zerolog.DebugLevel)

	lm.ChangeLevel(int(zerolog.ErrorLevel))

	assert.Equal(t, lm.logger.GetLevel(), zerolog.ErrorLevel)
}

func TestLogManager_ChangeOutput(t *testing.T) {
	lm := LogManager{
		logger: zerolog.Logger{},
	}

	assert.Equal(t, ErrInvalidLogFormat, lm.ChangeOutput("blah"))
	assert.Nil(t, lm.ChangeOutput("console"))
	assert.Nil(t, lm.ChangeOutput("json"))
}

func TestUpdateLoggerLevelPreRun(t *testing.T) {
	tests := []struct {
		name        string
		ctx         contextStub
		expectedErr error
	}{
		{
			name: "error is returned when config is not a pointer",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: struct {
						Field1 string
					}{
						Field1: "blah",
					},
				},
			},
			expectedErr: ErrConfigMustBeAPointer,
		},
		{
			name: "error is returned when config does not point to struct",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &[]string{
						"meh",
					},
				},
			},
			expectedErr: ErrConfigMustPointToAStruct,
		},
		{
			name: "no errors if LogLevel field is not in config",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						Field1 string
					}{
						Field1: "blah",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "error is returned if LogLevel is not int",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						LogLevel string
					}{
						LogLevel: "blah",
					},
				},
			},
			expectedErr: ErrLogLevelMustBeInt,
		},
		{
			name: "log level is changed correctly",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						LogLevel int
					}{
						LogLevel: int(zerolog.DebugLevel),
					},
					LogManagerContextKey: &LogManager{
						logger: zerolog.New(new(bytes.Buffer)).Level(zerolog.ErrorLevel),
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			err := updateLoggerLevelPreRun(test.ctx)

			assert.Equal(tt, test.expectedErr, err)

			if lm, ok := test.ctx.Vals[LogManagerContextKey].(*LogManager); ok {
				assert.Equal(t, lm.logger.GetLevel(), zerolog.DebugLevel)
			}
		})
	}
}

func TestUpdateLoggerFormatPreRun(t *testing.T) {
	tests := []struct {
		name        string
		ctx         contextStub
		expectedErr error
	}{
		{
			name: "error is returned when config is not a pointer",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: struct {
						Field1 string
					}{
						Field1: "blah",
					},
				},
			},
			expectedErr: ErrConfigMustBeAPointer,
		},
		{
			name: "error is returned when config does not point to struct",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &[]string{
						"meh",
					},
				},
			},
			expectedErr: ErrConfigMustPointToAStruct,
		},
		{
			name: "no errors if LogFormat field is not in config",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						Field1 string
					}{
						Field1: "blah",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "error is returned if LogFormat is not int",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						LogFormat int
					}{
						LogFormat: 4,
					},
				},
			},
			expectedErr: ErrLogFormatMustBeString,
		},
		{
			name: "log format is changed without error for valid format",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						LogFormat string
					}{
						LogFormat: "json",
					},
					LogManagerContextKey: &LogManager{
						logger: zerolog.New(new(bytes.Buffer)),
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "error is returned for invalid format",
			ctx: contextStub{
				Vals: map[interface{}]interface{}{
					ConfigContextKey: &struct {
						LogFormat string
					}{
						LogFormat: "blah",
					},
					LogManagerContextKey: &LogManager{
						logger: zerolog.New(new(bytes.Buffer)),
					},
				},
			},
			expectedErr: ErrInvalidLogFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			err := updateLoggerFormatPreRun(test.ctx)

			assert.Equal(tt, test.expectedErr, err)
		})
	}
}
