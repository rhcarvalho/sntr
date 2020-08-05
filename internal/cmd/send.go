package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send ORGANIZATION_SLUG PROJECT_SLUG",
		Short: "Send events",
		Long:  `Send events to Sentry.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runSend,
	}
	return cmd
}

func runSend(cmd *cobra.Command, args []string) error {
	orgSlug, projSlug := args[0], args[1]

	s, err := getMultiple(fmt.Sprintf("projects/%s/%s/keys", orgSlug, projSlug))
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

	return waitEvent(orgSlug, string(*id))
}

func waitEvent(orgSlug, id string) error {
	var err error
	for i := 0; i < 5; i++ {
		err = GetOrganizationEvent(orgSlug, id)
		if err == nil || !strings.Contains(err.Error(), "404") {
			break
		}
		time.Sleep(1 << i * time.Second)
	}
	return err
}
