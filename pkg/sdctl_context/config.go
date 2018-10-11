package sdctl_context

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

var (
	UserTokenKey      = "token"
	APIURLKey         = "api"
	SDJWTKey          = "jwt"
	CurrentContextKey = "current_context"
	ContextsKey       = "contexts"
)

// SdctlContext is param to use Screwdriver.cd API
type SdctlContext struct {
	UserToken string `json:"token"`
	APIURL    string `json:"api"`
	SDJWT     string `json:"jwt"`
}

// SdctlConfig represents the context of Screwdriver.cd
type SdctlConfig struct {
	CurrentContext string                  `json:"current_context"`
	SdctlContexts  map[string]SdctlContext `json:"contexts"`
}

func LoadConfig(configPath string, force bool) (SdctlConfig, error) {
	if _, err := os.Stat(configPath); err != nil || force {
		return initConfigFile(configPath)
	}

	return getSdctlConfigFromFile(configPath)
}

func getSdctlConfigFromFile(path string) (sdctlConfig SdctlConfig, err error) {
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(confFile, &sdctlConfig)
	return
}

func initConfigFile(configPath string) (SdctlConfig, error) {
	context := SdctlContext{
		UserToken: "",
		APIURL:    "",
		SDJWT:     "",
	}
	config := SdctlConfig{
		CurrentContext: "default",
		SdctlContexts:  make(map[string]SdctlContext),
	}
	config.SdctlContexts["default"] = context

	f, _ := json.Marshal(config)
	err := ioutil.WriteFile(configPath, f, 0660)

	return config, err
}

func (sc *SdctlConfig) Update(configPath string) error {
	f, err := json.Marshal(sc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, f, 0660)
}

func (sc *SdctlConfig) PrintParam(paramName string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	var s string
	sdctx := sc.SdctlContexts[sc.CurrentContext]

	switch {
	case paramName == UserTokenKey:
		s = sdctx.UserToken + "\n"
	case paramName == APIURLKey:
		s = sdctx.APIURL + "\n"
	case paramName == SDJWTKey:
		s = sdctx.SDJWT + "\n"
	case paramName == CurrentContextKey:
		s = sc.CurrentContext + "\n"
	case paramName == ContextsKey:
		var contexts []string
		for k := range sc.SdctlContexts {
			contexts = append(contexts, k)
		}
		sort.Slice(contexts, func(i, j int) bool {
			return contexts[i] < contexts[j]
		})
		for _, k := range contexts {
			if k == sc.CurrentContext {
				s = s + "* " + k + "\n"
			} else {
				s = s + "  " + k + "\n"
			}
		}
	default:
		s = "no such config"
	}

	fmt.Fprint(w, s)
}

func (sc *SdctlConfig) SetParam(paramName, param string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	sdctx := sc.SdctlContexts[sc.CurrentContext]

	switch {
	case paramName == UserTokenKey:
		sdctx.UserToken = param
		sc.SdctlContexts[sc.CurrentContext] = sdctx
	case paramName == APIURLKey:
		sdctx.APIURL = param
		sc.SdctlContexts[sc.CurrentContext] = sdctx
	case paramName == SDJWTKey:
		sdctx.SDJWT = param
		sc.SdctlContexts[sc.CurrentContext] = sdctx
	case paramName == CurrentContextKey:
		if _, ok := sc.SdctlContexts[param]; !ok {
			sc.SdctlContexts[param] = SdctlContext{
				UserToken: "",
				APIURL:    "",
				SDJWT:     "",
			}
		}
		sc.CurrentContext = param
	}

	fmt.Fprintf(w, "'%v' is set\n", paramName)
}
