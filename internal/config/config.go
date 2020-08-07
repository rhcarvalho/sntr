package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/getsentry/sntr/internal/client"
)

const defaultSentryURL = "https://sentry.io"

type Config struct {
	// AuthToken must be provided, see https://sentry.io/settings/account/api/auth-tokens/.
	AuthToken string `json:"auth_token"`
	// SentryURL is https://sentry.io or an alternative target for API calls.
	SentryURL string `json:"sentry_url"`

	// The next fields are internal and managed automatically.

	User          string          `json:"user"`
	Organizations []*Organization `json:"organizations"`

	ActiveOrganization string `json:"active_organization"`

	Client     *client.Client `json:"-"`
	Path       string         `json:"-"`
	APIRoot    string         `json:"-"`
	AuthString string         `json:"-"`
}

func LoadDefault() (*Config, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(root, "sntr", "config.json")
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Path: path}, nil
		}
		return nil, err
	}
	dec := json.NewDecoder(f)
	var cfg *Config
	err = dec.Decode(&cfg)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("configuration file %s is empty", path)
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			err = fmt.Errorf("configuration file %s is invalid: %w", path, err)
		}
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			err = fmt.Errorf("configuration file %s is invalid: at offset %d: got %s, want object", path, typeErr.Offset, typeErr.Value)
		}
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			err = fmt.Errorf("configuration file %s is invalid: at offset %d: %w", path, syntaxErr.Offset, syntaxErr)
		}
		return nil, err
	}
	cfg.Path = path
	if cfg.SentryURL == "" {
		cfg.SentryURL = defaultSentryURL
	}
	cfg.APIRoot = cfg.SentryURL + "/api/0"
	cfg.Client = &client.Client{APIRoot: cfg.APIRoot}
	cfg.SetAuthToken(cfg.AuthToken)
	return cfg, nil
}

func (c *Config) Save() error {
	dir := filepath.Dir(c.Path)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(c)
}

func (c *Config) SetAuthToken(token string) {
	c.AuthToken = token
	c.Client.AuthString = fmt.Sprintf("Bearer %s", token)
}

func (c *Config) ObfuscatedAuthToken() string {
	token := c.AuthToken
	if len(token) < 16 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}

func (c *Config) VerifyToken() error {
	m, err := c.Client.GetSingle("")
	if err != nil {
		return err
	}
	c.User = m["user"].(map[string]interface{})["email"].(string)

	var requiredScopes = map[string]bool{
		"org:read":     true,
		"project:read": true,
		"event:read":   true,
	}
	for _, scope := range m["auth"].(map[string]interface{})["scopes"].([]interface{}) {
		delete(requiredScopes, scope.(string))
	}
	if len(requiredScopes) > 0 {
		var missing []string
		for scope := range requiredScopes {
			missing = append(missing, scope)
		}
		fmt.Fprintf(os.Stderr, "Warning: authentication token %s missing required permissions: %s", c.ObfuscatedAuthToken(), strings.Join(missing, ", "))
	}
	return nil
}

type Organization struct {
	Slug     string     `json:"slug"`
	Projects []*Project `json:"projects"`
}

type Project struct {
	Slug string `json:"slug"`
}
