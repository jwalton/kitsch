# genSchema

gneSchema is a JSON schema generator for go structs.  This is not a complete implementation - it only does the things we need for kitsch prompt.

There are other libraries out there that generate JSON schema for structs using reflection.  The biggest advantage of the approach used here is that it happens instead at compile time - there's no need to have an instance of the object around and no need for any runtime reflection.

## Usage

To use, add a `//go:generate go run ../genSchema/main.go [options] [args...]` comment to the file where the struct is defined.  Each "arg" is the name of a structure defined in the same file.

This will generate a file will the same name, suffixed with "_schema", containing a schema definition for each

If `--pkg packagename` is specified in options, then generated schemas will be place in a child package with the given name.  If `--private` is specified, then generated schemas will not be exported.

## Field Types

genSchema supports the basic field types (`string`, `bool`, `int*`, `uint*`, `float*`).  Maps are supported, although keys are always assumed to be strings.

## Struct Tags

Struct tags for each field are parsed out in `parseStructTags()`.

### ref

This can be used on a child struct to use `{ "$ref": "#/definitions/[StructName]" }` instead of inlining the struct.  Note that it is up to you to make sure such a definition exists.

Example:

```go
type MyModule struct {
    Conditions condition.Conditions `yaml:"conditions" jsonschema:",ref"`
    Other      other.Foo            `yaml:"conditions" jsonschema:",ref=FooType"`
}
```

### required

Marks a field as required in the JSON Schema.

```go
type MyModule struct {
    Content string `yaml:"content" jsonschema:",required"`
}
```

### enum

Enum values should be ":" separated.  This will result in the field being of the "string" type in the JSON schema, regardless of the underlying type.

Example:

```go
type MyModule struct {
    Content `yaml:"content" jsonschema:",enum=text:json:toml:yaml"`
}
```
