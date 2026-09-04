package sqlite

import(
        "database/sql"
        "fmt"

        "github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Function to retrieve all of a Characters Moves from internal SQL DB
func GetMoves(db *sql.DB, characterID int64) ([]models.Move, error) {
	
	// Query Data base, save to rows
	rows, err := db.Query(`
		SELECT id, input, name
		FROM moves
		WHERE character_id = ?
		ORDER BY id
		`, characterID)

	if err != nil {
		return nil, fmt.Errorf("get moves: %w", err)
	}

	defer rows.Close()

	// Initialise Slice to hold moves collected from SQL table
	var moves []models.Move

	// Iterate through rows
	for rows.Next() {
		var move models.Move
		var	moveID int64

		// Assign values to internal move struct
		if err := rows.Scan(
			&moveID,
			&move.Input,
			&move.Name,
			&move.Description,
			); err != nil {
				return nil, fmt.Errorf("scan move: %w", err)
			}

			move.FrameDataGrids, err = GetMoveFrameDataGrids(db, moveID)
			if err != nil {
				return nil, err
			}

			// Append move struct to moves
			moves = append(moves, move)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate moves: %w", err)
	}
	
	return moves, nil
}

// Function to collect data for an individual move
func GetMove(db *sql.DB, moveID int64) (models.Move, error) {
	// Intialise Move struct
	var move models.Move

	// Query Moves table, assign element values with Scan()
	err := db.QueryRow(`
		SELECT input, name
		FROM moves
		WHERE id = ?
	`, moveID).Scan(
		&move.Input, 
		&move.Name,
	) 

	if err != nil {
		return models.Move{}, fmt.Errorf("get move: %w", err)
	}

	// call framedata helper to query tables and assign framedata to move struct
	move.FrameDataGrids, err = GetMoveFrameDataGrids(db, moveID)
	if err != nil {
		return models.Move{}, fmt.Errorf("get move frame data rows: %w", err)
	}

	return move, nil
}

// Function to query and collect framedata for a single cell
// in the SQL frame_data grid.
func GetMoveFrameData(db *sql.DB, moveID int64) (map[string]models.Cell, error) {
	// Query framedata table, assign to rows
	rows, err := db.Query(`
		SELECT property, value
		FROM frame_data
		WHERE move_id = ?
	`, moveID)

	if err != nil {
		return nil, fmt.Errorf("get frame data: %w", err) 
	}
	
	defer rows.Close()

	// Intialise framedata map to match Move struct elements
	frameData := make(map[string]models.Cell)

	// Iterate over queried data in rows
	for rows.Next() {
		var property, value string

		// Use Scan to assign values to the framedata map on each pass
		if err := rows.Scan(
			&property,
			&value,
		); err != nil {
			return nil, fmt.Errorf("scan frame data: %w", err)
		}

		// Assign cell value to the framedata property key
		frameData[property] = models.Cell{Value: value}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate frame data: %w", err)
	}

	return frameData, nil
}


// Helper to return frame data rows.
// Data is stored as a flat row, and this function
// rebuilds the row -> cell structure
func GetMoveFrameDataRows(db *sql.DB, moveID int64) ([]models.FrameDataRow, error) {

	// Query rows using input moveID
	rows, err := db.Query(`
		SELECT row_index, property, value, tooltip, position
		FROM frame_data
		WHERE move_id = ?
		ORDER BY row_index, position
		`, moveID)
		if err != nil {
			return nil, fmt.Errorf("get frame data rows: %w", err)
		}

		defer rows.Close()
		
		var frameDataRows []models.FrameDataRow

		// Loop over rows
		for rows.Next() {
			// Initalise variables matching DB rows
			var rowIndex, position int
			var property, value, tooltip string

			// Assign values to variables using .Scan()
			if err := rows.Scan(
				&rowIndex,
				&property,
				&value,
				&tooltip,
				&position,
			); err != nil {
				return nil, fmt.Errorf("scan frame data row: %w", err)
			}

			// Append any rows less than the found rowIndex to frameDataRows slice
			for len(frameDataRows) <= rowIndex {
				frameDataRows = append(frameDataRows, models.FrameDataRow{})
			}

			// Append the found value and tooltip to a cell at the rowIndex position in the frameDataRows slice
			frameDataRows[rowIndex].Cells = append(
				frameDataRows[rowIndex].Cells,
				models.Cell{
					Value: value,
					Tooltip: tooltip,
				},
			)

		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterate frame data rows: %w", err)
		}
		
	return frameDataRows, nil 
	}


// Helper to return frame data grids for a move
// The DB stores each cell as a flat row so this function rebuilds
// the nested structure as grids -> rows -> cells
func GetMoveFrameDataGrids(db *sql.DB, moveID int64) ([]models.FrameDataGrid, error) {

	// Query database for all saved frame data information
	// Ordering ensures cells are processed by grid, row, then column
	rows, err := db.Query(`
		SELECT grid_index, row_index, property, value, tooltip, position
		FROM frame_data
		WHERE move_id = ?
		ORDER BY grid_index, row_index, position
		`, moveID)

	if err != nil {
		return nil, fmt.Errorf("get frame data grids: %w", err)
	}

	defer rows.Close()

	var grids []models.FrameDataGrid

	// Iterate over rows data, creating varaibles and scanning in relevent value
	for rows.Next() {
		var gridIndex, rowIndex, position int
		var property, value, tooltip string

		if err := rows.Scan(
			&gridIndex,
			&rowIndex,
			&property,
			&value,
			&tooltip,
			&position,
		); err!= nil {
			return nil, fmt.Errorf("scan frame data grid: %w", err)
		}

		// Ensure the outer sice contains an entry for the database's grid index
		// Empty entries are added when a later grid index is encountered
		for len(grids) <= gridIndex {
			grids = append(grids, models.FrameDataGrid {
				Headers: []string{},
				Rows:	 []models.FrameDataRow{},
			})
		}

		// Select the current grid index so that its headers and rows can be populated
		grid := &grids[gridIndex]

		// Ensure that the headers slice is large enough for the cell's column position
		for len(grid.Headers) <= position {
			grid.Headers = append(grid.Headers, "")
		}

		// Store the property name as the header for this column
		grid.Headers[position] = property

		// Ensure that the grid contains the row represented by the rowIndex
		for len(grid.Rows) <= rowIndex {
			grid.Rows = append(grid.Rows, models.FrameDataRow{
				Cells: []models.Cell{},
			})
		}

		// Select the current row so its cells can be populated
		row := &grid.Rows[rowIndex]

		// Ensure the row contains a cell at this column position
		for len(row.Cells) <= position {
			row.Cells = append(row.Cells, models.Cell{})
		}

		// Store the DB value and tooltip at the original column position as a cell
		row.Cells[position] = models.Cell{
			Value: value,
			Tooltip: tooltip,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate frame data grids: %w", err)
	}

	return grids, nil
}

