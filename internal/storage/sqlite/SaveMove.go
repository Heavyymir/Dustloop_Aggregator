package sqlite


import(
        "database/sql"
        "fmt"

        "github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Function to save a characters moves to the DB
func SaveMoves(db *sql.DB, characterID int64, moves []models.Move) error {

	// Start the SQL transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin move transaction: %w", err)
	}

	// Defer Rollback to catch issues with commit if they occur
	defer tx.Rollback()

	// Delete existing frame data for the selected characters moves
	_, err = tx.Exec(`
		DELETE FROM frame_data
		WHERE move_id IN (SELECT id FROM moves WHERE character_id = ?)
		`, characterID)
	if err != nil {
		return fmt.Errorf("clear existing frame data: %w", err)
	}

	// Delete existing moves for this character
	_, err = tx.Exec(`
		DELETE FROM moves
		WHERE character_id = ?
		`, characterID)
	if err != nil {
		return fmt.Errorf("clear existing moves: %w", err)
	}

	// Loop over moves slice and call saveMove on each entry
	for _, move := range moves {
		if _, err := saveMove(tx, characterID, move); err != nil {
			return err
		}
	}

	// Once the loop completes, call .Commit() to finalise changes to the SQL table
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit moves: %w", err)
	}
	return nil
}
				
// Helper fo SaveMoves to handle the Insertion of individal moves into the DB 
func saveMove(tx *sql.Tx, characterID int64, move models.Move) (int64, error) {

	// Call .Exec() method on the input SQL transaction to insert data into DB
	result, err := tx.Exec(`
	INSERT INTO moves (character_id, input, name, description)
	VALUES (?, ?, ?, ?)
	`, characterID, move.Input, move.Name, move.Description)

	if err != nil {
		return 0, fmt.Errorf("save move: %w", err)
	}

	// Assign a Move ID to the insertion
	moveID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get move ID: %w", err)
	}
	// Nested loops to pull values from HTML
	for gridIndex, grid := range move.FrameDataGrids {
		for rowIndex, row := range grid.Rows {
			for position, cell := range row.Cells {
				if position >= len(grid.Headers) {
					continue
				}
	
				property := grid.Headers[position]

	
				_, err := tx.Exec(`
					INSERT INTO frame_data
					(move_id, grid_index, row_index, property, value, tooltip, position)
					VALUES (?, ?, ?, ?, ?, ?, ?)
				`,
					moveID,
					gridIndex,
					rowIndex,
					property,
					cell.Value,
					cell.Tooltip,
					position,
				)
	
				if err != nil {
					return 0, fmt.Errorf(
						"save frame data: %w", err,
					)
				}
			}
		}
	}
	return moveID, nil
}

