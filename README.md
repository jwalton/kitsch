# kitch-prompt

Kitch-prompt is a cross-platform tool for displaying a shell prompt, which can be extensively customized both in terms of what is shown, and in terms of how it is shown.  Kitch-prompt makes it easy to render your prompt with

## Installation

If you have go development toolchain installed, you can install kitsch-prompt by cloning this repo and running:

```sh
$ make install
```

## Setup

### ZSH

Add the following to your ~/.zshrc:

```sh
if which kitsch-prompt > /dev/null; then
    eval "$(kitsch-prompt init zsh)"
fi
```
