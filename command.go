package clapp

import (
	"context"

	"github.com/spf13/cobra"
)

type ValueType string

const StringFlag ValueType = "string"
const StringSliceFlag ValueType = "stringslice"
const IntFlag ValueType = "int"
const IntSliceFlag ValueType = "intslice"
const BoolFlag ValueType = "bool"

type HandlerFunc func(*cobra.Command, []string) error

type Descriptions struct {
	Long  string
	Short string
}

type Flag struct {
	Name        string
	Short       string
	Description string
	ValueRef    interface{}
	Type        ValueType
	Required    bool
}

type Command struct {
	Name                string
	Descriptions        Descriptions
	LocalFlags          []Flag
	PersistentFlags     []Flag
	Handle              HandlerFunc
	CustomConfiguration func(*cobra.Command)
	Children            []Command
}

type Executor interface {
	Run(c Command, ctx context.Context, cfg interface{}) error
}

type builder interface {
	Build(cmd Command, cfg interface{}) (interface{}, error)
}

type builderCallback func() builder
