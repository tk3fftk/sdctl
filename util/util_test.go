package util_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/tk3fftk/sdctl/util"
)

func TestReadYaml(t *testing.T) {
	cases := map[string]struct {
		filePATH string
		expect   string
	}{
		"file exists": {
			filePATH: "../testdata/screwdriver.yaml",
			expect:   readFile("../testdata/screwdriver.yaml"),
		},
		"file is missing": {
			filePATH: "missing",
			expect:   "",
		},
	}

	for k, v := range cases {
		k := k
		v := v
		t.Run(k, func(t *testing.T) {
			actual, _ := util.ReadYaml(v.filePATH)
			if actual != v.expect {
				t.Errorf("actual should be %s, but this is %s", v.expect, actual)
			}
		})
	}
}

func readFile(path string) string {
	b, _ := ioutil.ReadFile(path)
	return fmt.Sprintf("%q", string(b[:]))
}
