# Styles

Styling text in kitsch-prompt is done using "style strings". A style string can contain a color, a background color, and multiple modifiers. For example: `brightCyan bg:red bold` would make bold bright-cyan text on a red background.

## Colors

A valid color is any of the following:

- A color name - one of "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white".  Basic color names will set color using "16-color" ANSI.
- A color name prefixed with "bright" (e.g. "brightRed").
- A hex color (e.g. "#fff" or "#30a9fb").
- A CSS3 color name (other than those listed above).
- A CSS-style linear gradient (e.g. "linear-gradient(#000, #fff)") See more on "linear-gradients" below.
- Any of the above prefixed with "bg:" to set the background color.

The style string "red" would set the foreground color to red, the string "bg:red" would set the background to red.

Note that when using color names, the color will be set with a 16-color ANSI code.  The exact color that will be shown for "red" will depend on your terminal.  iTerm2 on a Mac would default "red" to "#c91b00", where the Windows 10 Console would show "red" as "#c50f1f".

### linear-gradient

A linear-gradient is specified exactly like one would be in CSS.  The only difference is that you may not set the direction of the gradient - it is always left-to-right.  A linear-gradient can have any number of stops, and stop positions may be specified as relative positions (e.g. "20%") or with absolute positions (e.g. "3px" - each character is considered 1px wide, since we can only set the color of an entire character), or even with a mix of the two.  Gradients can be applied to the background by prefixing them with "bg:", like any other color.

If color names are used inside a linear-gradient, they will be [CSS Color Level 3 colors](https://www.w3.org/TR/2018/REC-css-color-3-20180619/#svg-color).  A word of caution; this can cause some unexpected behavior when mixing a base color name in a style and in a linear-gradient.  As mentioned above, when using a color name like "red" in a style, the style will use a 16-color ANSI color, and the exact color will depend on the terminal you are using and how it is configured.  Since kitsch-prompt has no way to know what the exact color "red" represents in your terminal, if you use the color "red" inside a linear-gradient, it will likely not show as the same color as "red" in a style:

```yaml
colors:
  $red1: linear-gradient(red, red) # <== This red is #ff0000
  $red2: red # <== This red is a 16-color ANSI color, and is
             # whatever the terminal says it is.
```

To get around this, use hex colors instead of the base color names to get the exact color you want.  Also note that "bright" colors like "brightRed" cannot be used in a linear-gradient.

If the user is using a terminal that only supports 256 colors, linear-gradients will be gracefully down sampled to the ANSI 256 color pallette.

### Custom Colors

Custom colors can be specified at the top of a configuration file.  Custom colors must start with a "$".

```yaml
colors:
    $foreground: "#10fa29"
    $background: "linear-gradient(#090, #010)"
```

Note that you should explicitly quote your hex colors, otherwise YAML will think they are comments.

These custom colors can be used anywhere in the configuration where you would normally use a color:

```yaml
modules:
  - type: username
    style: "$foreground bg:$background"
```

## Modifiers

The following are all valid modifiers. Note that some modifiers are not supported on some terminals:

- `Reset` - Resets the current color chain.
- `Bold` - Make text bold.
- `Dim` - Emitting only a small amount of light.
- `Italic` - Make text italic. _(Not widely supported)_
- `Underline` - Make text underline. _(Not widely supported)_
- `Inverse`- Inverse background and foreground colors.
- `Hidden` - Prints the text, but makes it invisible.
- `Strikethrough` - Puts a horizontal line through the center of the text. _(Not widely supported)_
- `Visible`- Prints the text only when gchalk has a color level > 0. Can be useful for things that are purely cosmetic.
