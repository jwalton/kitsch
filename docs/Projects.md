# Projects

The "project" module works out what kind of project the current folder represents, and displays the current tooling versions. This is done through the "projects" top-level configuration item in `${configdir}/kitsch.yaml`. (If you're not sure where your configdir is, run `kitsch-prompt configdir`). Configuring projects in your config file means that if you have custom project types, they'll automatically apply to any theme you install. Let's have a look at one of the default project types:

```yaml
projectTypes:
  - name: node
    conditions:
      ifFiles: ["package.json"]
    toolSymbol: Node
    packageManagerSymbol: npm
    toolVersion:
      type: custom
      from: "node --version"
    packageManagerVersion:
      type: custom
      from: "npm --version"
    packageVersion:
      type: ancestorFile
      from: package.json
      template: "{{ .version }}"
```

Here we can see that each item in the the "projectTypes" list contains:

- `name` - the name of this project type.
- `conditions` - conditions for when to activate this project.
- `toolSymbol` - the default symbol to show for this project type in the "project" module.
- `toolVersion` - the "getter" to use to get the version of the tool.
- `pacakgeManagerSymbol` (optional) - the name of the package manager used (e.g. "npm").
- `packageManagerSymbol` and `packageVersion` (optional) - "getters" for the package manager version and the package version.

A given folder may be ambiguous in terms of project type - for example if a folder contains a "package.json" and a "go.mod", should we treat it as a node project or as a go project? The "project" module will go through the list of project types, and will return the first one that matches - the order in `projectTypes` defines the precedence.

You can add your own project types. All the default project types will be automatically added in after any custom ones you define. You can also re-order the existing project types, or redefine them in part or in whole. This example will re-order project types so "node" takes precedence over "go", and both will take precedence over all other default project types (which will be added in after these):

```yaml
projectTypes:
  - name: node
  - name: go
```

## Conditions

The following conditions may be specified on a project type:

- `ifFiles: ["file1", "file2", "file3"]` - The project type will be selected if one or more of the files is present in the current folder.
- `ifAncestorFiles: ["file1", "file2", "file3"]` - The project type will be selected if one or more of the files is present in the current folder, or any folder higher up the directory hierarchy.
- `ifExtensions: ["js", "jsx", "ts", "tsx"]` - The project type will be selected if one or more of the extensions is present in the current folder.
- `ifOS: ["darwin", "linux"] - The project will only be shown if the OS is one of the OSs listed.
- `ifNotOS: ["windows"]` - The project will only be shown if the OS is not one of the listed OSs.

Note that "ifFile", "ifAncestorFile", and "ifExtension" are "or"ed together - if any one of these conditions is met, the project type will be selected. The "OS" conditions are "and" conditions - if they are not met, the project will not be selected, even if files or extensions match. Also, note that most of the built-in project types do not use "ifAncestorFile" - this means if you're in a subdirectory of a node project, then the "project" module won't show that you're in a node project. This is considered an "acceptable performance tradeoff" - most of the time developers spend time in the root of a project anyways, since that's where you'd run `make`, or `npm test`, or `go build`, or whatever command you'd use to build your project. You can also easily override this, for example adding this to "config.yaml" will fix this for node projects:

```yaml
projectTypes:
  - name: node
    conditions:
      ifAncestorFile: ["package.json"]
```

### Getters

The `tool`, `packageManager`, and `packageVersion` all take a "getter" object that describes how to get an object. There are many different types of getters, but each is a `{ type, from, as, valueTemplate, regex, cache }` object.

`type` and `from` together define the where the data will be fetched from:

- `{type: "custom", from: "npm --version"}` will execute the command "npm --version", and use the output as the value (with leading and trailing whitespace removed).
- `{type: "file", from: "version.txt"}` will load the contents of a file as the value, if the file exists. "from" is relative to the current folder, and may not start with "/" or contain "." or "..". Files may be located in subdirectories, in which case "/" should be used as a path delimiter, even on windows.
- `{type: "ancestorFile", from: "version.txt"}` will attempt to locate the file "version.txt" in the current folder or any folder higher up the directory hierarchy.
- `{type: "env", from: "GCC_VERSION"}` will get a value from an environment variable.

For any getter, you can also specify:

- `as:` - how to interpret the retrieved value - one of `text`, `json`, `toml`, or `yaml`. For `json`, `toml`, or `yaml`, the file will be parsed and the results passed to the `valueTemplate`. For `text`, the valueTemplate will get a `{ Text }` object. If `as` is specified and `valueTemplate` is not, then this getter will never return a value.
- `valueTemplate` - A template used to render the value. Note that the valueTemplate is passed the parsed object from "as" - `.Globals` are not available in this template, nor are style functions.
- `regex` - A regular expression used to extract a value - if there are any capturing groups in the regex, then the first capturing group will be returned. Otherwise, the matched text will be returned. If `regex` is specified, then `as` will be ignored. If both `regex` and `valueTemplate` are specified, then `valueTemplate` will be run after the regex.

An example of the `regex` option can be seen in this example to fetch the version from the `go version` command. This will print a value like "go version go1.17.1 darwin/amd64", and the regex will extract the "1.17.1" part:

```yaml
projectTypes:
  - name: go
    ifFiles: ["go.mod"]
    ifExtensions: ["go"]
    toolSymbol: go
    toolVersion:
      type: "custom"
      command: "go version"
      regex: "go version go(\d+\.\d+\.\d+)"
```

If the "tool" getter returns an empty value or an error, the project will not be selected, even if it otherwise would have been.  For security reasons, a "custom" getter will ignore "." in the PATH. For example, if the command is "npm --version", and there is an "npm" command in the current working directory, then "custom" will not run `./npm` even if "." is in the PATH.