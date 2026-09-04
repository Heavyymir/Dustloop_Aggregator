package commands

import (
	"errors"
	"fmt"

	"github.com/Heavyymir/Dustloop_Aggregator/catalog"
	"github.com/Heavyymir/Dustloop_Aggregator/config"
)

// Command to allow users to walk through selections of Wikis, Games and Characters
func commandSelect(cfg *config.Config, args ...string) error {
	switch {
	// Case for if a user attempts a wildcard search to see which wikis are available
	case len(args) == 1 && args[0] == "*":
		printWikis()
		return nil

	// Case to handle Wildcard game selection by printing available games to a user
	case len(args) == 2 && args[1] == "*":
		wiki, exists := catalog.Wikis[args[0]]
		if !exists {
			return fmt.Errorf("unknown wiki: %s", args[0])
		}

		printGames(wiki)
		return nil

	// Case to check for the games under a wiki that the program can handle.
	case len(args) == 2 && args[1] != "*":
		return selectWikiAndGame(cfg, args[0], args[1])

	default:
		return errors.New("usage: select * | select <wiki> * | select <wiki> <game>")
	}
}

// Wrapper to print wikis to console
func printWikis() {
	for key, wiki := range catalog.Wikis {
		fmt.Printf("%s: %s\n", key, wiki.Name)
	}
}

// Wrapper to print available games to console
func printGames(wiki catalog.Wiki) {
	for key, game := range wiki.Games {
		fmt.Printf("%s: %s\n", key, game.Name)
	}
}

func selectWikiAndGame(cfg *config.Config, wikiKey string, gameKey string) error {
	// Check to see if the selected Wiki is accepted by the program
	wiki, exists := catalog.Wikis[wikiKey]
	if !exists {
		return fmt.Errorf("unknown wiki: %s", wikiKey)
	}
	// Check to see if the selected game is covered by the selected Wiki
	game, exists := wiki.Games[gameKey]
	if !exists {
		return fmt.Errorf("unknown game: %s", gameKey)
	}

	// Store the selected wiki and game in config
	cfg.Wiki = wiki
	cfg.Game = game

	// Show the User their selected Wiki and Game.
	fmt.Printf("selected %s: %s\n", wiki.Name, game.Name)

	return nil
}
