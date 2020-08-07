package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

var completionTargets = []string{"bash", "zsh", "fish", "powershell"}

func NewCompletionCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("completion [%s]", strings.Join(completionTargets, "|")),
		Short: "Generate completion script",
		Long: `To load completions:

Bash/Zsh:

$ source <(sntr completion bash)
$ source <(sntr completion zsh) && compdef _sntr sntr`,
		DisableFlagsInUseLine: true,
		ValidArgs:             completionTargets,
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			out := root.OutOrStdout()
			switch args[0] {
			case "bash":
				root.GenBashCompletion(out)
			case "zsh":
				root.GenZshCompletion(out)
			case "fish":
				root.GenFishCompletion(out, true)
			case "powershell":
				root.GenPowerShellCompletion(out)
			}
		},
	}
	return cmd
}
