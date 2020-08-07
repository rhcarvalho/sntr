package cmd

import (
	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

func NewProxyCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Start an HTTP proxy",
		Long: `Start an HTTP proxy that intercepts and forwards events.
Typically used to intercept events from an SDK before sending them to Sentry.`,
		RunE: notImplemented,
	}
	return cmd
}
