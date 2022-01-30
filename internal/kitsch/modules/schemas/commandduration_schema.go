// Code generated by "genSchema --pkg schemas CmdDurationModule"; DO NOT EDIT.

package schemas

// CmdDurationModuleJSONSchema is the JSON schema for the CmdDurationModule struct.
var CmdDurationModuleJSONSchema = `{
  "type": "object",
  "properties": {
    "type": {"type": "string", "description": "Type is the type of this module.", "enum": ["command_duration"]},
    "minTime": {"type": "integer", "description": "MinTime is the minimum duration to show, in milliseconds."},
    "showMilliseconds": {"type": "boolean", "description": "ShowMilliseconds - If true, show milliseconds."}
  },
  "required": ["type"]}`

