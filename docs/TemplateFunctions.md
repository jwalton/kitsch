# Template Functions

## style(style, text)

The `style` template function applies a style to some text.

Examples:

```gotemplate
// Show name in red
{{ .name | style "red" }}

// Show name in bright cyan with with blue background
{{ .name | style "bgBlue brightCyan" }}
```

## fgColor(color, text)

The `fgColor` function applies a foreground color to some text.

```gotemplate
// Show name in red
{{ .name | fgColor "red" }}

// Still shows name in red
{{ .name | fgColor "bg:red" }}
```

## bgColor(color, text)

The `bgColor` function applies a background color to some text.

```gotemplate
// Show name on red background
{{ .name | bgColor "red" }}
```

## newPowerline(prefix, separator, suffix)

`newPowerline` returns a powerline object, used as a helper to render powerline prompts from a template.

The `powerline` object has the following functions:

- `$pl.Segment color text` - Render a new powerline segment with the specific background color and text. Between each segment, the powerline object will render a separator, consisting of the `prefix` (with the last segment's background color and the new segment's color as the foreground) then the `separator + suffix` (with the last segment's color as the foreground and the new segment's background color).
- `$pl.Finish` - Render a "prefix + separator" in the last segment's background color. If there are no previous segments, this will return an empty string.

For example:

```yaml
template: |
  {{- $pl := newPowerline " " "\ue0b0" " " -}}
  {{- with .Data.Modules -}}
    {{- printf " %s" .time.Text | style "$timeFg" | $pl.Segment "$timeBg" -}}
    {{- .directory.Text | style "$directoryFg" | $pl.Segment "$directoryBg" -}}
    {{- $pl.Finish -}}{{- " " -}}
  {{- end -}}
```

## newReversePowerline(prefix, separator, suffix)

`newReversePowerline` is the same as `newPowerline`, except the colors of the "separator" are flipped. This is useful when you want to use the "left-pointing powerline arrow" (`\ue0b0`) for an rprompt.