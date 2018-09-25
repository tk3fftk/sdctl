package sdctl_context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	testToken   = "test_token1"
	testAPIURL  = "test_url1"
	testSDJWT   = "test_jwt1"
	testContext = "test_context"
	newToken    = "new_token1"
	newAPIURL   = "new_url1"
	newSDJWT    = "new_jwt1"
	newContext  = "new_context"
	tmpPath     = "testdata/sdctl_test"
)

func createMockSdctlConfig() SdctlConfig {
	testContext1 := SdctlContext{
		UserToken: testToken,
		APIURL:    testAPIURL,
		SDJWT:     testSDJWT,
	}
	testContext2 := SdctlContext{
		UserToken: "test_token2",
		APIURL:    "test_url2",
		SDJWT:     "test_jwt2",
	}

	config := SdctlConfig{
		CurrentContext: "default",
		SdctlContexts:  make(map[string]SdctlContext),
	}
	config.SdctlContexts["default"] = testContext1
	config.SdctlContexts[testContext] = testContext2

	return config
}

func emptySdctlContext() SdctlContext {
	return SdctlContext{
		UserToken: "",
		APIURL:    "",
		SDJWT:     "",
	}
}

func TestSdctlConfig_LoadConfig(t *testing.T) {
	cases := map[string]struct {
		hasFile bool
		force   bool
	}{
		"do first init": {
			false,
			false,
		},
		"exist config file but force init": {
			true,
			true,
		},
		"use existing config file": {
			true,
			false,
		},
		"do first init but force init": {
			false,
			false,
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			testLoadConfig(t, v.hasFile, v.force)
		})
	}
}

func testLoadConfig(t *testing.T, hasFile, force bool) {
	t.Helper()

	context := SdctlContext{
		UserToken: "",
		APIURL:    "",
		SDJWT:     "",
	}
	expectConfig := SdctlConfig{
		CurrentContext: "default",
		SdctlContexts:  make(map[string]SdctlContext),
	}
	expectConfig.SdctlContexts["default"] = context

	if hasFile && !force {
		mockConfig := createMockSdctlConfig()
		f, err := json.Marshal(mockConfig)
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(tmpPath, f, 0660)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpPath)

		config, err := LoadConfig(tmpPath, force)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(mockConfig, config, nil) {
			t.Errorf("expected='%v', actual='%v'", mockConfig, config)
		}
	} else {
		config, err := LoadConfig(tmpPath, force)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpPath)

		if !cmp.Equal(expectConfig, config, nil) {
			t.Errorf("expected='%v', actual='%v'", expectConfig, config)
		}
	}
}

func TestSdctlConfig_PrintParam(t *testing.T) {
	config := createMockSdctlConfig()

	cases := map[string]struct {
		paramName string
		expected  string
	}{
		"print user token for default context": {
			UserTokenKey,
			testToken + "\n",
		},
		"print api url for default context": {
			APIURLKey,
			testAPIURL + "\n",
		},
		"print jwt for default context": {
			SDJWTKey,
			testSDJWT + "\n",
		},
		"print current context": {
			CurrentContextKey,
			"default" + "\n",
		},
		"print context list": {
			ContextsKey,
			"* default\n  " + testContext + "\n",
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			t.Helper()

			buf := new(bytes.Buffer)
			config.PrintParam(v.paramName, buf)
			if buf.String() != v.expected {
				t.Errorf("expect='%v', actual='%v'", v.expected, buf.String())
			}
		})
	}
}

func TestSdctlConfig_SetParam(t *testing.T) {
	cases := map[string]struct {
		paramName       string
		param           string
		expectedMessage string
	}{
		"set UserToken": {
			UserTokenKey,
			newToken,
			fmt.Sprintf("'%v' is set\n", UserTokenKey),
		},
		"set APIURL": {
			APIURLKey,
			newAPIURL,
			fmt.Sprintf("'%v' is set\n", APIURLKey),
		},
		"set SDJWT": {
			SDJWTKey,
			newSDJWT,
			fmt.Sprintf("'%v' is set\n", SDJWTKey),
		},
		"set CurrentContext and it's new context": {
			CurrentContextKey,
			newContext,
			fmt.Sprintf("'%v' is set\n", CurrentContextKey),
		},
		"set CurrentContext and it's existing context": {
			CurrentContextKey,
			testContext,
			fmt.Sprintf("'%v' is set\n", CurrentContextKey),
		},
	}

	for k, v := range cases {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {
			testSetParam(t, v.paramName, v.param, v.expectedMessage)
		})
	}
}

func testSetParam(t *testing.T, paramName, param, expectedMessage string) {
	t.Helper()

	config := createMockSdctlConfig()
	f, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(tmpPath, f, 0660)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpPath)

	buf := new(bytes.Buffer)
	config.SetParam(paramName, param, buf)
	switch {
	case paramName == UserTokenKey:
		if config.SdctlContexts[config.CurrentContext].UserToken != param {
			t.Errorf("expect='%v', actual='%v'", param, config.SdctlContexts[config.CurrentContext].UserToken)
		}
	case paramName == APIURLKey:
		if config.SdctlContexts[config.CurrentContext].APIURL != param {
			t.Errorf("expect='%v', actual='%v'", param, config.SdctlContexts[config.CurrentContext].APIURL)
		}
	case paramName == SDJWTKey:
		if config.SdctlContexts[config.CurrentContext].SDJWT != param {
			t.Errorf("expect='%v', actual='%v'", param, config.SdctlContexts[config.CurrentContext].SDJWT)
		}
	case paramName == CurrentContextKey:
		if config.CurrentContext != param {
			t.Errorf("expect='%v', actual='%v'", param, config.CurrentContext)
		}
		// TODO handle this testcase
		if param == newContext {
			if !cmp.Equal(emptySdctlContext(), config.SdctlContexts[config.CurrentContext], nil) {
				t.Errorf("expected='%v', actual='%v'", emptySdctlContext(), config)
			}
		} else {
			if !cmp.Equal(config.SdctlContexts[param], config.SdctlContexts[config.CurrentContext], nil) {
				t.Errorf("expected='%v', actual='%v'", emptySdctlContext(), config)
			}
		}

	}
