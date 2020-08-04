package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/getsentry/sntr/internal/cmd"
)

var organizationProjectsRegexp = regexp.MustCompile("^(?:org|organization)s?/([^/]+)/projects$")

func main() {
	flag.Parse()

	var err error
	switch arg := flag.Arg(0); arg {
	case "":
		err = fmt.Errorf("usage: sntr <command>")
	case "organizations", "orgs":
		err = cmd.ListOrganizations()
	case "projects":
		err = cmd.ListProjects()
	default:
		if m := organizationProjectsRegexp.FindStringSubmatch(arg); m != nil {
			err = cmd.ListOrganizationProjects(m[1])
		} else {
			err = fmt.Errorf("unknown command: %s", arg)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
