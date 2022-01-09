package initscripts

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/gchalk"
)

var text = gchalk.BrightCyan
var code = gchalk.WithBold().BrightWhite

var shellConfigFiles = map[string]string{
	"bash":       "~/.bashrc",
	"zsh":        "~/.zshrc",
	"powershell": "Microsoft.PowerShell_profile.ps1 (you can find the location of this file by running `echo $PROFILE`)",
}

// ShowSetupInstructions prints setup instructions for the given shell.
func ShowSetupInstructions(
	programName string,
	githubRepo string,
	website string,
	shell string,
	shellDetected bool,
) {
	shellSetupCommand := map[string]string{
		"bash":       `eval "$(` + programName + ` init bash)"`,
		"zsh":        `eval "$(` + programName + ` init zsh)"`,
		"powershell": `Invoke-Expression (&` + programName + ` init powershell)`,
	}

	shellLongSetupCommand := map[string]string{
		"bash": heredoc.Doc(`
			if command -v ` + programName + ` > /dev/null; then
			    eval "$(` + programName + ` init bash)"
			fi`),
		"zsh": heredoc.Doc(`
			if command -v ` + programName + ` > /dev/null; then
			    eval "$(` + programName + ` init zsh)"
			fi`),
		"powershell": `Invoke-Expression (&` + programName + ` init powershell)`,
	}

	_, shellSupported := shellConfigFiles[shell]

	fmt.Println()

	if !shellSupported {
		if !shellDetected {
			showShellNotDetected(githubRepo)
			os.Exit(1)
		}

		showShellNotSupported(programName)
		os.Exit(1)
	}

	fmt.Print(text("To try out " + programName + " in the current shell, run:\n\n"))
	fmt.Print(code("    " + shellSetupCommand[shell] + "\n\n"))

	fmt.Print(text(wordWrap(80, "To use "+programName+" by default for future shells, add "+
		"the following to the end of your "+shellConfigFiles[shell]+":\n\n")))
	fmt.Print(code(addIndent(shellLongSetupCommand[shell], "    ")))
	fmt.Print("\n\n")

	fmt.Print(text("To see these instructions again, run \"" + programName + " setup " + shell + "\".\n"))
	fmt.Print(text("To learn more about " + programName + ", visit " + website + ".\n\n"))
}

func showShellNotDetected(githubRepo string) {
	fmt.Print(text("Sorry, but it looks like your shell is unsupported.\n"))
	fmt.Print(text("At the moment, kitsch prompt supports the following shells:\n\n"))
	for supportedShell := range shellConfigFiles {
		fmt.Print(code("  " + supportedShell + "\n"))
	}
	fmt.Print(text("\nIf your shell is not supported, please raise an issue at\n"))
	fmt.Print(text(githubRepo + "/issues.\n\n"))
}

func showShellNotSupported(programName string) {
	fmt.Print(text("It looks like your shell is not supported or couldn't be detected.\n"))
	fmt.Print(text("Run " + code(programName+" setup [shell]") + " with your shell type to see how to setup\n"))
	fmt.Print(text(programName + ".\n\n"))
}

func addIndent(val string, indent string) string {
	parts := strings.Split(val, "\n")
	return indent + strings.Join(parts, "\n"+indent)
}

func wordWrap(width int, val string) string {
	out := strings.Builder{}

	spaceLeft := width
	lastWord := ""

	writeWord := func(word string) {
		if len(word) == 0 {
			return
		}

		if spaceLeft < (len(word) + 1) {
			out.WriteString("\n")
			spaceLeft = width
		} else if spaceLeft != width {
			out.WriteString(" ")
			spaceLeft--
		}
		out.WriteString(word)
		spaceLeft -= len(word)
	}

	for i := 0; i < len(val); i++ {
		if val[i] == ' ' {
			if lastWord == "" {
				out.WriteString(" ")
				spaceLeft--
			} else {
				writeWord(lastWord)
				lastWord = ""
			}
		} else {
			lastWord = lastWord + string(val[i])
		}
	}
	writeWord(lastWord)

	return out.String()
}
