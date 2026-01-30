package commands

import (
	"encoding/json"
	"fmt"
	"os"
)

type Command interface {
	Run(args []string) error
}

var commands = map[string]func() Command{
	"init":     func() Command { return &InitCmd{} },
	"auth":     func() Command { return &AuthCmd{} },
	"projects": func() Command { return &ProjectsCmd{} },
	"boards":   func() Command { return &BoardsCmd{} },
	"cards":    func() Command { return &CardsCmd{} },
	"card":     func() Command { return &CardCmd{} },
	"move":     func() Command { return &MoveCmd{} },
}

func Execute(args []string, version string) {
	if len(args) < 1 {
		printHelp(version)
		os.Exit(1)
	}

	cmd := args[0]

	if cmd == "version" || cmd == "--version" || cmd == "-v" {
		fmt.Println(version)
		return
	}

	if cmd == "help" || cmd == "--help" || cmd == "-h" {
		printHelp(version)
		return
	}

	factory, ok := commands[cmd]
	if !ok {
		PrintError(fmt.Errorf("unknown command: %s", cmd))
		os.Exit(1)
	}

	if err := factory().Run(args[1:]); err != nil {
		PrintError(err)
		os.Exit(1)
	}
}

func printHelp(version string) {
	fmt.Printf(`basecamp - Basecamp CLI %s

Usage: basecamp <command> [arguments] [flags]

Commands:
  init                              Configure credentials
  auth                              Authenticate with OAuth
  projects                          List all projects
  boards <project_id>               List card tables in a project
  cards <project_id> <board_id>     List cards (--column <name> to filter)
  card <project_id> <card_id>       View card details (--comments for comments)
  move <proj> <board> <card>        Move card (--to <column> required)
  version                           Show version

Examples:
  basecamp projects
  basecamp boards 12345678
  basecamp cards 12345678 87654321 --column "In Progress"
  basecamp card 12345678 44444444 --comments
  basecamp move 12345678 87654321 44444444 --to "Done"
`, version)
}

func PrintError(err error) {
	errJSON, _ := json.Marshal(map[string]string{"error": err.Error()})
	fmt.Fprintln(os.Stderr, string(errJSON))
}

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
