package main

import (
	"fmt"
	"os"

	"github.com/getsentry/sntr/internal/cmd"
	"github.com/getsentry/sntr/internal/config"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	defer cfg.Save()
	return cmd.Execute(cfg)
}
