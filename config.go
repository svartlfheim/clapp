package clapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

var ErrCannotUseNonPointerValue error = errors.New("cannot use non-pointer value")
var ErrConfigNotFound error = errors.New("config file does not exist")

type Config struct {
	filePath        string
	appName         string
	fs              afero.Fs
	configMustExist bool
}

type configOpt func(c *Config)

func FilePathOpt(path string) configOpt {
	return func(c *Config) {
		c.filePath = path
	}
}

func FileMustExistOpt() configOpt {
	return func(c *Config) {
		c.configMustExist = true
	}
}

func (c *Config) load(into interface{}) error {
	if exists, err := afero.Exists(c.fs, c.filePath); !exists || err != nil {
		return ErrConfigNotFound
	}

	cfgBytes, err := afero.ReadFile(c.fs, c.filePath)

	if err != nil {
		return ErrReadingFile{
			wrapped: err,
		}
	}

	if err := yaml.Unmarshal(cfgBytes, into); err != nil {
		return ErrUnmarshallingYAML{
			wrapped: err,
		}
	}

	return nil
}

func (c *Config) OverrideWithEnvVars(cfg interface{}) error {
	rval := reflect.ValueOf(cfg)

	if rval.Kind() != reflect.Ptr {
		return ErrCannotUseNonPointerValue
	}

	return envconfig.Process(c.appName, cfg)
}

func newConfigManager(ctx context.Context, cfg interface{}, appName string, filePath string, mustExist bool) (*Config, error) {
	c := &Config{
		appName:         appName,
		fs:              FsFromContext(ctx),
		configMustExist: false,
		filePath:        fmt.Sprintf("./%s.yaml", appName),
	}

	if filePath != "" {
		FilePathOpt(filePath)(c)
	}

	if mustExist {
		FileMustExistOpt()(c)
	}

	err := c.load(cfg)

	if !c.configMustExist && err == ErrConfigNotFound {
		return c, nil
	}

	return c, err
}
