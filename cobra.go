package clapp

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cobraBuilder struct {
	_cmd *cobra.Command
}

type CobraExecutor struct {
	_builder *cobraBuilder
}

func NewCobraExecutor() *CobraExecutor {
	return &CobraExecutor{
		_builder: newCobraBuilder().(*cobraBuilder),
	}
}

func newCobraBuilder() builder {
	return &cobraBuilder{
		_cmd: &cobra.Command{},
	}
}

func (e CobraExecutor) Run(c Command, ctx context.Context, cfg interface{}) error {
	cmd, err := e._builder.Build(c, cfg)

	if err != nil {
		return err
	}

	cobraCmd := cmd.(*cobra.Command)

	return cobraCmd.ExecuteContext(ctx)
}

func (b *cobraBuilder) setName(n string) {
	b._cmd.Use = n
}

func (b *cobraBuilder) setDescriptions(s string, l string) {
	b._cmd.Short = s
	b._cmd.Long = l
}

func (b *cobraBuilder) handleStringFlag(s *pflag.FlagSet, f Flag) error {
	ref, ok := f.ValueRef.(*string)

	if !ok {
		return ErrIncorrectValueRefForFlag{
			expectedType: "string",
		}
	}

	if f.Short != "" {
		s.StringVarP(ref, f.Name, f.Short, *ref, f.Description)
		return nil
	}

	s.StringVar(ref, f.Name, *ref, f.Description)

	return nil
}

func (b *cobraBuilder) handleStringSliceFlag(s *pflag.FlagSet, f Flag) error {
	ref, ok := f.ValueRef.(*[]string)

	if !ok {
		return ErrIncorrectValueRefForFlag{
			expectedType: "[]string",
		}
	}

	if f.Short != "" {
		s.StringSliceVarP(ref, f.Name, f.Short, *ref, f.Description)
		return nil
	}

	s.StringSliceVar(ref, f.Name, *ref, f.Description)

	return nil
}

func (b *cobraBuilder) handleIntFlag(s *pflag.FlagSet, f Flag) error {
	ref, ok := f.ValueRef.(*int)

	if !ok {
		return ErrIncorrectValueRefForFlag{
			expectedType: "int",
		}
	}

	if f.Short != "" {
		s.IntVarP(ref, f.Name, f.Short, *ref, f.Description)
		return nil
	}

	s.IntVar(ref, f.Name, *ref, f.Description)

	return nil
}

func (b *cobraBuilder) handleIntSliceFlag(s *pflag.FlagSet, f Flag) error {
	ref, ok := f.ValueRef.(*[]int)

	if !ok {
		return ErrIncorrectValueRefForFlag{
			expectedType: "[]int",
		}
	}

	if f.Short != "" {
		s.IntSliceVarP(ref, f.Name, f.Short, *ref, f.Description)
		return nil
	}

	s.IntSliceVar(ref, f.Name, *ref, f.Description)

	return nil
}

func (b *cobraBuilder) handleBoolFlag(s *pflag.FlagSet, f Flag) error {
	ref, ok := f.ValueRef.(*bool)

	if !ok {
		return ErrIncorrectValueRefForFlag{
			expectedType: "bool",
		}
	}

	if f.Short != "" {
		s.BoolVarP(ref, f.Name, f.Short, *ref, f.Description)
		return nil
	}

	s.BoolVar(ref, f.Name, *ref, f.Description)

	return nil
}

func (b *cobraBuilder) handleFlag(s *pflag.FlagSet, f Flag) error {
	switch f.Type {
	case StringFlag:
		return b.handleStringFlag(s, f)
	case StringSliceFlag:
		return b.handleStringSliceFlag(s, f)
	case IntFlag:
		return b.handleIntFlag(s, f)
	case IntSliceFlag:
		return b.handleIntSliceFlag(s, f)
	case BoolFlag:
		return b.handleBoolFlag(s, f)
	}

	return ErrFlagTypeNotImplemented{
		t: string(f.Type),
	}
}

func (b *cobraBuilder) addPersistentFlags(flags ...Flag) error {
	for _, f := range flags {
		if err := b.handleFlag(b._cmd.PersistentFlags(), f); err != nil {
			return err
		}

		if f.Required {
			err := b._cmd.MarkPersistentFlagRequired(f.Name)

			// untestable:
			// the only way this could return an error is if the
			// flag doesn't exist in the set. This isn't possible here as the flag
			// is added above or we bailed out the function. We'll leave this
			// here in case something very weird/unpredictable happens
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *cobraBuilder) addLocalFlags(flags ...Flag) error {
	for _, f := range flags {
		if err := b.handleFlag(b._cmd.Flags(), f); err != nil {
			return err
		}

		if f.Required {
			err := b._cmd.MarkFlagRequired(f.Name)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *cobraBuilder) setHandler(h HandlerFunc) {
	b._cmd.RunE = func(c *cobra.Command, args []string) error {
		if h != nil {
			return h(c, args)
		}

		return c.Help()
	}
}

func (b *cobraBuilder) addChildCommands(bC builderCallback, cfg interface{}, children ...Command) error {
	for _, c := range children {
		cmd, err := bC().Build(c, cfg)

		if err != nil {
			return err
		}

		cobraCmd, ok := cmd.(*cobra.Command)

		if !ok {
			return ErrCouldNotBuildRequiredCommandimplementation{
				requiredImplementation: "cobra",
			}
		}

		b._cmd.AddCommand(cobraCmd)
	}

	return nil
}

func (b *cobraBuilder) customConfigure(f func(interface{})) {
	f(b._cmd)
}

func (b *cobraBuilder) Build(cmd Command, cfg interface{}) (interface{}, error) {
	b.setName(cmd.Name)
	b.setDescriptions(cmd.Descriptions.Short, cmd.Descriptions.Long)

	err := b.addPersistentFlags(cmd.PersistentFlags...)

	if err != nil {
		return nil, err
	}

	err = b.addLocalFlags(cmd.LocalFlags...)

	if err != nil {
		return nil, err
	}

	b.setHandler(cmd.Handle)
	err = b.addChildCommands(newCobraBuilder, cfg, cmd.Children...)

	if err != nil {
		return nil, err
	}

	if cmd.CustomConfiguration != nil {
		b.customConfigure(cmd.CustomConfiguration)
	}

	return b._cmd, nil
}
