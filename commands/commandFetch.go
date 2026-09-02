package commands

import (
	"fmt"
	"strings"

	"github.com/Heavyymir/CharData_Aggregator/config"
	"github.com/Heavyymir/CharData_Aggregator/internal/models"
	"github.com/Heavyymir/CharData_Aggregator/internal/storage/sqlite"
	"github.com/Heavyymir/CharData_Aggregator/internal/parsers/bbcf"
	"github.com/Heavyymir/CharData_Aggregator/internal/parsers/ggst"
	"github.com/Heavyymir/CharData_Aggregator/internal/parsers/fat"
	//"github.com/Heavyymir/CharData_Aggregator/internal/parsers/sf3s"
)

// Command to Fetch character data once accessing a wiki
func commandFetch(cfg *config.Config, args ...string) error {

	// Ensure the entered strings aren't empty
	if cfg.Wiki.Name == "" || cfg.Game.Name == "" {
		return fmt.Errorf("select a wiki and game first")
	}

	if len(args) != 1 {
		return fmt.Errorf("usage: fetch <character>")
	}

	characterName := args[0]
	characterSlug := strings.ToLower(strings.ReplaceAll(characterName, " ", "_"))
	
	// Format page title from catalog URL patterns using characterName   
	pagePath := strings.Replace(cfg.Game.CharacterPath, "{character}", characterName, 1)
	requestURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.Wiki.URL, "/"), pagePath)

	data, err := cfg.CharDataClient.Fetch(requestURL)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", requestURL, err)
	}
	
	var moves []models.Move
	
	switch cfg.Game.Slug {
	case "bbcf":
		moves, err = bbcf.BBCFCharPageParser(data)
		if err != nil {
			return err
		}

	case "ggst":
		moves, err = ggst.GGSTCharPageParser(data)
		if err != nil {
			return err
		}

	case "sf6", "sf5", "usf4":
		moves, err = fat.FATJSONParser(data)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("no parser available for game: %s", cfg.Game.Name)
	}

	character := models.Character{
		Name:	characterName,
		Slug:	characterSlug,
		URL:	requestURL,
	}

	characterID, err := sqlite.SaveCharacter(
		cfg.DB,
		cfg.Game.Slug,
		character,
	)

	if err != nil {
		return err
	}

	for _, move := range moves {
		rowCount := 0
	
		for _, grid := range move.FrameDataGrids {
			rowCount += len(grid.Rows)
		}
	
	
		for gridIndex, grid := range move.FrameDataGrids {
			for rowIndex, row := range grid.Rows {
				fmt.Printf(
					"before save: grid=%d row=%d cells=%d\n",
					gridIndex,
					rowIndex,
					len(row.Cells),
				)
			}
		}
	}
    
	if err := sqlite.SaveMoves(cfg.DB, characterID, moves); err != nil {
		return err
	}
		
	printMoveTable(moves)

	// Return request data to the user so they can see it has been successful
	fmt.Printf("fetched %d bytes from %s\n", len(data), requestURL)

	return nil
}
