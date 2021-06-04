package clapp

import (
	"bytes"
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFSFromContext(t *testing.T) {
	expectedFs := afero.NewMemMapFs()
	ctx := contextWithFs(context.TODO(), expectedFs)

	assert.Equal(t, expectedFs, FsFromContext(ctx))
}

func TestConfigFromContext(t *testing.T) {
	cfg := &struct {
		Field string
	}{
		Field: "blah",
	}

	ctx := contextWithConfig(context.TODO(), cfg)

	assert.Equal(t, cfg, ConfigFromContext(ctx))
}

func TestLoggerFromContext(t *testing.T) {
	b := new(bytes.Buffer)
	l := zerolog.New(b)
	ctx := contextWithLogger(context.TODO(), l)

	assert.IsType(t, &LogManager{}, LogManagerFromContext(ctx))
	assert.Equal(t, l, LogManagerFromContext(ctx).logger)
	assert.Equal(t, l, LoggerFromContext(ctx))
}

func TestBuildContext(t *testing.T) {
	baseCtx := context.TODO()
	fs := afero.NewMemMapFs()
	l := zerolog.Logger{}
	cfg := struct {
		Field1 string
	}{
		Field1: "somevalue",
	}

	ctx := buildContext(baseCtx, fs, l, &cfg)

	assert.Equal(t, fs, FsFromContext(ctx))
	assert.Equal(t, l, LoggerFromContext(ctx))
	assert.Equal(t, &cfg, ConfigFromContext(ctx))
}

func TestBuildContext_FsIsCreatedIfNil(t *testing.T) {
	baseCtx := context.TODO()
	l := zerolog.Logger{}
	cfg := struct {
		Field1 string
	}{
		Field1: "somevalue",
	}

	ctx := buildContext(baseCtx, nil, l, &cfg)

	assert.Equal(t, afero.NewOsFs(), FsFromContext(ctx))
	assert.Equal(t, l, LoggerFromContext(ctx))
	assert.Equal(t, &cfg, ConfigFromContext(ctx))
}
