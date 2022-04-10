package modules

import (
	"strings"

	"github.com/jwalton/go-ansiparser"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
)

// We can't work out the width of a flexible space, until after we've rendered
// the complete prompt.  We use this marker, which is extremely unlikely to
// occur naturally in the output of any module, as a sentinel for the flexible
// space, and then replace it at the end after the entire prompt is rendered.
//
// An alternative to using a sentinel value here would be to let the block
// module handle a flexible space.  We'd have to somehow mark which children were
// flexible spaces, which is a very easy to solve problem.  However, if
// a block had a flexible space and a template, the flexible space wouldn't "survive"
// through the template (how would we represent a flexible space in the output
// of a template?) and would end up being ignored.  Second, with nested blocks,
// which block should handle the flexible space?  Should the nearest ancestor
// fill out the line to the terminal width?
//
// This approach is a lot easier to use, although it has the disadvantage that
// if some other module accidentally or maliciously produces the `flexibleSpaceMarker`,
// we'll end up printing a bunch of spaces where we shouldn't.
const flexibleSpaceMarker = "\t \u00a0/\\;:flex:;\\/\u00a0 \t"

//go:generate go run ../genSchema/main.go --pkg schemas FlexibleSpaceModule

// FlexibleSpaceModule inserts a flexible-width space.
//
type FlexibleSpaceModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=flexible_space"`
}

// Execute the flexible space module.
func (mod FlexibleSpaceModule) Execute(context *Context) ModuleResult {
	return ModuleResult{DefaultText: flexibleSpaceMarker, Data: map[string]interface{}{}}
}

func init() {
	registerModule(
		"flexible_space",
		registeredModule{
			jsonSchema: schemas.FlexibleSpaceModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				var module FlexibleSpaceModule = FlexibleSpaceModule{
					Type: "flexible_space",
				}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}

// processFlexibleSpaces replaces flexible space markers with spaces.
func processFlexibleSpaces(terminalWidth int, renderedPrompt string) string {
	result := ""

	// Split into lines.
	lines := strings.Split(renderedPrompt, "\n")
	for lineIndex, line := range lines {

		// For each line, split into segments around the FlexibleSpaceMarker.
		if strings.Contains(line, flexibleSpaceMarker) {
			segments := strings.Split(line, flexibleSpaceMarker)

			segmentsTotalLength := 0
			for _, segment := range segments {
				segmentsTotalLength += getPrintWidth(segment)
			}

			extraSpace := terminalWidth - segmentsTotalLength

			if extraSpace > 0 {
				spacesAdded := 0
				spacesPerSegment := extraSpace / (len(segments) - 1)

				line = ""
				for index := 0; index < len(segments)-2; index++ {
					line += segments[index]
					line += strings.Repeat(" ", spacesPerSegment)
					spacesAdded += spacesPerSegment
				}

				line += segments[len(segments)-2]
				line += strings.Repeat(" ", extraSpace-spacesAdded)
				line += segments[len(segments)-1]
			}
		}

		result += line
		if lineIndex < len(lines)-1 {
			result += "\n"
		}
	}

	return result
}

func getPrintWidth(str string) int {
	width := 0
	tokenizer := ansiparser.NewStringTokenizer(str)

	for tokenizer.Next() {
		token := tokenizer.Token()

		if token.Type == ansiparser.String {
			if token.IsASCII {
				width += len(token.Content)
			} else {
				width += runewidth.StringWidth(token.Content)
			}
		}
	}

	return width
}
