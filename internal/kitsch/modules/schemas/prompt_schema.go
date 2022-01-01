// Code generated by "genSchema --pkg schemas PromptModule"; DO NOT EDIT.

package schemas

// PromptModuleJSONSchema is the JSON schema for the PromptModule struct.
var PromptModuleJSONSchema = `{
  "type": "object",
  "properties": {
    "style": {"type": "string", "description": ""},
    "template": {"type": "string", "description": ""},
    "type": {"type": "string", "description": "Type is the type of this module.", "enum": ["prompt"]},
    "prompt": {"type": "string", "description": "Prompt is what to display as the prompt.  Defaults to \"$ \"."},
    "rootPrompt": {"type": "string", "description": "RootPrompt is what to display as the prompt if the current user is root.  Defaults to \"# \"."},
    "rootStyle": {"type": "string", "description": "RootStyle will be used in place of ` + "`" + `Style` + "`" + ` if the current user is root. If this style is empty, will fall back to Style."},
    "vicmdPrompt": {"type": "string", "description": "ViCmdPrompt is what to display as the prompt if the shell is in vicmd mode. Defaults to \": \"."},
    "vicmdStyle": {"type": "string", "description": "VicmdStyle will be used when the shell is in vicmd mode."},
    "errorStyle": {"type": "string", "description": "ErrorStyle will be used when the previous command failed."}
  },
  "additionalProperties": false
}`

