package modules

import (
	"gopkg.in/yaml.v3"
)

func moduleFromYAML(data string) (Module, error) {
	var moduleSpec ModuleSpec
	err := yaml.Unmarshal([]byte(data), &moduleSpec)
	if err != nil {
		return nil, err
	}
	return moduleSpec.Module, nil
}

func moduleFromYAMLMust(data string) Module {
	module, err := moduleFromYAML(data)
	if err != nil {
		panic(err)
	}
	return module
}
