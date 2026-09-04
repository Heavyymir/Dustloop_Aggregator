package commands

import(
	"fmt"

	"github.com/Heavyymir/Dustloop_Aggregator/config"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/storage/sqlite"
)

// Command to list all character entries in the DB based on a selected game
func commandList(cfg *config.Config, args... string) error {
	if cfg.DB == nil {
		return fmt.Errorf("Database is not configured")
	}

	if cfg.Game.Name == "" {
		return fmt.Errorf("select a game first")
	}

	// Call ListCharacters SQL function, assign to characters
	characters, err := sqlite.ListCharacters(cfg.DB, cfg.Game.Slug)
	if err != nil {
		return err
	}

	// Walk through characters, print each character name and their game slug
	for _, character := range characters{
		fmt.Printf("%s (%s)\n", character.Name, character.Slug)
	}

	// Return a confirmation line
	fmt.Printf("Stored %d characters for %s\n", len(characters), cfg.Game.Name)
	return nil
}
