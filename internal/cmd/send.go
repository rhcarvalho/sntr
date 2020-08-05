package cmd

import (
	"github.com/spf13/cobra"
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send events",
		Long:  `Send events to Sentry.`,
		RunE:  notImplemented,
	}
	return cmd
}
