# Modules

## Global Module Configuration

- `style` - Style to apply to this module. May be overridden by the module in certain situations (e.g. the "username" module will use `rootStyle` instead of `style` if the current user is root).
- `template` - Specify a golang template to override the rendering of this module.

## Hostname Module

### Configuration

- `always` - If true, hostname will always be shown. If false, hostname will only be shown
  if `.isRemoteSession` is true.

### Template Variables

- `.hostname` - The name of the host.
- `.isRemoteSession` - True if the user is currently logged in via SSH (SSH_CONNECTION, SSH_CLIENT, or SSH_TTY environment variables are set).

## Project Module

The "project" module is used to show information about the build tools used for the project in the current folder. For example, if this is a node.js project, we might want to show "via node@14.17.0" or "via node@14.17.0/npm@v6.14.13", or if this is a rust project we might show "via Rust v1.47.0". The project module detects what kind of project the current folder is using the "projects" configuration.

Note that the "projects" module does not show the Python virtualenv or the current conda environment - these are handled by separate modules.

Unlike most modules, configuration for the project module is split between the module itself, and the [`projectTypes` top-level configuration item](./Projects.md).

### Configuration

- `projects` - a map where keys are project types (e.g. "nodejs", "go", "ruby", etc...), and values are `{ style, toolSymbol, packageManagerSymbol }` objects.
- `template` - defaults to `via {{ .Data.ToolSymbol }}@{{ .Data.ToolVersion }}.

### Template Variables

- `Name` is the name of the project type.
- `ToolSymbol` and `ToolVersion` for the compiler/interpreter.
- `PackageManagerSymbol` and `PackageManagerVersion` (optional) for the default package management utility, if any (e.g. "npm").
- `PackageVersion` (optional) for the version of the current package (e.g. from "package.json").
- `ProjectStyle` is the style associated with the project type from the configuration, in `projects[type].style`.
