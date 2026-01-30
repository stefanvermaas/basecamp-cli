# Basecamp CLI

A command-line interface for Basecamp written in Go with zero external dependencies.

## Installation

```bash
go install github.com/rzolkos/basecamp-cli/cmd/basecamp@latest
```

Or build from source:

```bash
make build
```

## Setup

1. Create a Basecamp OAuth app at https://launchpad.37signals.com/integrations
2. Run `basecamp init` to configure credentials
3. Run `basecamp auth` to authenticate

Configuration files (XDG Base Directory):
- `~/.config/basecamp/config.json` - client credentials
- `~/.local/share/basecamp/token.json` - OAuth token

## Usage

```bash
# List all projects
basecamp projects

# List card tables in a project
basecamp boards <project_id>

# List cards in a board
basecamp cards <project_id> <board_id>

# Filter cards by column
basecamp cards <project_id> <board_id> --column "In Progress"

# View card details
basecamp card <project_id> <card_id>

# View card with comments
basecamp card <project_id> <card_id> --comments

# Move a card to a different column
basecamp move <project_id> <board_id> <card_id> --to "Done"
```

## Output

All commands output JSON for easy parsing with `jq`:

```bash
basecamp projects | jq '.[] | select(.status == "active") | .name'
```

Errors are output as JSON to stderr:

```json
{"error": "not authenticated, run 'basecamp auth' first"}
```

## Development

```bash
# Run tests
make test

# Build
make build

# Install to $GOPATH/bin
make install
```
