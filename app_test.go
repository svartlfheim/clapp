package clapp

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrHandleError error = errors.New("some fake error")

type DummyExecutor struct {
	err error
}

func (e *DummyExecutor) Run(c Command, ctx context.Context, cfg interface{}) error {
	return e.err
}

func TestRun(t *testing.T) {
	tests := []struct {
		name            string
		app             App
		exec            Executor
		expectedErr     error
		expectedErrType error
	}{
		{
			name: "errors if cfg is not pointer",
			app: App{
				Config: testConf{},
			},
			exec:        &DummyExecutor{},
			expectedErr: ErrConfigMustBeAPointer,
		},
		{
			name: "errors if config must exist and is not found",
			app: App{
				Config:          &testConf{},
				ConfigPath:      "/tmp/missing/config.yaml",
				ConfigMustExist: true,
				Fs:              buildMockFs(),
				RootCommand: Command{
					Name: "testing",
				},
			},
			exec:        &DummyExecutor{},
			expectedErr: ErrConfigNotFound,
		},
		{
			name: "errors if the config could not be loaded (invalid yaml)",
			app: App{
				Config:          &testConf{},
				ConfigPath:      invalidConfigPath,
				ConfigMustExist: true,
				Fs:              buildMockFs(),
				RootCommand: Command{
					Name: "testing",
				},
			},
			exec:            &DummyExecutor{},
			expectedErrType: ErrUnmarshallingYAML{},
		},
		{
			name: "error is returned from envconfig override",
			app: App{
				Config: &struct {
					Field string `envconfig:"THIS_DONT_EXIST_MATE" required:"true"`
				}{
					Field: "blah",
				},
				ConfigPath:      validConfigPath,
				ConfigMustExist: false,
				Fs:              buildMockFs(),
				RootCommand: Command{
					Name: "testing",
				},
			},
			exec:            &DummyExecutor{},
			expectedErrType: ErrOverridingConfigWithEnvFailed{},
		},
		{
			name: "error is returned from handler",
			app: App{
				Config:          &testConf{},
				ConfigPath:      validConfigPath,
				ConfigMustExist: false,
				Fs:              buildMockFs(),
				RootCommand: Command{
					Name: "testing",
				},
			},
			exec: &DummyExecutor{
				err: ErrHandleError,
			},
			expectedErr: ErrHandleError,
		},
		{
			name: "successful run",
			app: App{
				Config:          &testConf{},
				ConfigPath:      validConfigPath,
				ConfigMustExist: false,
				Fs:              buildMockFs(),
				RootCommand: Command{
					Name: "testing",
				},
			},
			exec: &DummyExecutor{
				err: nil,
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			err := Run(test.app, test.exec)

			if test.expectedErrType != nil {
				assert.IsType(tt, test.expectedErrType, err)
				return
			}

			assert.Equal(tt, test.expectedErr, err)
		})
	}
}
