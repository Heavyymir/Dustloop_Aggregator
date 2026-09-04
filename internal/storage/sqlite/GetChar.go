package sqlite

import(
	"database/sql"
	"fmt"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Function to fetch character data from internal SQL db
func GetCharacter(db *sql.DB, game, slug string) (models.Character, error) {
	// Initialise a Character struct for return
	var character models.Character

	// Query the table for the required data, use Scan to assign the Character Struct elements
	err := db.QueryRow(`
	SELECT name, slug, source_url
	FROM characters
	WHERE game = ? AND slug = ?
	`, game, slug).Scan(
		&character.Name,
		&character.Slug,
		&character.URL,
	)

	if err != nil {
		return models.Character{}, fmt.Errorf("get character: %w", err)
	}

	return character, nil
}
