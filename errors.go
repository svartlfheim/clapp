package clapp

import "fmt"

type ErrIncorrectValueRefForFlag struct {
	expectedType string
}

func (e ErrIncorrectValueRefForFlag) Error() string {
	return fmt.Sprintf("expected ValueRef to be *%s type", e.expectedType)
}

type ErrIncorrectInitialValueForFlag struct {
	expectedType string
}

func (e ErrIncorrectInitialValueForFlag) Error() string {
	return fmt.Sprintf("expected InitialValue to be %s type", e.expectedType)
}

type ErrFlagTypeNotImplemented struct {
	t string
}

func (e ErrFlagTypeNotImplemented) Error() string {
	return fmt.Sprintf("flag type %s not implemented", e.t)
}

type ErrCouldNotBuildRequiredCommandimplementation struct {
	requiredImplementation string
}

func (e ErrCouldNotBuildRequiredCommandimplementation) Error() string {
	return fmt.Sprintf("command was not a %s implementation", e.requiredImplementation)
}

type ErrCommandNotCastable struct {
	CastingTo string
}

func (e ErrCommandNotCastable) Error() string {
	return fmt.Sprintf("could not cast commadn to type %s", e.CastingTo)
}

type ErrOverridingConfigWithEnvFailed struct {
	wrapped error
}

func (e ErrOverridingConfigWithEnvFailed) Error() string {
	return fmt.Sprintf("failed to override config with env vars: %s", e.wrapped.Error())
}

type ErrUnmarshallingYAML struct {
	wrapped error
}

func (e ErrUnmarshallingYAML) Error() string {
	return fmt.Sprintf("could not unmarshal config: %s", e.wrapped.Error())
}

type ErrReadingFile struct {
	wrapped error
}

func (e ErrReadingFile) Error() string {
	return fmt.Sprintf("could not read config: %s", e.wrapped.Error())
}
