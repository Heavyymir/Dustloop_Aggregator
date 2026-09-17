package commands

import(
	"fmt"
	"strings"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/utils"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/display"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
)

// Cli Command to pull framedata for a character 
func commandFrames(cfg *config.Config, args ...string) error {
	// Verify arguments
	if len(args) < 1 {
		return fmt.Errorf("usage frames <character_name> [--details]")
	}

	if cfg.Game.Name == "" || cfg.Wiki.Slug == ""{
			return fmt.Errorf("select a wiki and game first")
	}
	
	var characterName []string
	showDetails := false

	// Seperate flags for character name words
	for _, arg := range args {
		if arg == "--details" || arg == "-d" {
			showDetails = true
		} else {
			characterName = append(characterName, arg)
		}
	}

	if len(characterName) == 0 {
		return fmt.Errorf("please provide a character name")
	}

	rawName := strings.Join(characterName, " ")
	characterSlug := utils.FormatWikiSlug(rawName)

	// Get character ID from SQL table
	characterID, err := sqlite.GetCharacterID(
		cfg.DB,
		cfg.Game.Slug,
		characterSlug,
	)
	if err != nil {
		return err
	}

	// Get moves from SQL table, print framedata to console
	moves, err := sqlite.GetMoves(cfg.DB, characterID)
	if err != nil {
		return err
	}

	// Debug/Output
	fmt.Printf("moves loaded: %d\n", len(moves))


	// Print logic for SQL framedata grids
	for _, move := range moves {
		//  Print the move header
		var title string

		name := strings.TrimSpace(move.Name)
		input := strings.TrimSpace(move.Input)

		if input == "" || strings.EqualFold(input, name) {
			title = name
		} else if name == "" {
			title = input 
		} else {
			title = fmt.Sprintf("%s (%s)", name, input)
			}

		fmt.Printf("=== %s  ===\n", title)
			
		//  Print the frame data grids
		for i, grid := range move.FrameDataGrids {
			if i == 0 {
				fmt.Println("[Base Frame Data]")
			} else {
				fmt.Printf("[Additional Data - Grid %d]\n", i)
			}
			display.PrintGrid(grid)
			fmt.Println()
		}
			
		//  Print notes if they exist
		if len(move.Notes) > 0 {
			fmt.Println("Notes:")
			for _, note := range move.Notes {
				fmt.Printf("  • %s\n", note)
			}
			fmt.Println()
		}
	
		// Print description when flag is enabled AND description exists
		if showDetails && move.Description != "" {
			fmt.Println("Description:")
			fmt.Println(move.Description)
			fmt.Println()
		}
	
		fmt.Println(strings.Repeat("=", 60))
		}

	return nil
	
}
