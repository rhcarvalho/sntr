package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/getsentry/sntr/internal/config"
)

var (
	organizationProjectsRegexp             = regexp.MustCompile(`^(?:org|organization)s?/([^/]+)/projects$`)
	orgSlugSlashProjSlugRegexp             = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
	orgSlugSlashEventIDRegexp              = regexp.MustCompile(`^([^/]+)/([A-Fa-f0-9]{32})$`)
	orgSlugSlashProjSlugSlashEventIDRegexp = regexp.MustCompile(`^([^/]+)/([^/]+)/([A-Fa-f0-9]{32})$`)
)

type GetCommand struct {
	cfg   *config.Config
	query string
}

func NewGetCommand(cfg *config.Config) *cobra.Command {
	c := &GetCommand{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource",
		Long:  `Get resource.`,
		ValidArgs: []string{
			"organization", "organizations", "org", "orgs",
			"project", "projects", "proj", "projs",
		},
		RunE: checkUsage(c.Run),
	}
	cmd.Flags().StringVarP(&c.query, "query", "q", "", "Query to search events like in Discover")
	return cmd
}

func (c *GetCommand) Run(cmd *cobra.Command, args []string) error {
	if c.query != "" {
		return c.RunDiscover(cmd, args)
	}

	if len(args) == 0 {
		return UsageError{errors.New("missing resource type")}
	}
	var err error
	switch arg := args[0]; arg {
	case "organization", "organizations", "org", "orgs":
		err = ListOrganizations(c.cfg)
	case "project", "projects", "proj", "projs":
		err = ListProjects(c.cfg)
	default:
		if m := orgSlugSlashEventIDRegexp.FindStringSubmatch(arg); m != nil {
			err = GetOrganizationEvent(c.cfg, m[1], m[2])
			break
		}
		if m := orgSlugSlashProjSlugSlashEventIDRegexp.FindStringSubmatch(arg); m != nil {
			err = GetOrganizationProjectEvent(c.cfg, m[1], m[2], m[3])
			break
		}
		if m := organizationProjectsRegexp.FindStringSubmatch(arg); m != nil {
			err = ListOrganizationProjects(c.cfg, m[1])
			break
		}
		if m := orgSlugSlashProjSlugRegexp.FindStringSubmatch(arg); m != nil {
			err = ListProjectIssues(c.cfg, m[1], m[2])
			break
		}
		err = UsageError{fmt.Errorf("unknown command: %s", arg)}
	}
	return err
}

func (c *GetCommand) RunDiscover(cmd *cobra.Command, args []string) error {
	orgSlug := args[0]
	fields := "field=project&field=timestamp&field=title&sort=-timestamp"
	m, err := getSingle(c.cfg, fmt.Sprintf("organizations/%s/eventsv2/?query=%s&%s", orgSlug, url.QueryEscape(c.query), fields))
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 1, 8, 1, ' ', 0)
	fmt.Fprintln(tw, "ID\tPROJECT\tTIMESTAMP\tEVENT TITLE")
	for _, e := range m["data"].([]interface{}) {
		event := e.(map[string]interface{})
		id := event["id"].(string)
		project := event["project"].(string)
		timestamp := event["timestamp"].(string)
		title := event["title"].(string)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", id, project, timestamp, title)
	}
	return tw.Flush()
}
