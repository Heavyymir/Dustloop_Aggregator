package sqlite

import(
	"database/sql"
	"fmt"
	"strings"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
) 

// Function to insert character data into internal SQL tables
func SaveCharacter(db *sql.DB, gameSlug string, character models.Character) (int64, error) {
	normGame := strings.ToLower(gameSlug)
	normSlug := strings.ToLower(character.Slug)

	// Use db.Exec to insert a character into the Database using the internal Character Struct
	_, err := db.Exec(`
		INSERT INTO characters (game, name, slug, source_url)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(game, slug) DO UPDATE SET
			name = excluded.name,
			source_url = excluded.source_url
		`, 
		normGame, 
		character.Name, 
		normSlug, 
		character.URL,
		)

	if err != nil {
		return 0, fmt.Errorf("save character: %w", err)
		}

	// Initialise id
	var id int64

	// Query the id Column from characters, Scan to assign an ID to the new entry
	err = db.QueryRow(`
		SELECT id FROM characters
		WHERE game = ? AND slug = ?
		`,
		normGame,
		normSlug,
		).Scan(&id)
		
	if err != nil {
		return 0, fmt.Errorf("get character ID: %w", err)
	}

	// Return id
	return id, nil
}

