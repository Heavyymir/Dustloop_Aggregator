package sqlite

import(
	"database/sql"
	"fmt"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Function to list characters present in the SQL database based on a game
func ListCharacters(db *sql.DB, game string) ([]models.Character, error) {

	// Initialise rows, query characters for character names using input games string
	rows, err := db.Query(`
	SELECT name, slug, source_url
	FROM characters
	WHERE game = ?
	ORDER BY name
	`, game)
	if err != nil {
		return nil, fmt.Errorf("list character: %w", err)
	}

	defer rows.Close()

	// Initialise the external characters slic
	var characters []models.Character

	// Iterate over each entry in rows
	for rows.Next() {
		// Initialise an interal character struct
		var character models.Character

		// Use scan to assign values from the SQL table to the character stuct
		if err := rows.Scan(
			&character.Name,
			&character.Slug,
			&character.URL,
		); err != nil {
			return nil, fmt.Errorf("iterate characters %w", err)
		}

		// Append character to characters slice
		characters = append(characters, character)
	}

	return characters, nil
}
