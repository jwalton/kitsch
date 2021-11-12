package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/projects"
	"gopkg.in/yaml.v3"
)

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
	// PackageManagerSymbol is, if available, the symbol for this project's package manager.
	PackageManagerSymbol string
	// ProjectStyle is the style for this project, if any.
	ProjectStyle string
}

// PackageManagerVersion is, if available, the version of this project's package manager.
func (p projectModuleData) PackageManagerVersion() string {
	return p.projectInfo.PackageManagerVersion()
}

// PackageVersion returns, if available, the version of this project.
func (p projectModuleData) PackageVersion() string {
	return p.projectInfo.PackageVersion()
}

// ProjectModule prints information about the project in the current folder.
//
type ProjectModule struct {
	CommonConfig `yaml:",inline"`
	// Projects is project-specific configuration.
	Projects map[string]ProjectConfig `yaml:"projects"`
}

// Execute the module.
func (mod ProjectModule) Execute(context *Context) ModuleResult {
	directory := context.Directory
	projectInfo := projects.ResolveProjectType(context.ProjectTypes, directory)

	if projectInfo == nil {
		return ModuleResult{}
	}

	overrides, ok := mod.Projects[projectInfo.Name]
	if !ok {
		overrides = ProjectConfig{}
	}

	data := projectModuleData{
		projectInfo:          *projectInfo,
		Name:                 projectInfo.Name,
		ToolSymbol:           defaultString(overrides.ToolSymbol, projectInfo.ToolSymbol),
		ToolVersion:          projectInfo.ToolVersion,
		PackageManagerSymbol: defaultString(overrides.PackageManagerSymbol, projectInfo.PackageManagerSymbol),
		ProjectStyle:         overrides.Style,
	}

	projectStyle, err := context.Styles.Get(data.ProjectStyle)
	if err != nil {
		log.Warn("Invalid style " + data.ProjectStyle + ": " + err.Error())
		projectStyle, _ = context.Styles.Get("")
	}

	text := ""
	if data.ToolVersion != "" {
		text = "via " + projectStyle.Apply(data.ToolSymbol+"@"+data.ToolVersion)
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, text)
}

func init() {
	registerFactory("project", func(node *yaml.Node) (Module, error) {
		var module ProjectModule
		err := node.Decode(&module)
		return &module, err
	})
}
