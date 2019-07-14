package util

import (
	"io/ioutil"
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
