package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Heavyymir/Dustloop_Aggregator/catalog"
)

// Build the URL to get character HTML data for parsing

func CharacterURL(wiki catalog.Wiki, game catalog.Game, characterName string) (string, error) {
	if strings.TrimSpace(characterName) == "" {
		return "", fmt.Errorf("character cannot be empty")
	}

	// Escape the character name to build the intended URL.
	characterPath := strings.Replace(
		game.CharacterPath,
		"{character}",
		url.PathEscape(characterName),
		1,
	)
	return wiki.URL + characterPath, nil
}
