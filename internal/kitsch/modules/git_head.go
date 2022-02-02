package modules

import (
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas GitHeadModule

// GitHeadModule shows information about the HEAD of the current git repo.
//
type GitHeadModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=git_head"`
	// MaxTagsToSearch is the maximum number of tags to search when checking to
	// see if HEAD is a tagged release.  Defaults to 200.
	MaxTagsToSearch int `yaml:"maxTagsToSearch"`
}

type gitHeadResult struct {
	// HeadDescription is the name of the branch we are currently on if the head
	// is not detached.  If the head is detached, this will be the branch name
	// if we are in the middle of a rebase or merge, the tag name if the head is
	// at a tag, or the short hash otherwise.
	Description string
	// Detached is true if the head is detached.
	Detached bool
	// Hash is the current hash of the head.
	Hash string
	// ShortHash is the short version of the hash.
	ShortHash string
	// Upstream is the name of the upstream branch, or "" if there is no upstream,
	// of if the Head is detached.
	Upstream string
}

// Execute runs a git module.
func (mod GitHeadModule) Execute(context *Context) ModuleResult {
	git := context.Git()

	if git == nil {
		return ModuleResult{DefaultText: "", Data: gitHeadResult{}}
	}

	head, err := git.Head(mod.MaxTagsToSearch)
	if err != nil {
		return ModuleResult{DefaultText: "???", Data: gitHeadResult{
			Description: "???",
			Detached:    true,
			Hash:        "???",
			Upstream:    "",
		}}
	}

	upstream := ""
	if !head.Detached {
		upstream = git.GetUpstream(head.Description)
	}

	shortHash := head.Hash
	if len(shortHash) > 7 {
		shortHash = shortHash[0:7]
	}

	return ModuleResult{DefaultText: head.Description, Data: gitHeadResult{
		Description: head.Description,
		Detached:    head.Detached,
		Hash:        head.Hash,
		ShortHash:   shortHash,
		Upstream:    upstream,
	}}
}

func init() {
	registerModule(
		"git_head",
		registeredModule{
			jsonSchema: schemas.GitHeadModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := GitHeadModule{
					Type:            "git_head",
					MaxTagsToSearch: 200,
				}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
