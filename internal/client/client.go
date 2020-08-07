package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

var userAgent = fmt.Sprintf("sntr go/%s", runtime.Version()[2:])

type Client struct {
	APIRoot    string
	AuthString string

	Debug bool
	JSON  bool
}

func (c *Client) Get(path string, ret interface{}) error {
	endpoint := c.EndpointFor(path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.AuthString)
	req.Header.Add("User-Agent", userAgent)

	if c.Debug {
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

	if c.Debug {
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

	if c.JSON {
		// Dump response body to stdout
		_, err = io.Copy(os.Stdout, resp.Body)
	} else {
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(ret)
	}
	return err
}

func (c *Client) GetSingle(path string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := c.Get(path, &ret)
	return ret, err
}

func (c *Client) GetMultiple(path string) ([]map[string]interface{}, error) {
	var ret []map[string]interface{}
	err := c.Get(path, &ret)
	return ret, err
}

func (c *Client) EndpointFor(path string) string {
	var b strings.Builder
	b.WriteString(c.APIRoot)
	if !strings.HasPrefix(path, "/") {
		b.WriteByte('/')
	}
	b.WriteString(path)
	if u, err := url.Parse(b.String()); err == nil && u.RawQuery != "" {
		return b.String()
	}
	if path != "" && !strings.HasSuffix(path, "/") {
		b.WriteByte('/')
	}
	return b.String()
}
