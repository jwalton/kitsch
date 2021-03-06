package modules

import (
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"github.com/jwalton/kitsch/internal/kitsch/projects"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas ProjectModule

// ProjectConfig represents configuration overrides for individual project types
// within the project module.
type ProjectConfig struct {
	// Style is the style to apply to this project.
	Style string `yaml:"style"`
	// ToolSymbol is the symbol to show for this project's build tool.
	ToolSymbol string `yaml:"toolSymbol"`
	// PackageManagerSymbol is the symbol to show for this project's package manager.
	PackageManagerSymbol string `yaml:"packageManagerSymbol"`
}

type projectModuleData struct {
	projectInfo projects.ProjectInfo

	// Name is the name of the matched project type.
	Name string
	// ToolSymbol is the symbol for this project's build tool.
	ToolSymbol string
	// ToolVersion is the version of this project's build tool.
	ToolVersion string
	// PackageManagerSymbol is the symbol for this project's package manager, or "" if unavailable.
	PackageManagerSymbol string
	// ProjectStyle is the style for this project, or "" if none.
	ProjectStyle string
}

// PackageManagerVersion returns the version of the package manager, or "" if unavailable.
func (p projectModuleData) PackageManagerVersion() string {
	return p.projectInfo.PackageManagerVersion()
}

// PackageVersion returns the version of the package in the current folder, or "" if unavailable.
func (p projectModuleData) PackageVersion() string {
	return p.projectInfo.PackageVersion()
}

// ProjectModule prints information about the project in the current folder.
//
type ProjectModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=project"`
	// Projects is project-specific configuration.
	Projects map[string]ProjectConfig `yaml:"projects"`
	// DefaultProjectStyle is the style to use if no project-specific style is specified.
	DefaultProjectStyle string `yaml:"defaultProjectStyle"`
}

// Execute the module.
func (mod ProjectModule) Execute(context *Context) ModuleResult {
	projectInfo := projects.ResolveProjectType(context.ProjectTypes, context)

	if projectInfo == nil {
		return ModuleResult{}
	}

	overrides, ok := mod.Projects[projectInfo.Name]
	if !ok {
		overrides = ProjectConfig{}
	}

	projectStyleString := overrides.Style
	if projectStyleString == "" {
		projectStyleString = mod.DefaultProjectStyle
	}
	if projectStyleString == "" {
		projectStyleString = projectInfo.Style
	}
	projectStyle := context.GetStyle(projectStyleString)

	data := projectModuleData{
		projectInfo:          *projectInfo,
		Name:                 projectInfo.Name,
		ToolSymbol:           defaultString(overrides.ToolSymbol, projectInfo.ToolSymbol),
		ToolVersion:          projectInfo.ToolVersion,
		PackageManagerSymbol: defaultString(overrides.PackageManagerSymbol, projectInfo.PackageManagerSymbol),
		ProjectStyle:         projectStyleString,
	}

	text := ""
	if data.ToolVersion != "" {
		text = "w/" + projectStyle.Apply(data.ToolSymbol+"@"+data.ToolVersion)
	}

	return ModuleResult{DefaultText: text, Data: data}
}

func init() {
	registerModule(
		"project",
		registeredModule{
			jsonSchema: schemas.ProjectModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := ProjectModule{Type: "project"}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
