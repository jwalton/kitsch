// Package projects is used to detect the project type of a directory.
//
package projects

import (
	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
)

// ProjectInfo represents resolved information about the current project.
type ProjectInfo struct {
	projectType   ProjectType
	getterContext getters.GetterContext

	// Name is the name of the matched project type.
	Name string
	// ToolSymbol is the symbol for this project's build tool.
	ToolSymbol string
	// ToolVersion is the version of this project's build tool.
	ToolVersion string
	// PackageManagerSymbol is, if available, the symbol for this project's package manager.
	PackageManagerSymbol        string
	packageManagerVersion       string
	packageVersion              string
	packageManagerVersionLoaded bool
	packageVersionLoaded        bool
}

// PackageManagerVersion is, if available, the version of this project's package manager.
func (projectInfo *ProjectInfo) PackageManagerVersion() string {
	if !projectInfo.packageManagerVersionLoaded {
		str, _ := getStringValue(
			projectInfo.projectType.PackageManagerVersion,
			projectInfo.getterContext,
		)
		projectInfo.packageManagerVersion = str
		projectInfo.packageManagerVersionLoaded = true
	}
	return projectInfo.packageManagerVersion
}

// PackageVersion returns, if available, the version of this project.
func (projectInfo *ProjectInfo) PackageVersion() string {
	if !projectInfo.packageVersionLoaded {
		str, _ := getStringValue(
			projectInfo.projectType.PackageVersion,
			projectInfo.getterContext,
		)
		projectInfo.packageVersion = str
		projectInfo.packageVersionLoaded = true
	}
	return projectInfo.packageVersion
}

func getStringValue(getter getters.Getter, getterContext getters.GetterContext) (string, error) {
	if getter == nil {
		return "", nil
	}

	value, err := getter.GetValue(getterContext)
	if err != nil {
		return "", err
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return "", nil
}

// ResolveProjectType returns the project type for the specified folder, or nil
// if the project type cannot be determined.
func ResolveProjectType(
	projectTypes []ProjectType,
	getterContext getters.GetterContext,
) *ProjectInfo {
	for _, projectType := range projectTypes {
		if !projectType.Conditions.Matches(getterContext.GetWorkingDirectory()) {
			continue
		}

		toolVersion, err := getStringValue(projectType.ToolVersion, getterContext)
		if err != nil || toolVersion == "" {
			// If we can't get a toolVesrion, skip this project type.
			log.Info("Could not get tool version for project type", projectType.Name)
			continue
		}

		return &ProjectInfo{
			projectType:          projectType,
			getterContext:        getterContext,
			Name:                 projectType.Name,
			ToolSymbol:           projectType.ToolSymbol,
			ToolVersion:          toolVersion,
			PackageManagerSymbol: projectType.PackageManagerSymbol,
		}
	}

	return nil
}
