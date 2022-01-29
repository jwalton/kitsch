package modules

import (
	"gopkg.in/yaml.v3"
)

func moduleFromYAML(data string) Module {
	return moduleSpecFromYAML(data).Module
}

func moduleSpecFromYAML(data string) ModuleSpec {
	var moduleSpec ModuleSpec
	err := yaml.Unmarshal([]byte(data), &moduleSpec)
	if err != nil {
		panic(err)
	}
	return moduleSpec
}
