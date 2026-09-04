package discovery

import(
	"encoding/json"
	"os"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/CharDataCache"
)

func LoadCharCache(path string) (CharDataCache.CharacterCache, error) {
	file, err := os.Open(path)
	if err != nil {
		return CharDataCache.CharacterCache{}, err
	}

	defer file.Close()

	var cache CharDataCache.CharacterCache

	err = json.NewDecoder(file).Decode(&cache)
	if err != nil {
		return CharDataCache.CharacterCache{}, nil
	}

	return cache, nil
}
