package cmd

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	organizationProjectsRegexp             = regexp.MustCompile(`^(?:org|organization)s?/([^/]+)/projects$`)
	orgSlugSlashProjSlugRegexp             = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
	orgSlugSlashEventIDRegexp              = regexp.MustCompile(`^([^/]+)/([A-Fa-f0-9]{32})$`)
	orgSlugSlashProjSlugSlashEventIDRegexp = regexp.MustCompile(`^([^/]+)/([^/]+)/([A-Fa-f0-9]{32})$`)
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource",
		Long:  `Get resource.`,
		RunE:  checkUsage(runGet),
	}
	return cmd
}

func runGet(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return UsageError{errors.New("missing resource type")}
	}
	var err error
	switch arg := args[0]; arg {
	case "organizations", "orgs":
		err = ListOrganizations()
	case "projects":
		err = ListProjects()
	default:
		if m := orgSlugSlashEventIDRegexp.FindStringSubmatch(arg); m != nil {
			err = GetOrganizationEvent(m[1], m[2])
			break
		}
		if m := orgSlugSlashProjSlugSlashEventIDRegexp.FindStringSubmatch(arg); m != nil {
			err = GetOrganizationProjectEvent(m[1], m[2], m[3])
			break
		}
		if m := organizationProjectsRegexp.FindStringSubmatch(arg); m != nil {
			err = ListOrganizationProjects(m[1])
			break
		}
		if m := orgSlugSlashProjSlugRegexp.FindStringSubmatch(arg); m != nil {
			err = ListProjectIssues(m[1], m[2])
			break
		}
		err = UsageError{fmt.Errorf("unknown command: %s", arg)}
	}
	return err
}
