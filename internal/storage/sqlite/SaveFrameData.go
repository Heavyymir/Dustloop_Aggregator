package sqlite

import(
        "database/sql"
        "fmt"

        "github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Func to save framedata to the sql table
func SaveFrameData(db *sql.DB, moveID int64, position int, property string, cell models.Cell) error {

		_, err := db.Exec(`
			INSERT INTO frame_data
				(move_id, property, value, tootip, position)
			VALUES (?, ?, ?, ?, ?)
			`, 
			moveID, 
			property, 
			cell.Value,
			)

			if err != nil {
				return fmt.Errorf("save frame data: %w", err)
		}
	return nil
}
