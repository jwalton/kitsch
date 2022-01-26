// Code generated by "genSchema --pkg schemas KubernetesModule"; DO NOT EDIT.

package schemas

// KubernetesModuleJSONSchema is the JSON schema for the KubernetesModule struct.
var KubernetesModuleJSONSchema = `{
  "type": "object",
  "properties": {
    "style": {"type": "string", "description": ""},
    "template": {"type": "string", "description": ""},
    "symbol": {"type": "string", "description": "Symbol is a symbol to show if a Kubernetes context is detected.  Defaults to \"☸ \""},
    "contextAliases": {"type": "object", "description": "ContextAliases is a map where keys are context names and values are the value we want to show.  If the value is an empty string, we will not show anything.", "additionalProperties": {"type": "string", "description": ""}},
    "configFile": {"type": "string", "description": "ConfigFile is the path to the kubectl config file.  Defaults to \"~/.kube/config\"."}
  },
  "additionalProperties": false
}`

