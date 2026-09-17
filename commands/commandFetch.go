package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/utils"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/display"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/parsers/dustloop"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/parsers/mizuumi"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/discovery"
)

// Command to Fetch character data once accessing a wiki
func commandFetch(cfg *config.Config, args ...string) error {

	// Ensure the entered strings aren't empty
	if cfg.Wiki.Name == "" || cfg.Game.Name == "" {
		return fmt.Errorf("select a wiki and game first")
	}

	// Seperate character names from the verbose flag
	var charWords []string
	verbose := false
	showDetails := false
	force := false

	for _, arg := range args {
		trimmed := strings.TrimSpace(arg)
		switch trimmed {
		case "-v", "--verbose":
			verbose = true
		case "-d", "--details":
			showDetails = true
			verbose = true
		case "-f", "--force":
			force = true
		default:
			if trimmed != "" {
				charWords = append(charWords, trimmed)
			}
		}
	}
	
	if len(charWords) == 0 {
		return fmt.Errorf("usage: fetch <character> [-v|--verbose] [-d|--details]")
	}

	characterInput := strings.Join(charWords, " ")
	targetSlug := characterInput

	// Try loading discovered characters to get the exact casing from saved <wiki>_characters.json
	filename := strings.ToLower(cfg.Game.Slug) + "_characters.json"
	cache, err := discovery.LoadCharCache(filename)
	if err == nil {
		for _, char := range cache.Characters {
			if strings.EqualFold(char.Name, characterInput) || strings.EqualFold(char.Slug, characterInput) {
				targetSlug = char.Slug
				characterInput = char.Name
				break
			}
		}
	}
	
	// Fallback if cache wasn't found: Capitalize first letter for MediaWiki
	if targetSlug == characterInput && len(characterInput) > 0 {
		// Convert slug into url-friendly varient
		targetSlug = utils.FormatWikiSlug(characterInput)
	}

	// Build the URL using the matched slug
	pagePath := strings.Replace(cfg.Game.CharacterPath, "{character}", targetSlug, 1)
	requestURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.Wiki.URL, "/"), pagePath)

	// Check if character exists in local SQLite DB
	existingID, err := sqlite.GetCharacterID(cfg.DB, cfg.Game.Slug, targetSlug)
	if err == nil && existingID > 0 && !force {
		msg := fmt.Sprintf("'%s' is already in your database. Refetch and overwrite?", characterInput)

		// Use readline from config. If the answer to the above prompt is no, gracefully end function
		if !confirmPrompt(cfg.RL, msg) {
			fmt.Println("Fetch cancelled.")
			return nil
		}
		
	}
	
	spinnerMsg := fmt.Sprintf("Fetching data for %s from %s...", characterInput, cfg.Wiki.Name)
	sp := display.StartSpinner(spinnerMsg)


	var data []byte
	
	switch cfg.Wiki.Slug {
	case "mizuumi", "supercombo":
    	// Uses chromedp to wait out Cloudflare / Anubis PoW
    	fmt.Println("Using headless fetch method for character page, please wait a moment while challenges are passed")
   		data, err = cfg.CharDataClient.FetchHTMLHeadless(requestURL, 30 * time.Second)
	default:
   		// Fast standard HTTP for Dustloop, FAT JSON, etc.
   		data, err = cfg.CharDataClient.Fetch(requestURL)
	}

	sp.Stop()

	if err != nil {
    	return fmt.Errorf("fetch %s: %w", requestURL, err)
	}
	
	
	var moves []models.Move

	// Select game parser based on game slug
	switch strings.ToLower(cfg.Game.Slug) {
	case "bbcf":
		moves, err = dustloop.BBCFCharPageParser(data)
		if err != nil {
			return err
		}

	case "ggst":
		moves, err = dustloop.GGSTCharPageParser(data)
		if err != nil {
			return err
		}

	case "ggxrd-r2", "ggacr":
		moves, err = dustloop.ParseLegacyGG(data)
		if err != nil {
			return err
		}

	case "gbvsr":
		moves, err = dustloop.ParseGBVSR(data)
		if err != nil {
			return err
		}
		
	case "dbfz":
		moves, err = dustloop.ParseDBFZ(data)
		if err != nil {
			return err
		}

	case "bbtag":
		moves, err = dustloop.ParseBBTAG(data)
		if err != nil {
			return err
		}

	case "mtfs":
		moves, err = dustloop.ParseMTFS(data)
		if err != nil {
			return err
		}
		
	case "uni2":
		moves, err = mizuumi.Uni2CharDataParser(data)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("no parser available for gameslug '%s' (name: %s)", cfg.Game.Slug, cfg.Game.Name)
	}

	character := models.Character{
		Name:	characterInput,
		Slug:	targetSlug,
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
	
	if err := sqlite.SaveMoves(cfg.DB, characterID, moves); err != nil {
		return err
	}
		
	fmt.Printf("Successfully fetched and saved %d moves for %s!\n", len(moves), characterInput)

	if verbose {
		for _, move := range moves {
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
			fmt.Printf("\n=== %s ===\n", title)

			for i, grid := range move.FrameDataGrids {
				if i == 0 {
					fmt.Println("[Base Frame Data]")
				} else {
					fmt.Printf("[Additional Data - Grid %d]\n", i)
				}
				display.PrintGrid(grid)
				fmt.Println()
			}
			if showDetails {
				if len(move.Notes) > 0 {
					fmt.Println("Notes:")
					for _, note := range move.Notes {
						fmt.Printf("  • %s\n", note)
					}
					fmt.Println()
				}
				if move.Description != "" {
					fmt.Println("")
					fmt.Println("Description:")
					fmt.Println(move.Description)
					fmt.Println()
				}
			}
			fmt.Println(strings.Repeat("=", 60))
		}
	} else {
		fmt.Printf("Tip: View frame data anytime with 'frames %s'\n", characterInput)
	}

	return nil
}

