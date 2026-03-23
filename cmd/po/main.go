package main

import (
	"os"

	"github.com/wangke19/po/internal/pocmd"
)

func main() {
	code := pocmd.Main()
	os.Exit(int(code))
}
