# clapp

Clapp was built to make it quicker to bootstrap a golang application that runs on command line. It is very opinionated, using specific packages for certain features which may or may not fit needs of different individuals. 

---

## Why?

The main drivers for building this were the repetition of adding the following components to command line applications:

- logger
- filesystem abstraction
- configuration loader

### Logger

It should be pretty obvious why a logger is a common requirement for any app, not just a command line app. In this case I wanted a simple way to pass a logger around the application. 

The requirements for the logger were:
- that it allowed additional contextual information to be added to log messages
- allowed logging at different levels which could be toggled easily for debugging purposes
- a simple interface
- efficient performance

For this reasons the choice was [zerolog](https://github.com/rs/zerolog).

As mentioned previously this is opinionated and there are many options that would satisfy the above criteria (e.g. [zaplog](https://github.com/uber-go/zap), [logrus](https://github.com/sirupsen/logrus)). 

The final choice for zerolog came down to personal preference.

### Filesystem abstraction

Go provides some pretty simple ways to interact with the local filesystem, so why use an abstration on top of it? I chose to add an abstraction layer in order to allow us to switch out the underlying implementation, with little effort, if required; currently we utilise this during automated testing.

[Afero](https://github.com/spf13/afero) provides a nice abstraction for working with the filesystem, coming with an in memory driver out of the box.

### Configuration loader

When building any tool, portability makes life alot easier. Depending on where/how you plan to use the CLI tool, there are various methods of injecting configuration which are more convenient. When running in a container, typically environment variables are simplest. But when running as a system daemon it may be more convenient to load configuration from a file somewhere on the host. When trying to debug a tool, it can be more convenient to pass values via command line flags.

Bearing this in mind I typically end up creating a three tier configuration system, firstly loading a configuration file, then overriding values from the file with values from the environment, then finally overriding these with any flags passed to the command. 

Thats what this provides, albeit still with some effort, more simply than rewriting this every time. A struct can be provided to store the configuration, then using the [envconfig](https://github.com/kelseyhightower/envconfig) library, any values in the struct will be overridden (provided the `envconfig` tags exist), then pointers to the values from the struct can be passed as flags for each command. 

They will be loaded using the following order of precedence (top is highest precedence):

1. flags
1. env vars
1. config file

> **Viper**
>
> Considering we're using a couple of tools from spf13 already ([cobra](https://github.com/spf13/cobra), and [afero](https://github.com/spf13/afero)), you may be wondering why not use [viper](https://github.com/spf13/viper). I initially planned to use viper, but came across issues when loading arrays from yaml. It would load the yaml array `[1, 2, 3]` as a string with the value of `[1 2 3]`. This proved to be an issue with the yaml v2 library, so I opted to load the config file manually using yaml v3, then override with envconfig.

---

## How?

How to use this library is best explained with an example, although it is very similar to [cobra](https://github.com/spf13/cobra) itself, I have built a layer of abstraction on top of this.

The easiest way to show what this package does is with an example. You can pull down this repository and run `make runtime` (you will need docker available on your host), and run some of the examples provided in `./example/main.go`.

