---
sidebar_position: 1
---

# Modules

## Common Module Configuration

There are certain configuration items that are available on all modules:

- `style` a the [style string](/docs/styles) to apply to the entire module output.
- `template` is a golang template used to render the result of the module.

TODO: Add documentation about templates here.

## block

The "block" module is used to group a collection of modules together, and concatenate their results. By default, the block module will execute all child modules, then join together their output with " "s inbetween. Any child module that produces no output will be ignored.

The "join" can be specified using a template, so you can control how child modules are joined together. The block module also allows you to combine output from multiple modules using a single template; the `.Data.Modules` object is a map where keys are the `id`s (or `type` for modules that don't specify an `id`) of child modules, and the values are the output from those modules. For example:

```yaml
type: block
  modules:
    - type: hostname
    - type: username
      id: user
  template: |
    {{- printf "%s@%s" .Data.hostname.Hostname .Data.user.Username -}}
```

would print the hostname and username, joined by a "@".

Configuration:

- `modules` is an array of other modules to instantiate. Each module in `modules` may have an `id`, which is a string that will uniquely identify the module.
- `join=""` is a template to use to join together modules. Note that if `template` is specified, the `join` parameter will be ignored. The join template will be passed the following parameters:
  - `Globals` are the global variables available in any template.
  - `PrevColors` is an `{FG, BG}` object containing color strings for the previous module's end style.
  - `NextColors` is an `{FG, BG}` object containing color strings for the next module's start style.
  - `Index (int)` is the index of the next module in the Modules array.

Outputs:

- `Modules` is a map of results from executing each child module, indexed by module ID. Only modules that actually generated output will be included. If a module does not have an ID, then the module's `type` will be used to index the module results.
- `ModuleArray` is an array of results from executing each child module. Only modules that actually generated output will be included.

## custom

The "custom" module runs a command and returns the result. If the `as` parameter is specified as "json", "toml", or "yaml", then the output of the command will be parsed according to the specified format. In this case, you must provide a `template` parameter to extract the values you need out of the data.

Configuration:

- `from` is the command to run (e.g. "docker --version").
- `as="text"` indicates how the output should be interpreted. This must be one of "text", "json", "toml", or "yaml".
- `regex=""` is a regular expression used to parse values out of the result of the getter (e.g. "^Docker version (._), build ._$"). If specified, then "as" will be ignored.
- `cache={ enabled: false }` controls caching. If `cache.enabled` is true, then the module will resolve the full path of the executable (following any sym-links), and then use the full path, the last modified date, the size of the command, and the arguments as a cache key. Caches are written to the "cache" subfolder in your configuration directory. Caching means that, if we're interested in what version of npm is installed, we only need to run `npm --version` if and when the `npm` executable changes.

Outputs:

The `.Data` value returned from a custom module depends on the `as` configuration. If `as` is "text", then `.Data` will be a `{ Text: [string] }` object, containing the text returned from the command (with leading and trailing whitespace automatically stripped). If `as` is any other value, then the `.Data` object will be the parsed results of the output. For example if `as="json"`, and the returned value was '{"foo": "bar"}', then `.Data.foo` would be "bar".

## command_duration

The "command_duration" module shows the amount of time the previous command took to execute.

Configuration:

- `minTime=2000` the minimum duration to show, in milliseconds.
- `showMilliseconds=false` if true, show to millisecond precision instead of to second precision.

Outputs:

- `Duration (int64)` is the duration the command took, in milliseconds.
- `PrettyDuration (string)` is the duration the command took, in a human-readable format (e.g. "3m21s").

## directory

The "directory" module shows the current working directory. In the default configuration, the directory module will truncate the path if you are more than three directories deep. For example, if you were in "/tmp/foo/bar/baz/qux", ths would show `…/bar/baz/qux`. On windows machines, the volume will always be shown (e.g. `C:\…\bar\baz\qux`). If you are currently in a git directory, everything before the root of the git directory will be stripped.

Configuration:

- `homeSymbol="~"` is the symbol to replace the home directory with when you are in a subdirectory.
- `readOnlySymbol="🔒"` is the symbol to append to the directory if it is read-only.
- `truncateToRepo=true` controls whether or not we truncate to the root of a source code repository. If this is true, and you are in a git repo, we'll remove everything before the root of the source code repository, and prepend `RepoSymbol`.
- `repoSymbol=""` is a string that will be added as a prefix when we truncate to a repo.
- `truncationLength=3` is the maximum number of directories to show. If 0, truncation will be disabled.
- `truncationSymbol="…"` will be added to the start of the string in place of any paths that were removed.

Outputs:

- `Path (string)` is the path that will be shown to the user.
- `PathSeparator (string)` is the system defined path separator.
- `ReadOnly (boolean)` is true if the current directory is read-only.
- `ReadOnlySymbol (string)` is the same as ReadOnlySymbol from the module configuration.

## file

The "file" module reads a file and uses the contents to produce an output. The configuration and outputs of the "file" module are identical to the ["custom"](#custom) module, except that `from` should be the name of a file in the current folder, or the relative path of a file in a subdirectory of the current folder.

## git_status

The git_status module shows information about the status of the current git repo. The default output of this module is based on [posh-git](https://github.com/dahlbyk/posh-git) and [posh-git-sh](https://github.com/lyze/posh-git-sh). If you're in a git repo that has changes, the output will be something like:

```text
+A ~B -C !D | +E ~F -G
```

Where:

- `+A` is the number of unstaged new files.
- `~B` is the number of unstaged modified files.
- `-C` is the number of unstaged removed files.
- `!D` is the number of unmerged/conflicting files.
- `+E` is the number of staged new files.
- `~F` is the number of staged modified files.
- `-G` is the number of staged removed files.

The number of unmerged paths is not shown if it is 0. The "unstaged" and "staged" halves of this are also hidden if all values are zero. By default, unstaged counts are shown in red an staged in green, to mimic the output colors of `git status`.

TODO: Show how you could use a template to make this look like starship prompt.

Configuration:

- `indexStyle (string)` is the style to use for the staged status.
- `unstagedStyle (string)` is the style to use for the unstaged file status.
- `stashStyle (string)` is the style to use for the stash count.

Outputs:

- `Index` is a `{ Added, Modified, Deleted, Total }` object. Each is an `int` representing the number of staged files in that state.
- `Unstaged` is a `{ Added, Modified, Deleted, Total }` object. Each is an `int` representing the number of unstaged files in that state.
- `Unmerged (int)` is the total number of unmerged paths in the git repo.
- `StashCount (int)` is the number of stashes in the git repo.

## git

TODO

## hostname

The hostname module shows the current hostname. By default, this will only display anything if the user is currently logged in via SSH.

Configuration:

- `showAlways=false` will cause the hostname to always be shown. If false, then the hostname will only be shown if the current session is an SSH session.

Outputs:

- `Hostname (string)` is the current hostname.
- `IsSSH (bool)` is true if this is an SSH session, false otherwise.
- `Show (bool)` is true if we should show the hostname, false otherwise.

## jobs

The jobs module shows the current count of running background jobs. If the number of running jobs is greater than or equal to `SymbolThreshold` then the `Symbol` will be shone. If the number is greater than or equal to `CountThreshold` then the count of running jobs will be shown.

Configuration:

- `symbol="+"` is the symbol to show when there are background jobs.
- `symbolThreshold=1` is the threshold for showing the symbol.
- `countThreshold=2` is the threshold for showing the count of background jobs.

Outputs:

- `Jobs (int)` is the count of running jobs.
- `ShowSymbol (bool)` is true if the symbol should be shown.
- `ShowCount (bool)` is true if the count should be shown.

## project

The project module works out what kind of project the current folder represents, and displays the current tooling versions. This is done through the ["projects" top-level configuration item](../projects.mdx) in `${configdir}/kitsch.yaml`.

Configuration:

- `projects` is a map where keys are project names, and values are `{ style, toolSymbol, packageManagerSymbol }` objects, which can be used to provide a custom style and symbols for existing projects on a theme-by-theme basis.

Outputs:

- `Name (string)` is the name of the matched project type.
- `ToolSymbol (string)` is the symbol for this project's build tool.
- `ToolVersion (string)` is the version of this project's build tool
- `PackageManagerSymbol (string)` is the symbol for this project's package manager, or "" if unavailable.
- `ProjectStyle (string)` is the style for this project, or "" if none.
- `PackageManagerVersion (string)` is the version of the package manager, or "" if unavailable.
- `PackageVersion (string)` is the version of the package in the current folder, or "" if unavailable.

## prompt

The prompt module displays a "$", or a "#" if the current user is root.

Configuration:

- `prompt="$ "` is what to display as the prompt.
- `rootPrompt="# "` is what to display as the prompt if the current user is root.
- `rootStyle=""` will be used in place of `style` if the current user is root. If this style is empty, will fall back to `style`.
- `viCmdPrompt=": "` is what to display as the prompt if the shell is in vicmd mode.
- `vicmdStyle=""` will be used in place of `style` when the shell is in vicmd mode.
- `errorStyle=""` will be used in place of `style` when the previous command failed.

Outputs:

- `PromptString (string)` is the chosen prompt string, before styling.
- `PromptStyle (string)` is the chosen prompt style.
- `ViCmdMode (bool)` is true if the shell is in vicmd mode (when `.Globals.Keymap == "vicmd").

## text

The text module shows some text.

Configuration:

- `text=""` is the text to show.

Outputs:

- `Text (string)` is the text to show.

## time

The time module shows the current time.

Configuration:

- `layout="15:04:05"` is the format to show the time in. Layout defines the format by showing how the reference time, defined to be `Mon Jan 2 15:04:05 -0700 MST 2006`. The default, "15:04:05" shows the time in 24-hour time. See [the Go time package](https://golang.org/pkg/time/#Time.Format) for more details.

Outputs:

- `Time (time.Time)` is the current time, as a `time.Time` object.
- `Unix (int64)` is the number of seconds since the Unix epoch.
- `TimeStr (string)` is the current time as a formatted string.

## username

The username module shows the current user's username. By default, this will only display anything if the user is currently logged in via SSH. The username is looked up by first checking the `USER` environment variable. If this is empty, the user will be looked up from the OS.

Configuration:

- `showAlways=false` will cause the hostname to always be shown. If false, then the hostname will only be shown if the current session is an SSH session.
- `rootStyle=""` will be used in place of `style` if the current user is root. If this style is empty, will fall back to `style`.

Outputs:

- `Username (string)` is the current user's username.
- `IsSSH (bool)` is true if this is an SSH session, false otherwise.
- `Show (bool)` is true if we should show the hostname, false otherwise.
