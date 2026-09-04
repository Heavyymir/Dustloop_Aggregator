package sqlite

import(
	"database/sql"
	"fmt"
)

// Helper for commands to retrieve Character IDs from SQL tables
func GetCharacterID(db *sql.DB, game, slug string) (int64, error) {
	var id int64

	// Query DB, pull character ID using game and slug 
	err := db.QueryRow(`
		SELECT id
		FROM characters
		WHERE game = ? AND slug = ?
		`,
		game,
		slug,
		).Scan(&id)
		
		if err != nil {
			return 0, fmt.Errorf("get character ID: %w", err)
		}
	
	return id, nil
}
