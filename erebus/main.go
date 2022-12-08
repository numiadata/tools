package main

import (
	"os"

	"github.com/numiadata/tools/erebus/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
