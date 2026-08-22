// Package main is the build entry point of revive.
package main

import (
	"os"

	"github.com/mgechev/revive/cli"
)

func main() {
	os.Exit(cli.RunRevive())
}
