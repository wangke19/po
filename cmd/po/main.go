// Package main is the entry point for the po CLI.
package main

import (
	"os"

	"github.com/wangke19/po/internal/pocmd"
)

func main() {
	code := pocmd.Main()
	os.Exit(int(code))
}
