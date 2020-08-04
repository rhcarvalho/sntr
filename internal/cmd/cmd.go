package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	flagDebug = flag.Bool("debug", false, "Write debug messages to stderr")
	flagJSON  = flag.Bool("json", false, "Set output format to JSON")
)

const apiRoot = "https://sentry.io/api/0"

var userAgent = fmt.Sprintf("sntr go/%s", runtime.Version()[2:])

const configPath = "config.json"

var auth string

func init() {
	err := loadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("missing required configuration file %s", configPath)
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig() error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(f)
	var cfg map[string]string
	err = dec.Decode(&cfg)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("configuration file %s is empty", configPath)
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			err = fmt.Errorf("configuration file %s is invalid: %w", configPath, err)
		}
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			err = fmt.Errorf("configuration file %s is invalid: at offset %d: got %s, want object", configPath, typeErr.Offset, typeErr.Value)
		}
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			err = fmt.Errorf("configuration file %s is invalid: at offset %d: %w", configPath, syntaxErr.Offset, syntaxErr)
		}
		return err
	}
	token := cfg["token"]
	if token == "" {
		return errors.New(`configuration file missing "token" field`)
	}
	auth = fmt.Sprintf("Bearer %s", cfg["token"])
	return nil
}

func endpointFor(path string) string {
	var b strings.Builder
	b.WriteString(apiRoot)
	if !strings.HasPrefix(path, "/") {
		b.WriteByte('/')
	}
	b.WriteString(path)
	if !strings.HasSuffix(path, "/") {
		b.WriteByte('/')
	}
	return b.String()
}

func getMultiple(path string) ([]map[string]interface{}, error) {
	var ret []map[string]interface{}
	err := get(path, &ret)
	return ret, err
}

func getSingle(path string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := get(path, &ret)
	return ret, err
}

func get(path string, ret interface{}) error {
	endpoint := endpointFor(path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", auth)
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

func ListOrganizations() error {
	orgs, err := getMultiple("organizations")
	if err != nil {
		return err
	}
	for _, org := range orgs {
		fmt.Println(org["slug"])
	}
	return nil
}

func ListProjects() error {
	projs, err := getMultiple("projects")
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

func ListOrganizationProjects(slug string) error {
	// TODO: limit to projects isMember=true
	projs, err := getMultiple(fmt.Sprintf("organizations/%s/projects", slug))
	if err != nil {
		return err
	}
	for _, proj := range projs {
		fmt.Printf("%s/%s\n", slug, proj["slug"])
	}
	return nil
}

func ListProjectIssues(orgSlug, projSlug string) error {
	issues, err := getMultiple(fmt.Sprintf("projects/%s/%s/issues", orgSlug, projSlug))
	if err != nil {
		return err
	}
	for _, issue := range issues {
		fmt.Printf("%s: %s\n", issue["shortId"], issue["title"])
	}
	return nil
}

func GetOrganizationEvent(orgSlug, id string) error {
	m, err := getSingle(fmt.Sprintf("organizations/%s/eventids/%s", orgSlug, id))
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

func GetOrganizationProjectEvent(orgSlug, projSlug, id string) error {
	event, err := getSingle(fmt.Sprintf("projects/%s/%s/events/%s", orgSlug, projSlug, id))
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
