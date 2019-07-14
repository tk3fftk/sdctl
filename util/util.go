package util

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

// ReadYaml reads yaml file
func ReadYaml(yamlPath string) (yaml string, err error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return
	}
	yaml = fmt.Sprintf("%q", string(yamlFile[:]))
	return
}

// ConfigPATH gets config file path for sdctl
func ConfigPATH() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	configPATH := filepath.Join(usr.HomeDir, "/.sdctl")

	return configPATH, nil
}
