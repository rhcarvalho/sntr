package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

type LoginCommand struct {
	cfg   *config.Config
	force bool
}

func NewLoginCommand(cfg *config.Config) *cobra.Command {
	c := &LoginCommand{cfg: cfg}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Sentry",
		Long: `Login to Sentry with an authentication token.

You can manage tokens at https://sentry.io/settings/account/api/auth-tokens/.`,
		RunE: c.Run,
	}
	cmd.Flags().BoolVarP(&c.force, "force", "f", false, "force a new login")
	return cmd
}

func (c *LoginCommand) Run(cmd *cobra.Command, args []string) error {
	if c.cfg.AuthToken != "" {
		c.printInfo(cmd.OutOrStdout(), true)
		if !c.force {
			return nil
		}
	}
	fmt.Fprintln(cmd.OutOrStdout(), "You'll need an authentication token with read permissions.")
	fmt.Fprintln(cmd.OutOrStdout(), "Tokens are managed in https://sentry.io/settings/account/api/auth-tokens/.")
	var token string
	for token == "" {
		fmt.Fprint(cmd.OutOrStdout(), "Paste your token here: ")
		if _, err := fmt.Fscanln(cmd.InOrStdin(), &token); errors.Is(err, io.EOF) {
			return nil
		}
	}
	c.cfg.SetAuthToken(token)
	if err := c.cfg.VerifyToken(); err != nil {
		c.cfg.SetAuthToken("")
		return fmt.Errorf("token %s could not be verified: %w", c.cfg.ObfuscatedAuthToken(), err)
	}
	c.printInfo(cmd.OutOrStdout(), false)
	return nil
}

func (c *LoginCommand) printInfo(w io.Writer, verify bool) {
	var b strings.Builder
	fmt.Fprintf(&b, "Logged in with token %s", c.cfg.ObfuscatedAuthToken())
	if verify && c.cfg.User == "" {
		c.cfg.VerifyToken()
	}
	if c.cfg.User != "" {
		fmt.Fprintf(&b, " (%s)", c.cfg.User)
	} else {
		fmt.Fprint(&b, " (not verified)")
	}
	fmt.Fprintln(w, b.String())
}
