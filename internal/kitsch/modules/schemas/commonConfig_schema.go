// Code generated by "genSchema --pkg schemas CommonConfig"; DO NOT EDIT.

package schemas

// CommonConfigJSONSchema is the JSON schema for the CommonConfig struct.
var CommonConfigJSONSchema = `{
  "type": "object",
  "properties": {
    "type": {"type": "string", "description": "Type is the type of this module."},
    "id": {"type": "string", "description": "ID is a unique identifier for this module.  IDs are unique only within the parent block."},
    "style": {"type": "string", "description": "Style is the style to apply to this module."},
    "template": {"type": "string", "description": "Template is a golang template to use to render the output of this module."},
    "conditions":     {
      "type": "object",
      "properties": {
        "ifAncestorFiles": {"type": "array", "description": "", "items": {"type": "string", "description": ""}},
        "ifFiles": {"type": "array", "description": "", "items": {"type": "string", "description": ""}},
        "ifExtensions": {"type": "array", "description": "", "items": {"type": "string", "description": ""}},
        "ifOS": {"type": "array", "description": "", "items": {"type": "string", "description": ""}},
        "ifNotOS": {"type": "array", "description": "", "items": {"type": "string", "description": ""}}
      },
      "additionalProperties": false}
  }}`

