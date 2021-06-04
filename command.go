package clapp

import (
	"context"
)

type ValueType string

var StringFlag ValueType = "string"
var StringSliceFlag ValueType = "stringslice"
var IntFlag ValueType = "int"
var IntSliceFlag ValueType = "intslice"
var BoolFlag ValueType = "bool"

type HandlerFunc func(interface{}, []string) error

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
	Handle              func(interface{}, []string) error
	CustomConfiguration func(interface{})
	Children            []Command
}

type Executor interface {
	Run(c Command, ctx context.Context, cfg interface{}) error
}

type builder interface {
	Build(cmd Command, cfg interface{}) (interface{}, error)
}

type builderCallback func() builder
