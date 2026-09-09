package commands

import (
	"fmt"
	"strings"
	"time"
	"os"
	"path/filepath"

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

	// Determine if a JSON file will be saved by the command
	saveJSON := true
	for _, arg := range args {
		switch strings.ToLower(arg) {
			case "--no-json", "-n":
			saveJSON = false
		}
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

	// Change request method based on wiki to resolve anti-bot challenges
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

	if saveJSON {
		dataDir, err := ensureDataDir(cfg)
		if err != nil {
			return err
		}

		// Set filename and filePath for JSON save 
		filename := fmt.Sprintf("%s_characters.json", strings.ToLower(cfg.Game.Slug))
		fullPath := filepath.Join(dataDir, filename)

		if err := CharDataCache.SaveCharacters(
			filename,
			cfg.Game.Slug,
			characters,
		); err != nil {
			return err
		}

		fmt.Printf("Saved discovered chars to: %s\n", fullPath)

		cache, err := discovery.LoadCharCache(filename)
		if err != nil {
			fmt.Printf("Loaded %d characters from cache\n", len(cache.Characters))
		}
	} else {
		fmt.Println("Skipped saving discovered characters to JSON file")
	}

	
	fmt.Printf("Discovered %d characters for %s\n", len(characters), cfg.Game.Name)
	return nil	
}


// Func to set datadir for JSON and other file saves.
func ensureDataDir(cfg *config.Config) (string, error) {
	if cfg.DataDir == "" {
		if cfg.RL != nil {
			// Set up a temporary prompt for directory question
			cfg.RL.SetPrompt("Enter path to preferred save location for JSON files (press enter for default './data'): ")
			input, err := cfg.RL.Readline()
			// Restore the default REPL prompt
			cfg.RL.SetPrompt("CharData > ")

			if err != nil {
				return "", err
			}

			trimmed := strings.TrimSpace(input)
			if trimmed == "" {
				cfg.DataDir = "./data"
			} else {
				cfg.DataDir = trimmed
			}
		} else {
			cfg.DataDir = "./data"
		}
	}

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return "", fmt.Errorf("create data directory `%s`: %w", cfg.DataDir, err)
	}

	return cfg.DataDir, nil
}


