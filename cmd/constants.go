package cmd

// programName is the name of the kitsch executable.
const programName = "kitsch"

const website = "https://kitschprompt.com"
const githubRepo = "https://github.com/jwalton/kitsch"

var supportedShells = []string{"bash", "zsh"}

// These are set by goreleaser.
var version string = "unknown"
var commit string = "unknown"
