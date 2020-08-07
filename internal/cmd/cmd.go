package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/getsentry/sntr/internal/config"
)

var (
	flagDebug = flag.Bool("debug", false, "Write debug messages to stderr")
	flagJSON  = flag.Bool("json", false, "Set output format to JSON")
)

func ListOrganizations(cfg *config.Config) error {
	orgs, err := cfg.Client.GetMultiple("organizations")
	if err != nil {
		return err
	}
	var out []string
	for _, org := range orgs {
		out = append(out, fmt.Sprint(org["slug"]))
	}
	sort.StringSlice(out).Sort()
	fmt.Println(strings.Join(out, "\n"))
	return nil
}

func ListProjects(cfg *config.Config) error {
	projs, err := cfg.Client.GetMultiple("projects")
	if err != nil {
		return err
	}
	var out []string
	for _, proj := range projs {
		org, ok := proj["organization"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("organization is not a JSON object: %#v", proj["organization"])
		}
		out = append(out, fmt.Sprintf("%s/%s", org["slug"], proj["slug"]))
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.Replace(out[i], "/", "\x00", -1) < strings.Replace(out[j], "/", "\x00", -1)
	})
	fmt.Println(strings.Join(out, "\n"))
	return nil
}

func ListOrganizationProjects(cfg *config.Config, slug string) error {
	// TODO: limit to projects isMember=true
	projs, err := cfg.Client.GetMultiple(fmt.Sprintf("organizations/%s/projects", slug))
	if err != nil {
		return err
	}
	for _, proj := range projs {
		fmt.Printf("%s/%s\n", slug, proj["slug"])
	}
	return nil
}

func ListProjectIssues(cfg *config.Config, orgSlug, projSlug string) error {
	issues, err := cfg.Client.GetMultiple(fmt.Sprintf("projects/%s/%s/issues", orgSlug, projSlug))
	if err != nil {
		return err
	}
	for _, issue := range issues {
		fmt.Printf("%s: %s\n", issue["shortId"], issue["title"])
	}
	return nil
}

func GetOrganizationEvent(cfg *config.Config, orgSlug, id string) error {
	m, err := cfg.Client.GetSingle(fmt.Sprintf("organizations/%s/eventids/%s", orgSlug, id))
	if err != nil {
		return err
	}
	event, ok := m["event"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("event is not a JSON object: %#v", m["event"])
	}
	keys := make([]string, 0, len(event))
	for k := range event {
		keys = append(keys, k)
	}
	fmt.Printf("%s: %s\n", id, strings.Join(keys, ", "))
	return nil
}

func GetOrganizationProjectEvent(cfg *config.Config, orgSlug, projSlug, id string) error {
	event, err := cfg.Client.GetSingle(fmt.Sprintf("projects/%s/%s/events/%s", orgSlug, projSlug, id))
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(event))
	for k := range event {
		keys = append(keys, k)
	}
	fmt.Printf("%s: %s\n", id, strings.Join(keys, ", "))
	return nil
}
