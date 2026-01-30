package tests

import (
	"testing"

	"github.com/rzolkos/basecamp-cli/e2e/harness"
)

func TestVersion(t *testing.T) {
	h := harness.New(t)

	result := h.Run("version")

	if !result.Success() {
		t.Errorf("expected success, got exit code %d", result.ExitCode)
	}

	if result.Stdout == "" {
		t.Error("expected version output")
	}
}

func TestHelp(t *testing.T) {
	h := harness.New(t)

	result := h.Run("help")

	if !result.Success() {
		t.Errorf("expected success, got exit code %d", result.ExitCode)
	}

	if result.Stdout == "" {
		t.Error("expected help output")
	}
}

func TestProjects(t *testing.T) {
	h := harness.New(t)

	result := h.Run("projects")

	if !result.Success() {
		t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
	}

	if result.JSONArray == nil {
		t.Fatalf("expected JSON array, got: %s", result.Stdout)
	}

	if len(result.JSONArray) == 0 {
		t.Error("expected at least one project")
	}

	// Check first project has required fields
	first := result.JSONArray[0]
	if first["id"] == nil {
		t.Error("project missing 'id' field")
	}
	if first["name"] == nil {
		t.Error("project missing 'name' field")
	}
	if first["status"] == nil {
		t.Error("project missing 'status' field")
	}
}

func TestBoards(t *testing.T) {
	h := harness.New(t)

	t.Run("with explicit project_id", func(t *testing.T) {
		result := h.Run("boards", h.ProjectID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("project_id") == 0 {
			t.Error("expected project_id in response")
		}
		if result.GetString("project_name") == "" {
			t.Error("expected project_name in response")
		}
		if result.GetInt("board_id") == 0 {
			t.Error("expected board_id in response")
		}
	})
}

func TestCards(t *testing.T) {
	h := harness.New(t)

	if h.BoardID == "" {
		t.Skip("BASECAMP_TEST_BOARD_ID not set")
	}

	t.Run("list all cards", func(t *testing.T) {
		result := h.Run("cards", h.ProjectID, h.BoardID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("board_id") == 0 {
			t.Error("expected board_id in response")
		}
		if result.GetString("board_title") == "" {
			t.Error("expected board_title in response")
		}

		columns := result.GetNested("columns")
		if columns == nil {
			t.Error("expected columns in response")
		}
	})

	t.Run("filter by column", func(t *testing.T) {
		// Get a column name first
		result := h.Run("boards", h.ProjectID)
		if !result.Success() {
			t.Skip("could not get board info")
		}

		columns, ok := result.GetNested("columns").([]any)
		if !ok || len(columns) == 0 {
			t.Skip("no columns found")
		}

		firstCol, ok := columns[0].(map[string]any)
		if !ok {
			t.Skip("could not parse column")
		}
		colName, ok := firstCol["title"].(string)
		if !ok {
			t.Skip("could not get column name")
		}

		// Filter by that column
		result = h.Run("cards", h.ProjectID, h.BoardID, "--column", colName)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}
	})
}

func TestCard(t *testing.T) {
	h := harness.New(t)

	if h.CardID == "" {
		t.Skip("BASECAMP_TEST_CARD_ID not set")
	}

	t.Run("view card details", func(t *testing.T) {
		result := h.Run("card", h.ProjectID, h.CardID)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetInt("id") == 0 {
			t.Error("expected id in response")
		}
		if result.GetString("title") == "" {
			t.Error("expected title in response")
		}
		if result.GetString("creator") == "" {
			t.Error("expected creator in response")
		}
	})

	t.Run("view card with comments", func(t *testing.T) {
		result := h.Run("card", h.ProjectID, h.CardID, "--comments")

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		// Comments field should exist (may be empty array or have items)
		// The field is only present when --comments is used and there are comments
	})
}

func TestMoveCard(t *testing.T) {
	h := harness.New(t)

	if h.BoardID == "" || h.CardID == "" {
		t.Skip("BASECAMP_TEST_BOARD_ID and BASECAMP_TEST_CARD_ID required for move test")
	}

	// Get available columns
	result := h.Run("boards", h.ProjectID)
	if !result.Success() {
		t.Fatalf("could not get board info: %s", result.Stderr)
	}

	columns, ok := result.GetNested("columns").([]any)
	if !ok || len(columns) < 2 {
		t.Skip("need at least 2 columns to test move")
	}

	// Get first two column names
	col1, _ := columns[0].(map[string]any)
	col2, _ := columns[1].(map[string]any)
	colName1, _ := col1["title"].(string)
	colName2, _ := col2["title"].(string)

	if colName1 == "" || colName2 == "" {
		t.Skip("could not get column names")
	}

	// Find which column the card is currently in
	cardResult := h.Run("cards", h.ProjectID, h.BoardID)
	if !cardResult.Success() {
		t.Fatalf("could not get cards: %s", cardResult.Stderr)
	}

	// Determine target column (pick one the card isn't in, or just use col2)
	targetCol := colName2
	returnCol := colName1

	t.Run("move card to different column", func(t *testing.T) {
		result := h.Run("move", h.ProjectID, h.BoardID, h.CardID, "--to", targetCol)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.JSON == nil {
			t.Fatalf("expected JSON object, got: %s", result.Stdout)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})

	t.Run("move card back", func(t *testing.T) {
		result := h.Run("move", h.ProjectID, h.BoardID, h.CardID, "--to", returnCol)

		if !result.Success() {
			t.Fatalf("expected success, got exit code %d\nstderr: %s", result.ExitCode, result.Stderr)
		}

		if result.GetString("status") != "ok" {
			t.Error("expected status=ok")
		}
	})
}

func TestInvalidCommand(t *testing.T) {
	h := harness.New(t)

	result := h.Run("notacommand")

	if result.Success() {
		t.Error("expected failure for invalid command")
	}

	if result.ErrorMessage() == "" {
		t.Error("expected error message")
	}
}

func TestMissingArguments(t *testing.T) {
	h := harness.New(t)

	t.Run("cards without board_id", func(t *testing.T) {
		result := h.Run("cards", h.ProjectID)

		if result.Success() {
			t.Error("expected failure when board_id missing")
		}
	})

	t.Run("move without --to", func(t *testing.T) {
		if h.BoardID == "" || h.CardID == "" {
			t.Skip("need board and card IDs")
		}

		result := h.Run("move", h.ProjectID, h.BoardID, h.CardID)

		if result.Success() {
			t.Error("expected failure when --to missing")
		}
	})
}
