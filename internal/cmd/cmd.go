package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sntr/internal/config"
)

var (
	flagDebug = flag.Bool("debug", false, "Write debug messages to stderr")
	flagJSON  = flag.Bool("json", false, "Set output format to JSON")
)

var userAgent = fmt.Sprintf("sntr go/%s", runtime.Version()[2:])

func getMultiple(cfg *config.Config, path string) ([]map[string]interface{}, error) {
	var ret []map[string]interface{}
	err := get(cfg, path, &ret)
	return ret, err
}

func getSingle(cfg *config.Config, path string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := get(cfg, path, &ret)
	return ret, err
}

func get(cfg *config.Config, path string, ret interface{}) error {
	endpoint := cfg.EndpointFor(path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", cfg.AuthString)
	req.Header.Add("User-Agent", userAgent)

	if *flagDebug {
		b, err := httputil.DumpRequest(req, false)
		if err != nil {
			return err
		}
		os.Stderr.Write(b)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if *flagDebug {
		// Dump response minus body to stderr
		b, err := httputil.DumpResponse(resp, false)
		if err != nil {
			return err
		}
		os.Stderr.Write(b)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	if *flagJSON {
		// Dump response body to stdout
		_, err = io.Copy(os.Stdout, resp.Body)
	} else {
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(ret)
	}
	return err
}

func ListOrganizations(cfg *config.Config) error {
	orgs, err := getMultiple(cfg, "organizations")
	if err != nil {
		return err
	}
	for _, org := range orgs {
		fmt.Println(org["slug"])
	}
	return nil
}

func ListProjects(cfg *config.Config) error {
	projs, err := getMultiple(cfg, "projects")
	if err != nil {
		return err
	}
	for _, proj := range projs {
		org, ok := proj["organization"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("organization is not a JSON object: %#v", proj["organization"])
		}
		fmt.Printf("%s/%s\n", org["slug"], proj["slug"])
	}
	return nil
}

func ListOrganizationProjects(cfg *config.Config, slug string) error {
	// TODO: limit to projects isMember=true
	projs, err := getMultiple(cfg, fmt.Sprintf("organizations/%s/projects", slug))
	if err != nil {
		return err
	}
	for _, proj := range projs {
		fmt.Printf("%s/%s\n", slug, proj["slug"])
	}
	return nil
}

func ListProjectIssues(cfg *config.Config, orgSlug, projSlug string) error {
	issues, err := getMultiple(cfg, fmt.Sprintf("projects/%s/%s/issues", orgSlug, projSlug))
	if err != nil {
		return err
	}
	for _, issue := range issues {
		fmt.Printf("%s: %s\n", issue["shortId"], issue["title"])
	}
	return nil
}

func GetOrganizationEvent(cfg *config.Config, orgSlug, id string) error {
	m, err := getSingle(cfg, fmt.Sprintf("organizations/%s/eventids/%s", orgSlug, id))
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
	event, err := getSingle(cfg, fmt.Sprintf("projects/%s/%s/events/%s", orgSlug, projSlug, id))
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
