package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

func init() {
	cobra.EnableCommandSorting = false
}

type UsageError struct {
	error
}

func NewRootCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sntr",
		Short:        "Sentry command-line tool",
		Long:         `The sntr tool gives you quick access to your data in Sentry.`,
		SilenceUsage: true,
	}
	cmd.AddCommand(
		NewGetCommand(cfg),
		NewSendCommand(cfg),
		NewProxyCommand(cfg),
		NewExecCommand(cfg),
		NewLoginCommand(cfg),
		NewCompletionCommand(cfg),
	)
	return cmd
}

// Execute executes the root command.
func Execute(cfg *config.Config) error {
	return NewRootCommand(cfg).Execute()
}

func notImplemented(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("command %q is not implemented yet", cmd.Name())
}

type runFunc func(cmd *cobra.Command, args []string) error

func checkUsage(f runFunc) runFunc {
	return func(cmd *cobra.Command, args []string) error {
		err := f(cmd, args)
		if _, ok := err.(UsageError); ok {
			cmd.Usage()
		}
		return err
	}
}
