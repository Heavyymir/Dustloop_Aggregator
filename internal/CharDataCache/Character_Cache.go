package CharDataCache

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

type CharacterCache struct {
	GameSlug   string             `json:"game_slug"`
	Characters []models.Character `json:"characters"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

func SaveCharacters(filename string, gameSlug string, characters []models.Character) error {
	cache := CharacterCache{
		GameSlug:   gameSlug,
		Characters: characters,
		UpdatedAt:  time.Now(),
	}

	data, err := json.MarshalIndent(cache, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
