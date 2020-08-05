package main

import (
	"os"

	"github.com/getsentry/sntr/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
