package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/svartlfheim/clapp"
)

func configPath() string {
	if env, found := os.LookupEnv("EXAMPLE_CONFIG_PATH"); found {
		return env
	}

	return ""
}


type myconfig struct {
	GlobalVar string `envconfig:"GLOBAL_VAR" yaml:"global"`

	MyConfigVar string `envconfig:"MY_CONFIG_VAR" yaml:"my_config_var"`
}

var appConf myconfig = myconfig{
	// default value, if no config, env or flag override provided
	GlobalVar: "global",

	// default value, if no config, env or flag override provided
	MyConfigVar: "ping", 
}

var anotherVar string = "anothervar-default"
var required int

var app clapp.App = clapp.App{
	// The value must be a pointer to a struct (of any type)
	Config: &appConf,

	// The path to load the config file from
	ConfigPath: configPath(), 
	
	// An error will be returned if the path defined is not found and this is true
	// When this is false, the file will be used if found, but simply ignored if not
	ConfigMustExist: false,

	// If you wish to have a specific context that will be decorated, pass it here
	// if this is nil, a new context will be generated
	InitialContext: context.Background(),

	// This allows you to create an FS to be used by the app
	// if this is nil a new fs will be created for the host filesystem e.g. afero.NewOsFs()
	Fs: afero.NewOsFs(),

	// Create you default logger here
	// If certain variables are found with the Config struct, the level and format will be changed later
	// See: logger.go
	Logger: zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel),

	// Define the root command of the app
	// This is an abstraction on top of the underlying library (typically cobra, but may be extended)
	RootCommand: clapp.Command{

		// The name of the application
		// This value will also be used by envconfig for the prefix
		// envconfig tags in the Config will be prefix with EXAMPLE_ in this case
		Name: "example",

		// This text is shown when using -h or running help
		// This behaviour is inherited from cobra, but can be overridden
		Descriptions: clapp.Descriptions{
			Short: "An example clapp command",
			Long: `This command has sub commands:

child:
	shows values being used for configuration
	child has a subcommand
	
	grandchild:
		shows all configuration
`,
		},

		// See: https://github.com/spf13/cobra#working-with-flags
		// Local flags are only applicable on the current command
		LocalFlags: []clapp.Flag{
			{
				// The name of the option on the cli e.g. --my-config-var=blah
				Name: "my-config-var",
	
				// The short name on the cli e.g. -m blah
				Short: "m",

				// Will be shown in help text
				Description: "This flag sets the MyConfigVar value.",

				// A pointer to a variable to store the flags value in
				// Here we point it to the property of the struct that should be overridden
				ValueRef: &appConf.MyConfigVar,

				// The type of value this flag should be
				// See the ValueType constants in ./command.go for available types
				Type: clapp.StringFlag,

				// Whether this flag needs to be supplied
				// An error will be shown if a required flag is not supplied
				Required: false,
			},
			{
				Name: "another-var",
	
				// An empty value here means no shorthand is available
				Short: "",

				Description: "This flag does not change anything in config.",

				// This will not affect the config, but the value will be stored in this variable
				ValueRef: &anotherVar,

				Type: clapp.StringFlag,

				Required: false,
			},
		},

		// See: https://github.com/spf13/cobra#working-with-flags
		// Persistent flags can be passed to any command that is a child of this one
		PersistentFlags: []clapp.Flag{
			{
				Name: "global-var",

				Description: "A value that can be set globally",

				ValueRef: &appConf.GlobalVar,

				Type: clapp.StringFlag,
			},
		},

		// The function that is actually called when this command is run.
		Handle: func(cmd *cobra.Command, args []string) error {
			// Now we get the config from the context
			// Again the config could be of any type, so we need to typeassert here
			cfg := clapp.ConfigFromContext(cmd.Context()).(*myconfig)

			// Print out the values of the config
			fmt.Printf("GlobalVar is:       %s\n", cfg.GlobalVar)
			fmt.Printf("MyConfigVar var is: %s\n", cfg.MyConfigVar)
			fmt.Printf("AnotherVar var is:  %s\n", anotherVar)
			fmt.Printf("Command version is: %s\n", cmd.Version)

			return nil
		},

		// This allows you to run a function on the underlying command
		// It provides a way to do more complex operations on the underlying command
		// We don't have to provide every property of the command by allowing this
		CustomConfiguration: func(cmd *cobra.Command) {
			// Just an example of modifying the underlying command
			cmd.Version = "1.1.1"
		},

		Children: []clapp.Command{
			{

				// The name of the child command
				Name: "child",
		
				Descriptions: clapp.Descriptions{
					Short: "I'm a child of example",
					Long: `I will just say hello and show the values set for this command.`,
				},
		
				LocalFlags: []clapp.Flag{
					{
						Name: "required",
			
						// An empty value here means no shorthand is available
						Short: "r",
		
						Description: "This bool flag must be supplied.",
		
						// This will not affect the config, but the value will be stored in this variable
						ValueRef: &required,
		
						Type: clapp.IntFlag,
		
						Required: true,
					},
				},
		
				PersistentFlags: []clapp.Flag{
					// This command will inherit the global-var flag from the RootCommand
				},
		
				// The function that is actually called when this command is run.
				Handle: func(cmd *cobra.Command, args []string) error {
					// Now we get the config from the context
					// Again the config could be of any type, so we need to typeassert here
					cfg := clapp.ConfigFromContext(cmd.Context()).(*myconfig)
		
					// Print out the values of the config
					fmt.Printf("GlobalVar is:       %s\n", cfg.GlobalVar)
					fmt.Printf("MyConfigVar var is: %s\n", cfg.MyConfigVar)
					fmt.Printf("AnotherVar var is:  %s\n", anotherVar)
					fmt.Printf("Command version is: %s\n", cmd.Version)
					fmt.Printf("Required value is:  %d\n", required)
		
					return nil
				},
		
				// This allows you to run a function on the underlying command
				// It provides a way to do more complex operations on the underlying command
				// We don't have to provide every property of the command by allowing this
				CustomConfiguration: func(cmd *cobra.Command) {
					// Just an example of modifying the underlying command
					cmd.Version = "1.3.1"
				},
			},
		},
	},
}

func main() {
	if err := clapp.Run(app, clapp.NewCobraExecutor()); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

/*
Try running the following commands to see what values are in the configuration:

	go run main.go
		shows the default values for everything

	EXAMPLE_GLOBAL_VAR=blah go run main.go
		defaults for everything except global-var which has a value of blah (overridden by env var)

	EXAMPLE_GLOBAL_VAR=blah go run main.go --global-var meh
		defaults for everything except global-var which has a value of meh (the flag overrides the env var)

	EXAMPLE_CONFIG_PATH=./config1.yaml go run main.go
		global-var and my-config-var show the values defined in the config1.yaml file based on the yaml tags in the myconfig struct

	EXAMPLE_CONFIG_PATH=./config1.yaml go run main.go --my-config-var meh
		my-config-var show the value from the flag
		global-var shows the value from config1.yaml file based on the yaml tags in the myconfig struct

Have a look in example.yaml to see how clapp loads a file by default based on the name of the root command.

Try as many combinations of flags and configs as you like.
*/