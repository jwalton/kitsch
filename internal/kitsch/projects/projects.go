// Package projects is used to detect the project type of a directory.
//
package projects

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/log"
)

// ProjectInfo represents resolved information about the current project.
type ProjectInfo struct {
	projectType   ProjectType
	getterContext getters.GetterContext

	// Name is the name of the matched project type.
	Name string
	// Style is the default style for the matched project type.
	Style string
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

func getStringValue(getter []getters.Getter, getterContext getters.GetterContext) (string, error) {
	if len(getter) == 0 {
		return "", nil
	}

	for index := range getter {

		value, err := getter[index].GetValue(getterContext)
		if err != nil {
			log.Warn("Error running getter:", err)
			continue
		}

		if str, ok := value.(string); ok {
			return str, nil
		}
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
			// If we can't get a toolVersion, skip this project type.
			log.Info("Could not get tool version for project type", projectType.Name)
			continue
		}

		return &ProjectInfo{
			projectType:          projectType,
			getterContext:        getterContext,
			Name:                 projectType.Name,
			Style:                projectType.Style,
			ToolSymbol:           projectType.ToolSymbol,
			ToolVersion:          toolVersion,
			PackageManagerSymbol: projectType.PackageManagerSymbol,
		}
	}

	return nil
}

//JSONSchemaDefinitions is a string containing JSON schema definitions for objects in the projects package.
var JSONSchemaDefinitions = "\"ProjectType\": " + projectTypeJSONSchema + ",\n" +
	heredoc.Doc(`"GetterList": {
	  "oneOf": [
	    { "$ref": "#/definitions/CustomGetter" },
	    { "type": "array", "items": { "$ref": "#/definitions/CustomGetter" } }
	  ]
	}`)
