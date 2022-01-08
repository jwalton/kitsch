package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Prints instructions about how to setup " + programName + ".",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var shell string

		if len(args) > 0 {
			shell = args[0]
		} else {
			shellType := os.Getenv("SHELL")
			if strings.HasSuffix(shellType, "/zsh") {
				shell = "zsh"
			} else if strings.HasSuffix(shellType, "/bash") {
				shell = "bash"
			} else {
				shell = "unknown"
			}
		}

		shellSupported := false
		for _, supportedShell := range supportedShells {
			if shell == supportedShell {
				shellSupported = true
			}
		}

		text := gchalk.BrightCyan
		code := gchalk.WithBold().BrightWhite

		fmt.Println()

		if !shellSupported {
			if len(args) == 1 {
				fmt.Print(text("Sorry, but it looks like your shell is unsupported.\n"))
				fmt.Print(text("At the moment, kitsch prompt supports the following shells:\n\n"))
				for _, supportedShell := range supportedShells {
					fmt.Print(code("  " + supportedShell + "\n"))
				}
				fmt.Print(text("\nIf your shell is not supported, please raise an issue at\n"))
				fmt.Print(text(githubRepo + "/issues.\n\n"))
				os.Exit(1)
			}

			fmt.Print(text("It looks like your shell is not supported or couldn't be detected.\n"))
			fmt.Print(text("Run " + code(programName+" setup [shell]") + " with your shell type to see how to setup\n"))
			fmt.Print(text(programName + ".\n\n"))
			os.Exit(1)
		}

		shellConfigFiles := map[string]string{
			"bash": "~/.bashrc",
			"zsh":  "~/.zshrc",
		}

		shellSetupCommand := map[string]string{
			"bash": `eval "$(` + programName + ` init bash)"`,
			"zsh":  `eval "$(` + programName + ` init zsh)"`,
		}

		fmt.Print(text("To try out " + programName + " in the current shell, run:\n\n"))
		fmt.Print(code("    " + shellSetupCommand[shell] + "\n\n"))

		fmt.Print(text("To use " + programName + " by default for future shells, add the\nfollowing to the end of your " + shellConfigFiles[shell] + ":\n\n"))
		fmt.Print(code(
			addIndent(heredoc.Doc(`
				if command -v `+programName+` > /dev/null; then
				    `+shellSetupCommand[shell]+`
				fi`),
				"    ",
			),
		))
		fmt.Print("\n\n")

		fmt.Print(text("To see these instructions again, run \"" + programName + " setup " + shell + "\".\n"))
		fmt.Print(text("To learn more about " + programName + ", visit " + website + ".\n\n"))

		// TODO: Add link to documentation here.
	},
}

func addIndent(val string, indent string) string {
	parts := strings.Split(val, "\n")
	return indent + strings.Join(parts, "\n"+indent)
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
