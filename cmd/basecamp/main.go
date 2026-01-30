package main

import (
	"os"

	"github.com/rzolkos/basecamp-cli/internal/commands"
)

var version = "dev"

func main() {
	commands.Execute(os.Args[1:], version)
}
