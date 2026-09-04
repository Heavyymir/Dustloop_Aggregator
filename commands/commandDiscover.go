package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/CharDataCache"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/discovery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
)

// Command to find all available characters for a game and save them to the local SQL DB 
func commandDiscover(cfg *config.Config, args ...string) error {
	if cfg.Wiki.Name == "" || cfg.Game.Name == "" {
		return fmt.Errorf("Select a wiki and game first")
	}

	// Format the wikislug to use for the URL
	wikiSlug := strings.ToUpper(cfg.Game.Slug)

	// Create the indexURL
	indexURL := cfg.Wiki.URL + wikiSlug

	// Debug the index URL
	fmt.Printf("index URL: %s\n", indexURL)

	// Fetch the HTML page using the Index URL
	var data []byte
	var err error
	
	switch strings.ToLower(cfg.Wiki.Slug) {
	case "mizuumi", "supercombo":
		fmt.Println("Using headless fetch method for roster page...")
		data, err = cfg.CharDataClient.FetchHTMLHeadless(indexURL, 30 * time.Second)
	default:
		data, err = cfg.CharDataClient.Fetch(indexURL)
	}
	
	if err != nil {
		return err
	}

	// Debug for data fetching
	fmt.Printf("fetched %d bytes\n", len(data))
	
	// Call DiscoveredChars to create a characters slice
	characters, err := discovery.DiscoveredChars(data, wikiSlug)
	if err != nil {
		return err
	}

	// Iterate over the slice, print characters to console and Save them to the internal SQL DB
	for _, character := range characters {
		fmt.Printf("%s -> %s\n", character.Name, character.URL)
		_, err := sqlite.SaveCharacter(
			cfg.DB,
			cfg.Game.Slug,
			character,
		)
		if err != nil {
			return fmt.Errorf("save character %s: %w", character.Name, err)
		}
	}

	filename := cfg.Game.Slug + "_characters" + ".json"

	if err := CharDataCache.SaveCharacters(
		filename,
		cfg.Game.Slug,
		characters,
	); err != nil {
		return err
	}

	cache, err := discovery.LoadCharCache(filename)
	if err != nil {
		return err
	}


	fmt.Printf("Discovered %d characters for %s\n", len(characters), cfg.Game.Name)
	fmt.Printf("Loaded %d characters\n", len(cache.Characters))
	return nil
}
