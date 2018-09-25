package sdctl_context

import (
	"bytes"
	"testing"
)

func TestSdctlConfig_PrintParam(t *testing.T) {
	test_context1 := SdctlContext{
		UserToken: "test_token1",
		APIURL:    "test_url",
		SDJWT:     "test_jwt",
	}
	test_context2 := SdctlContext{
		UserToken: "test_token2",
		APIURL:    "test_url2",
		SDJWT:     "test_jwt2",
	}

	config := SdctlConfig{
		CurrentContext: "default",
		SdctlContexts:  make(map[string]SdctlContext),
	}
	config.SdctlContexts["default"] = test_context1
	config.SdctlContexts["context2"] = test_context2

	cases := map[string]struct {
		paramName string
		expected  string
	}{
		"print user token for default context": {
			UserTokenKey,
			test_context1.UserToken + "\n",
		},
		"print api url for default context": {
			APIURLKey,
			test_context1.APIURL + "\n",
		},
		"print jwt for default context": {
			SDJWTKey,
			test_context1.SDJWT + "\n",
		},
		"print current context": {
			CurrentContextKey,
			"default" + "\n",
		},
		"print context list": {
			ContextsKey,
			"  context2\n* default\n",
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			buf := new(bytes.Buffer)
			config.PrintParam(v.paramName, buf)
			if buf.String() != v.expected {
				t.Errorf("expect='%v', actual='%v'", v.expected, buf.String())
			}
		})
	}
}
