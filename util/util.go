package util

import (
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
	yaml = string(yamlFile)

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
