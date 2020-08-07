package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

type ProxyCommand struct {
	cfg *config.Config
}

func NewProxyCommand(cfg *config.Config) *cobra.Command {
	c := &ProxyCommand{cfg: cfg}
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Start an HTTP proxy",
		Long: `Start an HTTP proxy that intercepts and forwards events.
Typically used to intercept events from an SDK before sending them to Sentry.`,
		RunE: c.Run,
	}
	return cmd
}

func (c *ProxyCommand) Run(cmd *cobra.Command, args []string) error {
	url, err := url.Parse(c.cfg.SentryURL)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	origDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		log.Println("Got request", r.Method, r.URL)
		if origDirector != nil {
			origDirector(r)
		}
	}
	origModifyResponse := proxy.ModifyResponse
	proxy.ModifyResponse = func(r *http.Response) error {
		log.Println("Got response", r.Status)
		if origModifyResponse != nil {
			return origModifyResponse(r)
		}
		return nil
	}
	const addr = "localhost:4000"
	log.Printf("Proxy to %s listening on http://%s", url, addr)
	return http.ListenAndServe(addr, proxy)
}
