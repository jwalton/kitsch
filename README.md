# Kitsch Prompt

Kitsch prompt is a cross-platform tool for displaying a shell prompt, which can be extensively customized both in terms of what is shown, and in terms of how it is shown.  Kitsch prompt makes it easy to render your prompt with

## Installation

On Linux or Mac, you can install kitsch by running:

```sh
$ curl https://raw.githubusercontent.com/jwalton/kitsch/master/install.sh | sh
```

If you have go development toolchain installed, you can install kitsch by cloning this repo and running:

```sh
$ make install
```

## Setup

### ZSH

Add the following to your ~/.zshrc:

```sh
if which kitsch > /dev/null; then
    eval "$(kitsch init zsh)"
fi
```
