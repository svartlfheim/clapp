package clapp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	fs := buildMockFs()
	ctx := contextStub{
		Vals: map[interface{}]interface{}{
			FsContextKey: fs,
		},
	}

	tests := []struct {
		name                  string
		inputConf             *testConf
		expectedConfStructure interface{}
		appName               string
		filePath              string
		fileMustExist         bool
		expectedConfig        Config
		expectedErr           error
		expectedErrTxt        string
		opts                  []configOpt
	}{
		{
			name:                  "returns config struct when no file is found",
			inputConf:             &testConf{},
			expectedConfStructure: testConf{},
			appName:               "blah",
			expectedConfig: Config{
				appName:         "blah",
				filePath:        "./blah.yaml",
				fs:              fs,
				configMustExist: false,
			},
			expectedErr: nil,
		},
		{
			name:                  "filepath option changes the config value",
			inputConf:             &testConf{},
			expectedConfStructure: testConf{},
			appName:               "blah",
			filePath:              "/tmp/some/dir/blah.yaml",
			expectedConfig: Config{
				appName:         "blah",
				filePath:        "/tmp/some/dir/blah.yaml",
				fs:              fs,
				configMustExist: false,
			},
			expectedErr: nil,
		},
		{
			name:                  "file must exist option causes error for missing file",
			inputConf:             &testConf{},
			expectedConfStructure: testConf{},
			appName:               "blah",
			fileMustExist:         true,
			expectedConfig: Config{
				appName:         "blah",
				filePath:        "./blah.yaml",
				fs:              fs,
				configMustExist: true,
			},
			expectedErr: ErrConfigNotFound,
		},
		{
			name:                  "multiple options are applied",
			inputConf:             &testConf{},
			expectedConfStructure: testConf{},
			appName:               "blah",
			filePath:              "/tmp/some/dir/blah.yaml",
			fileMustExist:         true,
			expectedConfig: Config{
				appName:         "blah",
				filePath:        "/tmp/some/dir/blah.yaml",
				fs:              fs,
				configMustExist: true,
			},
			expectedErr: ErrConfigNotFound,
		},
		{
			name:      "values should be overridden from file",
			inputConf: &testConf{},
			// See validYAML variable
			expectedConfStructure: testConf{
				ShouldEnableThis: "awesome-feature",
				ShouldDoThat:     "work-properly",
				ExternalEndpoint: testExtEndpoint{
					Protocol: "https",
					Domain:   "somefakesite.com",
					Port:     6738,
				},
				ListOfThings: []string{
					"thing1",
					"thing2",
					"thing3",
				},
			},
			appName:  "blah",
			filePath: validConfigPath,
			expectedConfig: Config{
				appName:         "blah",
				filePath:        validConfigPath,
				fs:              fs,
				configMustExist: false,
			},
			expectedErr: nil,
		},
		{
			name:      "error is returned if file is not valid yaml",
			inputConf: &testConf{},
			// See validYAML variable
			expectedConfStructure: testConf{},
			appName:               "blah",
			filePath:              invalidConfigPath,
			expectedConfig: Config{
				appName:         "blah",
				filePath:        invalidConfigPath,
				fs:              fs,
				configMustExist: false,
			},
			expectedErrTxt: "could not unmarshal config: ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			cfg, err := newConfigManager(ctx, test.inputConf, test.appName, test.filePath, test.fileMustExist)

			if test.expectedErrTxt == "" {
				assert.Equal(tt, test.expectedErr, err)
			} else {
				assert.Contains(tt, err.Error(), test.expectedErrTxt)
			}
			assert.Equal(tt, test.expectedConfig, *cfg)
			assert.Equal(tt, test.expectedConfStructure, *test.inputConf)
		})
	}
}

func TestNewConfig_OverrideWithEnv_FailsForNonPointerValue(t *testing.T) {
	ctx := contextStub{
		Vals: map[interface{}]interface{}{
			FsContextKey: buildMockFs(),
		},
	}
	inputConf := &testConf{}
	appName := "blah"
	cfg, err := newConfigManager(ctx, inputConf, appName, validConfigPath, false)

	assert.Nil(t, err)

	err = cfg.OverrideWithEnvVars(*inputConf)

	assert.Equal(t, ErrCannotUseNonPointerValue, err)
}

func TestNewConfig_OverrideWithEnv(t *testing.T) {
	ctx := contextStub{
		Vals: map[interface{}]interface{}{
			FsContextKey: buildMockFs(),
		},
	}
	inputConf := &testConf{}
	appName := "blah"
	cfg, err := newConfigManager(ctx, inputConf, appName, validConfigPath, false)

	assert.Nil(t, err)
	assert.Equal(t, testConf{
		ShouldEnableThis: "awesome-feature",
		ShouldDoThat:     "work-properly",
		ExternalEndpoint: testExtEndpoint{
			Protocol: "https",
			Domain:   "somefakesite.com",
			Port:     6738,
		},
		ListOfThings: []string{
			"thing1",
			"thing2",
			"thing3",
		},
	}, *inputConf)

	os.Setenv("BLAH_SHOULD_ENABLE_THIS", "some-override")
	defer os.Remove("BLAH_SHOULD_ENABLE_THIS")

	os.Setenv("BLAH_EXTERNAL_ENDPOINT_PROTOCOL", "tcp")
	defer os.Remove("BLAH_EXTERNAL_ENDPOINT_PROTOCOL")

	err = cfg.OverrideWithEnvVars(inputConf)

	assert.Nil(t, err)
	assert.Equal(t, testConf{
		// This changed
		ShouldEnableThis: "some-override",
		ShouldDoThat:     "work-properly",
		ExternalEndpoint: testExtEndpoint{
			// This changed
			Protocol: "tcp",
			Domain:   "somefakesite.com",
			Port:     6738,
		},
		ListOfThings: []string{
			"thing1",
			"thing2",
			"thing3",
		},
	}, *inputConf)
}
