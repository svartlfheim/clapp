package clapp

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func countFlags(s *pflag.FlagSet) int {
	i := 0
	s.VisitAll(func(_ *pflag.Flag) {
		i++
	})

	return i
}

func countRequiredFlags(s *pflag.FlagSet) int {
	i := 0
	s.VisitAll(func(f *pflag.Flag) {
		if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok {
			i++
		}
	})

	return i
}

func TestCobraBuilder_SetName(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	b.setName("blah")
	assert.Equal(t, "blah", b._cmd.Use)
}

func TestCobraBuilder_SetDescriptions(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	b.setDescriptions("short", "long")
	assert.Equal(t, "short", b._cmd.Short)
	assert.Equal(t, "long", b._cmd.Long)
}

func TestCobraBuilder_handleStringFlag_failsIfValueRefIsNotPointer(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	err := b.handleStringFlag(b._cmd.Flags(), Flag{
		ValueRef: "not a pointer",
	})

	assert.Equal(t, ErrIncorrectValueRefForFlag{
		expectedType: "string",
	}, err)
}

func TestCobraBuilder_handleStringSliceFlag_failsIfValueRefIsNotPointer(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	err := b.handleStringSliceFlag(b._cmd.Flags(), Flag{
		ValueRef: []string{"not a pointer"},
	})

	assert.Equal(t, ErrIncorrectValueRefForFlag{
		expectedType: "[]string",
	}, err)
}

func TestCobraBuilder_handleIntFlag_failsIfValueRefIsNotPointer(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	err := b.handleIntFlag(b._cmd.Flags(), Flag{
		ValueRef: 7346,
	})

	assert.Equal(t, ErrIncorrectValueRefForFlag{
		expectedType: "int",
	}, err)
}

func TestCobraBuilder_handleIntSliceFlag_failsIfValueRefIsNotPointer(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	err := b.handleIntSliceFlag(b._cmd.Flags(), Flag{
		ValueRef: []int{3},
	})

	assert.Equal(t, ErrIncorrectValueRefForFlag{
		expectedType: "[]int",
	}, err)
}

func TestCobraBuilder_handleBoolFlag_failsIfValueRefIsNotPointer(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	err := b.handleBoolFlag(b._cmd.Flags(), Flag{
		ValueRef: false,
	})

	assert.Equal(t, ErrIncorrectValueRefForFlag{
		expectedType: "bool",
	}, err)
}

func TestCobraBuilder_handleStringFlag_addsFlagWithoutAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.Str("meh")

	err := b.handleStringFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "meh", f.DefValue)
	assert.Equal(t, "meh", f.Value.String())
	assert.Equal(t, "string", f.Value.Type())
}

func TestCobraBuilder_handleStringFlag_addsFlagWithAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.Str("meh")

	err := b.handleStringFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Short:       "b",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, f.Shorthand, "b")
	assert.Equal(t, "meh", f.DefValue)
	assert.Equal(t, "meh", f.Value.String())
	assert.Equal(t, "string", f.Value.Type())
}

func TestCobraBuilder_handleIntFlag_addsFlagWithoutAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.Int(4637)

	err := b.handleIntFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "4637", f.DefValue)
	assert.Equal(t, "4637", f.Value.String())
	assert.Equal(t, "int", f.Value.Type())
}

func TestCobraBuilder_handleIntFlag_addsFlagWithAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.Int(4637)

	err := b.handleIntFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Short:       "b",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, f.Shorthand, "b")
	assert.Equal(t, "4637", f.Value.String())
	assert.Equal(t, "4637", f.DefValue)
	assert.Equal(t, "int", f.Value.Type())
}

func TestCobraBuilder_handleStringSliceFlag_addsFlagWithoutAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.StrSlice([]string{"meh"})

	err := b.handleStringSliceFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "[meh]", f.DefValue)
	assert.Equal(t, "[meh]", f.Value.String())
	assert.Equal(t, "stringSlice", f.Value.Type())
}

func TestCobraBuilder_handleStringSliceFlag_addsFlagWithAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.StrSlice([]string{"meh"})

	err := b.handleStringSliceFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Short:       "b",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, f.Shorthand, "b")
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "[meh]", f.DefValue)
	assert.Equal(t, "[meh]", f.Value.String())
	assert.Equal(t, "stringSlice", f.Value.Type())
}

func TestCobraBuilder_handleIntSliceFlag_addsFlagWithoutAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.IntSlice([]int{4635})

	err := b.handleIntSliceFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "[4635]", f.DefValue)
	assert.Equal(t, "[4635]", f.Value.String())
	assert.Equal(t, "intSlice", f.Value.Type())
}

func TestCobraBuilder_handleIntSliceFlag_addsFlagWithAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.IntSlice([]int{4635})

	err := b.handleIntSliceFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Short:       "b",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, f.Shorthand, "b")
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "[4635]", f.DefValue)
	assert.Equal(t, "[4635]", f.Value.String())
	assert.Equal(t, "intSlice", f.Value.Type())
}

func TestCobraBuilder_handleBoolFlag_addsFlagWithoutAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.Bool(false)

	err := b.handleBoolFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, "false", f.DefValue)
	assert.Equal(t, "false", f.Value.String())
	assert.Equal(t, "bool", f.Value.Type())
}

func TestCobraBuilder_handleBoolFlag_addsFlagWithAlias(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	val := pointTo.Bool(true)

	err := b.handleBoolFlag(b._cmd.Flags(), Flag{
		Name:        "blah",
		Short:       "b",
		Description: "my string desc",
		ValueRef:    val,
	})

	assert.Equal(t, nil, err)
	f := b._cmd.Flags().Lookup("blah")

	assert.NotNil(t, f)

	assert.Equal(t, "blah", f.Name)
	assert.Equal(t, "my string desc", f.Usage)
	assert.Equal(t, f.Shorthand, "b")
	assert.Equal(t, "true", f.Value.String())
	assert.Equal(t, "true", f.DefValue)
	assert.Equal(t, "bool", f.Value.Type())
}

func TestCobraBuilder_handleFlag(t *testing.T) {
	tests := []struct {
		name        string
		set         *pflag.FlagSet
		f           Flag
		expectedErr error
	}{
		{
			name: "adds a string flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.Str("meh"),
				Description: "some desc",
				Type:        StringFlag,
			},
			expectedErr: nil,
		},
		{
			name: "adds an int flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.Int(23486),
				Description: "some desc",
				Type:        IntFlag,
			},
			expectedErr: nil,
		},
		{
			name: "adds a []string flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.StrSlice([]string{"meh"}),
				Description: "some desc",
				Type:        StringSliceFlag,
			},
			expectedErr: nil,
		},
		{
			name: "adds a []int flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.IntSlice([]int{287364}),
				Description: "some desc",
				Type:        IntSliceFlag,
			},
			expectedErr: nil,
		},
		{
			name: "adds a bool flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.Bool(false),
				Description: "some desc",
				Type:        BoolFlag,
			},
			expectedErr: nil,
		},
		{
			name: "returns error for string flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.IntSlice([]int{287364}),
				Description: "some desc",
				Type:        StringFlag,
			},
			expectedErr: ErrIncorrectValueRefForFlag{
				expectedType: "string",
			},
		},
		{
			name: "returns error for int flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.IntSlice([]int{287364}),
				Description: "some desc",
				Type:        IntFlag,
			},
			expectedErr: ErrIncorrectValueRefForFlag{
				expectedType: "int",
			},
		},
		{
			name: "returns error for []string flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.IntSlice([]int{287364}),
				Description: "some desc",
				Type:        StringSliceFlag,
			},
			expectedErr: ErrIncorrectValueRefForFlag{
				expectedType: "[]string",
			},
		},
		{
			name: "returns error for []int flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.Int(287364),
				Description: "some desc",
				Type:        IntSliceFlag,
			},
			expectedErr: ErrIncorrectValueRefForFlag{
				expectedType: "[]int",
			},
		},
		{
			name: "returns error for bool flag",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.Str("meh"),
				Description: "some desc",
				Type:        BoolFlag,
			},
			expectedErr: ErrIncorrectValueRefForFlag{
				expectedType: "bool",
			},
		},
		{
			name: "returns error for unhandled type",
			set:  pflag.NewFlagSet("test", pflag.PanicOnError),
			f: Flag{
				Name:        "blah",
				ValueRef:    pointTo.Int(287364),
				Description: "some desc",
				Type:        ValueType("meh"),
			},
			expectedErr: ErrFlagTypeNotImplemented{
				t: "meh",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			b := &cobraBuilder{
				_cmd: &cobra.Command{},
			}

			err := b.handleFlag(test.set, test.f)

			assert.Equal(tt, test.expectedErr, err)

			if test.expectedErr == nil {
				// Ensure the flag was actually added
				pf := test.set.Lookup(test.f.Name)
				assert.NotNil(tt, pf)
			}
		})
	}
}

var flagTestCases = []struct {
	name                  string
	flags                 []Flag
	expectedErr           error
	expectedFlags         []string
	expectedRequiredFlags []string
}{
	{
		name:        "nothing happens with no flags",
		flags:       []Flag{},
		expectedErr: nil,
	},
	{
		name: "error is returned for first flag",
		flags: []Flag{
			{
				Name:     "blah",
				Type:     ValueType("meh"),
				ValueRef: pointTo.Int(213478),
			},
		},
		expectedErr: ErrFlagTypeNotImplemented{
			t: "meh",
		},
	},
	{
		name: "error is returned for subsequent flag",
		flags: []Flag{
			{
				Name:     "blah",
				Type:     StringFlag,
				ValueRef: pointTo.Str("blah"),
			},
			{
				Name:     "nope",
				Type:     ValueType("meh"),
				ValueRef: pointTo.Str("blah"),
			},
		},
		expectedErr: ErrFlagTypeNotImplemented{
			t: "meh",
		},
	},
	{
		name: "flags are successfully added",
		flags: []Flag{
			{
				Name:     "blah",
				Type:     StringFlag,
				ValueRef: pointTo.Str("blah"),
			},
			{
				Name:     "yup",
				Type:     IntFlag,
				ValueRef: pointTo.Int(2734),
			},
		},
		expectedFlags: []string{
			"blah",
			"yup",
		},
		expectedErr: nil,
	},
	{
		name: "flags are successfully added and 1 is marked as required",
		flags: []Flag{
			{
				Name:     "blah",
				Type:     StringFlag,
				ValueRef: pointTo.Str("blah"),
			},
			{
				Name:     "yup",
				Type:     IntFlag,
				ValueRef: pointTo.Int(2734),
				Required: true,
			},
		},
		expectedFlags: []string{
			"blah",
			"yup",
		},
		expectedRequiredFlags: []string{
			"yup",
		},
		expectedErr: nil,
	},
	{
		name: "flags are successfully added and are marked as required",
		flags: []Flag{
			{
				Name:     "blah",
				Type:     StringFlag,
				ValueRef: pointTo.Str("blah"),
				Required: true,
			},
			{
				Name:     "yup",
				Type:     IntFlag,
				ValueRef: pointTo.Int(2734),
				Required: true,
			},
		},
		expectedFlags: []string{
			"blah",
			"yup",
		},
		expectedRequiredFlags: []string{
			"blah",
			"yup",
		},
		expectedErr: nil,
	},
}

func TestCobraBuilder_addPersistentFlags(t *testing.T) {
	for _, test := range flagTestCases {
		t.Run(test.name, func(tt *testing.T) {
			b := &cobraBuilder{
				_cmd: &cobra.Command{},
			}

			err := b.addPersistentFlags(test.flags...)

			assert.Equal(tt, test.expectedErr, err, "got unexpected error")

			if test.expectedErr != nil {
				return
			}

			fs := b._cmd.PersistentFlags()

			foundFlags := countFlags(fs)
			assert.Equal(tt, foundFlags, len(test.expectedFlags), "Expected %d flag(s), found %d", len(test.expectedFlags), foundFlags)

			for _, n := range test.expectedFlags {
				v := fs.Lookup(n)

				assert.NotNil(tt, v)
			}

			foundRequiredFlags := countRequiredFlags(fs)
			assert.Equal(tt, foundRequiredFlags, len(test.expectedRequiredFlags), "Expected %d required flag(s), found %d", len(test.expectedRequiredFlags), foundRequiredFlags)

			for _, n := range test.expectedRequiredFlags {
				v := fs.Lookup(n)

				assert.NotNil(tt, v)
				// This annotation should exist if it was marked as required
				_, ok := v.Annotations[cobra.BashCompOneRequiredFlag]

				assert.True(tt, ok)
			}
		})
	}
}

func TestCobraBuilder_addLocalFlags(t *testing.T) {
	for _, test := range flagTestCases {
		t.Run(test.name, func(tt *testing.T) {
			b := &cobraBuilder{
				_cmd: &cobra.Command{},
			}

			err := b.addLocalFlags(test.flags...)

			assert.Equal(tt, test.expectedErr, err, "got unexpected error")

			if test.expectedErr != nil {
				return
			}

			fs := b._cmd.Flags()

			foundFlags := countFlags(fs)
			assert.Equal(tt, foundFlags, len(test.expectedFlags), "Expected %d flag(s), found %d", len(test.expectedFlags), foundFlags)

			for _, n := range test.expectedFlags {
				v := fs.Lookup(n)

				assert.NotNil(tt, v)
			}

			foundRequiredFlags := countRequiredFlags(fs)
			assert.Equal(tt, foundRequiredFlags, len(test.expectedRequiredFlags), "Expected %d required flag(s), found %d", len(test.expectedRequiredFlags), foundRequiredFlags)

			for _, n := range test.expectedRequiredFlags {
				v := fs.Lookup(n)

				assert.NotNil(tt, v)
				// This annotation should exist if it was marked as required
				_, ok := v.Annotations[cobra.BashCompOneRequiredFlag]

				assert.True(tt, ok)
			}
		})
	}
}

func TestCobraBuilder_setHandler(t *testing.T) {
	b := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	var expectedErr error = errors.New("we should get this back")
	var h HandlerFunc = func(i interface{}, args []string) error {
		return expectedErr
	}

	b.setHandler(h)

	err := b._cmd.RunE(b._cmd, []string{})

	assert.Equal(t, expectedErr, err)
}

func TestCobraBuilder_addChildCommands(t *testing.T) {
	tests := []struct {
		name             string
		children         []Command
		cfg              interface{}
		builder          *cobraBuilderStub
		expectedErr      error
		expectedChildren []string
	}{
		{
			name:             "nothing happens with empty list",
			children:         []Command{},
			cfg:              testConf{},
			builder:          &cobraBuilderStub{},
			expectedErr:      nil,
			expectedChildren: []string{},
		},
		{
			name: "error is returned when first child build errors",
			children: []Command{
				{
					Name: "first",
				},
			},
			cfg: testConf{},
			builder: &cobraBuilderStub{
				BuildResults: []interface{}{
					"blah",
				},
				BuildErrors: []error{
					errors.New("fake error"),
				},
			},
			expectedErr: errors.New("fake error"),
		},
		{
			name: "error is returned when subsequent child build errors",
			children: []Command{
				{
					Name: "first",
				}, {
					Name: "second",
				},
			},
			cfg: testConf{},
			builder: &cobraBuilderStub{
				BuildResults: []interface{}{
					&cobra.Command{},
					"meh",
				},
				BuildErrors: []error{
					nil,
					errors.New("fake error"),
				},
			},
			expectedErr: errors.New("fake error"),
		},
		{
			name: "error is returned when child typecast errors",
			children: []Command{
				{
					Name: "first",
				},
			},
			builder: &cobraBuilderStub{
				BuildResults: []interface{}{
					"blah",
				},
				BuildErrors: []error{
					nil,
				},
			},
			expectedErr: ErrCouldNotBuildRequiredCommandimplementation{
				requiredImplementation: "cobra",
			},
		},
		{
			name: "command is successfully added to parent",
			children: []Command{
				{
					Name: "first",
				},
			},
			builder: &cobraBuilderStub{
				BuildResults: []interface{}{
					&cobra.Command{
						Use: "first",
					},
				},
				BuildErrors: []error{
					nil,
				},
			},
			expectedErr: nil,
			expectedChildren: []string{
				"first",
			},
		},
		{
			name: "multiple commands are successfully added to parent",
			children: []Command{
				{
					Name: "first",
				},
				{
					Name: "second",
				},
				{
					Name: "third",
				},
			},
			builder: &cobraBuilderStub{
				BuildResults: []interface{}{
					&cobra.Command{
						Use: "first",
					},
					&cobra.Command{
						Use: "second",
					},
					&cobra.Command{
						Use: "third",
					},
				},
				BuildErrors: []error{
					nil,
					nil,
					nil,
				},
			},
			expectedErr: nil,
			expectedChildren: []string{
				"first",
				"second",
				"third",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			bc := func() builder {
				return test.builder
			}

			b := &cobraBuilder{
				_cmd: &cobra.Command{},
			}

			err := b.addChildCommands(bc, test.cfg, test.children...)

			assert.Equal(tt, test.expectedErr, err)

			if test.expectedErr != nil {
				return
			}

			givenNames := []string{}

			for _, ch := range b._cmd.Commands() {
				givenNames = append(givenNames, ch.Use)
			}

			assert.Equal(tt, test.expectedChildren, givenNames)
		})
	}
}

func TestCobraBuidler_customConfigure(t *testing.T) {
	cb := &cobraBuilder{
		_cmd: &cobra.Command{},
	}

	cb.customConfigure(func(i interface{}) {
		cmd, ok := i.(*cobra.Command)

		assert.True(t, ok)

		cmd.Use = "blah"
	})

	assert.Equal(t, "blah", cb._cmd.Use)
}
