package commands

import (
	"fmt"
	"strings"
	"time"
	"os"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/parsers/bbcf"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/parsers/ggst"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/parsers/fat"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/parsers/mizuumi"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/discovery"
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

	characterInput := args[0]
	targetSlug := characterInput

	// 1. Try to load the discovered characters cache for this game
	filename := cfg.Game.Slug + "_characters.json"
	cache, err := discovery.LoadCharCache(filename)
	if err == nil {
    	// 2. Look for a case-insensitive match (e.g. "hyde" matches "Hyde", "RYU" matches "ryu")
    	for _, char := range cache.Characters {
        	if strings.EqualFold(char.Name, characterInput) || strings.EqualFold(char.Slug, characterInput) {
            	targetSlug = char.Slug
            	break
        	}
    	}
	}

	// 3. Build the URL using the matched slug
	pagePath := strings.Replace(cfg.Game.CharacterPath, "{character}", targetSlug, 1)
	requestURL := fmt.Sprintf("%s/%s", strings.TrimRight(cfg.Wiki.URL, "/"), pagePath)

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

	if err != nil {
    	return fmt.Errorf("fetch %s: %w", requestURL, err)
	}

	// 1. Log the fetched size and URL
	fmt.Printf("[DEBUG] Fetched %d bytes from %s\n", len(data), requestURL)
	
	// 2. Dump HTML to disk to inspect the raw response
	if err := os.WriteFile("debug_page.html", data, 0644); err != nil {
		fmt.Printf("[DEBUG] Failed to write dump: %v\n", err)
	}
	
	var moves []models.Move
	
	switch strings.ToLower(cfg.Game.Slug) {
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
		
	printMoveTable(moves)

	// Return request data to the user so they can see it has been successful
	fmt.Printf("fetched %d bytes from %s\n", len(data), requestURL)

	return nil
}
