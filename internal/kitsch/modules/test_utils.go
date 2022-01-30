package modules

import (
	"gopkg.in/yaml.v3"
)

func moduleFromYAML(data string) Module {
	return moduleWrapperFromYAML(data).Module
}

func moduleWrapperFromYAML(data string) ModuleWrapper {
	var moduleWrapper ModuleWrapper
	err := yaml.Unmarshal([]byte(data), &moduleWrapper)
	if err != nil {
		panic(err)
	}
	return moduleWrapper
}
