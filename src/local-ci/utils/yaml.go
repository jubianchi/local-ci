package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadYaml(path string) map[interface{}]interface{} {
	dat, err := ioutil.ReadFile(path)

	if nil != err {
		panic(err)
	}

	yml := make(map[interface{}]interface{})

	yaml.Unmarshal(dat, &yml)

	return yml
}
