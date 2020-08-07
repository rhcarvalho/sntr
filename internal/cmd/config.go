package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

func NewConfigCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `Manage configuration.`,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print path to config file",
		Long:  `Print the full path to the configuration file.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), cfg.Path)
			return nil
		},
	})
	return cmd
}
