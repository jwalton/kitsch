# Template Functions

## style(...styles, text)

The `style` template function takes one or more styles and applies it to some text.  Each style can be either a string containing a single style or multiple styles separated by spaces, or a style object.

Examples:

```gotemplate
// Show name in red
{{ .name | style "red" }}

// Show name in bright cyan with with blue background
{{ .name | style "bgBlue" "brightCyan" }}

// Same as above
{{ .name | style "bgBlue brightCyan" }}
```

## fgColor(color, text)

The `fgColor` function applies a foreground color to some text.

```gotemplate
// Show name in red
{{ .name | fgColor "red" }}

// Still shows name in red
{{ .name | fgColor "bgRed" }}
```

## bgColor(color, text)

The `bgColor` function applies a background color to some text.

```gotemplate
// Show name on red background
{{ .name | bgColor "red" }}

// Same as above
{{ .name | bgColor "bgRed" }}
```
