package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/getsentry/sntr/internal/cmd"
	"github.com/getsentry/sntr/internal/config"
)

func main() {
	cfg, err := config.LoadDefault()
	if err == nil && cfg.AuthToken == "" {
		err = errors.New(`configuration file missing "auth_token" field`)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := cmd.Execute(cfg); err != nil {
		os.Exit(1)
	}
}
