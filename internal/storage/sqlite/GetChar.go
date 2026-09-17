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
	WHERE LOWER(game) = LOWER(?) 
		AND (
			LOWER(slug) = LOWER(?) 
			OR LOWER(name) = LOWER(?)
			OR REPLACE (LOWER(slug), '_', ' ') = REPLACE(LOWER(?), '_', ' ')
			OR REPLACE (LOWER(slug), '-', ' ') = REPLACE(LOWER(?), '-', ' ')
		)
		LIMIT 1;
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

func GetCharacterNames(db *sql.DB, gameSlug string) ([]string, error) {
	query := `SELECT name FROM characters WHERE game = ? ORDER BY name ASC`
	rows, err := db.Query(query, gameSlug)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, rows.Err()
}
