package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableCommandSorting = false
}

type UsageError struct {
	error
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sntr",
		Short:        "Sentry command-line tool",
		Long:         `The sntr tool gives you quick access to your data in Sentry.`,
		SilenceUsage: true,
	}
	cmd.AddCommand(
		NewGetCommand(),
		NewSendCommand(),
		NewProxyCommand(),
		NewExecCommand(),
	)
	return cmd
}

// Execute executes the root command.
func Execute() error {
	return NewRootCommand().Execute()
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
