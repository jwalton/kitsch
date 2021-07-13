# kitch-prompt

Kitch-prompt is a tool for displaying a shell prompt, with a focus on extreme customization.

## Installation

kitch-prompt supports the following shell types:

Via `go`:

```sh
$ go get https://github.com/jwalton/kitch-prompt
```

## Setup

For ZSH, add the following to your ~/.zshrc:

```
if which kitsch-prompt > /dev/null; then
    eval "$(kitsch-prompt init zsh)"
fi
```