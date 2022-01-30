// Code generated by "genSchema --pkg schemas TextModule"; DO NOT EDIT.

package schemas

// TextModuleJSONSchema is the JSON schema for the TextModule struct.
var TextModuleJSONSchema = `{
  "type": "object",
  "properties": {
    "type": {"type": "string", "description": "Type is the type of this module.", "enum": ["text"]},
    "text": {"type": "string", "description": "Text is the text to print."}
  },
  "required": ["type", "text"]}`

