package clapp

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
)

var validYAML string = `# comment just to make it more readable on next line
should-enable-this: "awesome-feature"
should-do-that: "work-properly"
external-endpoint:
  protocol: "https"
  domain: somefakesite.com
  port: 6738
list-of-things:
  - thing1
  - thing2
  - thing3
`

var aintValidYAML = `{CW £863
	got tabs in, yaml no likey
			
					
	make sure IDE's don't reformat this...
	
		
		tabs everywhere!!!
`

var configDir string = "/config-test/meh/"
var validConfigPath string = fmt.Sprintf("%s/%s.yaml", configDir, "valid")
var invalidConfigPath string = fmt.Sprintf("%s/%s.yaml", configDir, "invalid")

type contextStub struct {
	DeadlineVal    time.Time
	DeadlineIsOkay bool
	Vals           map[interface{}]interface{}
	ErrVal         error
}

func (c contextStub) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (c contextStub) Deadline() (time.Time, bool) {
	return c.DeadlineVal, c.DeadlineIsOkay
}

func (c contextStub) Err() error {
	return c.ErrVal
}

func (c contextStub) Value(k interface{}) interface{} {
	return c.Vals[k]
}

type testExtEndpoint struct {
	Protocol string `yaml:"protocol" envconfig:"PROTOCOL"`
	Domain   string `yaml:"domain" envconfig:"DOMAIN"`
	Port     int    `yaml:"port" envconfig:"port"`
}

type testConf struct {
	ShouldEnableThis string          `yaml:"should-enable-this" envconfig:"SHOULD_ENABLE_THIS"`
	ShouldDoThat     string          `yaml:"should-do-that" envconfig:"SHOULD_DO_THAT"`
	ExternalEndpoint testExtEndpoint `yaml:"external-endpoint" envconfig:"EXTERNAL_ENDPOINT"`
	ListOfThings     []string        `yaml:"list-of-things" envconfig:"LIST_OF_THINGS"`
}

func buildMockFs() afero.Fs {
	fs := afero.NewMemMapFs()

	if err := fs.Mkdir("/config-test", os.ModeDir); err != nil {
		panic(err)
	}

	if err := fs.Mkdir("/config-test/meh", os.ModeDir); err != nil {
		panic(err)
	}

	valid, err := fs.Create(validConfigPath)

	if err != nil {
		panic(err)
	}
	invalid, err := fs.Create(invalidConfigPath)

	if err != nil {
		panic(err)
	}

	_, err = valid.WriteAt([]byte(validYAML), 0)

	if err != nil {
		panic(err)
	}

	_, err = invalid.WriteAt([]byte(aintValidYAML), 0)

	if err != nil {
		panic(err)
	}

	return fs
}

var pointTo = struct {
	Str      func(s string) *string
	StrSlice func(s []string) *[]string
	Int      func(s int) *int
	IntSlice func(s []int) *[]int
	Bool     func(b bool) *bool
}{
	Str: func(s string) *string {
		return &s
	},
	StrSlice: func(s []string) *[]string {
		return &s
	},
	Int: func(i int) *int {
		return &i
	},
	IntSlice: func(i []int) *[]int {
		return &i
	},
	Bool: func(b bool) *bool {
		return &b
	},
}

type cobraBuilderStub struct {
	T              *testing.T
	buildCallCount int
	BuildErrors    []error
	BuildResults   []interface{}
}

func (cb *cobraBuilderStub) Build(cmd Command, cfg interface{}) (i interface{}, err error) {
	if len(cb.BuildResults) == 0 {
		cb.T.Error("cobraBuilderStub.Build was not expected to be called")
		return
	}

	if cb.buildCallCount > len(cb.BuildResults)-1 {
		cb.T.Error("cobraBuilderStub.Build called too many times")
		return
	}

	i, err = cb.BuildResults[cb.buildCallCount], cb.BuildErrors[cb.buildCallCount]
	cb.buildCallCount++

	return
}
