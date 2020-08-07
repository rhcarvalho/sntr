package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

type SendCommand struct {
	cfg *config.Config
}

func NewSendCommand(cfg *config.Config) *cobra.Command {
	c := &SendCommand{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:   "send ORGANIZATION_SLUG PROJECT_SLUG",
		Short: "Send events",
		Long:  `Send events to Sentry.`,
		Args:  cobra.ExactArgs(2),
		RunE:  c.Run,
	}
	return cmd
}

func (c *SendCommand) Run(cmd *cobra.Command, args []string) error {
	if c.cfg.AuthToken == "" {
		return errors.New(`missing authentication token: run "sntr login" to setup`)
	}
	orgSlug, projSlug := args[0], args[1]

	s, err := c.cfg.Client.GetMultiple(fmt.Sprintf("projects/%s/%s/keys", orgSlug, projSlug))
	if err != nil {
		return err
	}
	m := s[0]

	dsn := func() string {
		defer func() { recover() }()
		return m["dsn"].(map[string]interface{})["public"].(string)
	}()
	if dsn == "" {
		return fmt.Errorf("unknown DSN: %#v", m)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
	if err != nil {
		return err
	}

	start := time.Now()
	id := sentry.CaptureMessage("test message from sntr")
	if id == nil {
		return errors.New("internal error: message was not sent :(")
	}
	log.Println("event sent:", *id)
	log.Println("waiting to read it back")
	defer func() {
		log.Println("took", time.Since(start))
	}()

	return c.WaitEvent(orgSlug, string(*id))
}

func (c *SendCommand) WaitEvent(orgSlug, id string) error {
	var err error
	for i := 0; i < 5; i++ {
		err = GetOrganizationEvent(c.cfg, orgSlug, id)
		if err == nil || !strings.Contains(err.Error(), "404") {
			break
		}
		time.Sleep(1 << i * time.Second)
	}
	return err
}
