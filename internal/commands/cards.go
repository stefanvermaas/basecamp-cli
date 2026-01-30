package commands

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type CardsCmd struct{}

type ColumnDetail struct {
	Title      string `json:"title"`
	CardsCount int    `json:"cards_count"`
	CardsURL   string `json:"cards_url"`
}

type CardTableDetail struct {
	ID    int            `json:"id"`
	Title string         `json:"title"`
	Lists []ColumnDetail `json:"lists"`
}

type Creator struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CardSummary struct {
	ID      int     `json:"id"`
	Title   string  `json:"title"`
	Creator Creator `json:"creator"`
}

type CardOutput struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Creator string `json:"creator"`
}

type ColumnCards struct {
	Column string       `json:"column"`
	Cards  []CardOutput `json:"cards"`
}

type CardsOutput struct {
	BoardID    int           `json:"board_id"`
	BoardTitle string        `json:"board_title"`
	Columns    []ColumnCards `json:"columns"`
}

func (c *CardsCmd) Run(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: basecamp cards <project_id> <board_id> [--column <name>]")
	}
	projectID := args[0]
	boardID := args[1]

	var columnFilter string
	for i := 2; i < len(args); i++ {
		if args[i] == "--column" && i+1 < len(args) {
			columnFilter = args[i+1]
			break
		}
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Get the card table
	data, err := cl.Get("/buckets/" + projectID + "/card_tables/" + boardID + ".json")
	if err != nil {
		return err
	}

	var cardTable CardTableDetail
	if err := json.Unmarshal(data, &cardTable); err != nil {
		return err
	}

	output := CardsOutput{
		BoardID:    cardTable.ID,
		BoardTitle: cardTable.Title,
		Columns:    []ColumnCards{},
	}

	for _, list := range cardTable.Lists {
		// Filter by column if specified
		if columnFilter != "" && !strings.Contains(strings.ToLower(list.Title), strings.ToLower(columnFilter)) {
			continue
		}

		if list.CardsCount == 0 {
			continue
		}

		// Fetch cards from this column
		cardsData, err := cl.Get(list.CardsURL)
		if err != nil {
			return err
		}

		var cards []CardSummary
		if err := json.Unmarshal(cardsData, &cards); err != nil {
			return err
		}

		columnCards := ColumnCards{
			Column: list.Title,
			Cards:  make([]CardOutput, len(cards)),
		}

		for i, card := range cards {
			creator := "Unknown"
			if card.Creator.Name != "" {
				creator = card.Creator.Name
			}
			columnCards.Cards[i] = CardOutput{
				ID:      card.ID,
				Title:   card.Title,
				Creator: creator,
			}
		}

		output.Columns = append(output.Columns, columnCards)
	}

	return PrintJSON(output)
}
