package cmd

import (
	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

func NewExecCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute a new process and monitor errors",
		Long:  `Execute a new process and monitor errors or crashes.`,
		RunE:  notImplemented,
	}
	return cmd
}
