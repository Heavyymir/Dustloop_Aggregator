package fat

import(
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Heavyymir/CharData_Aggregator/internal/models"
)

// FATJSONParser takes the raw JSON bytes from FAT and transforms it into models.Move slices
func FATJSONParser(data []byte) ([]models.Move, error) {
	// Attempt to decode the FAT data into a slice of moves first
	var rawMoves []FATMove
	if err := json.Unmarshal(data, &rawMoves); err != nil {
		// As a fallback, try decoding the wrapped { "moves": [...] } payload
		var wrapper FATCharacterPayload
		if errWrapper := json.Unmarshal(data, &wrapper); errWrapper != nil {
			return nil, fmt.Errorf("failed to unmarshal FAT json: %w", err)
		}
		rawMoves = wrapper.Moves
	}

	if len(rawMoves) == 0 {
		return nil, fmt.Errorf("no moves found in FAT data")
	}

	// Standard columns for FAT frame data grids
	standardHeaders := []string {
		"Input", "Damage", "Startup", "Active", "Recovery", "On Block", "On Hit", "Cancel",
	}	

	var parsedMoves []models.Move

	for i, rm := range rawMoves {
		name := strings.TrimSpace(rm.MoveName)
		if name == "" {
			name = rm.Input
		}
		if name == "" {
			name = fmt.Sprintf("Move %d", i)
		}

		// Constuct the row cells matching the standard headers
		cells := []models.Cell{
			{Value: rm.Input},
			{Value: rm.Damage},
			{Value: rm.Startup},
			{Value: rm.Active},
			{Value: rm.Recovery},
			{Value: rm.OnBlock},
			{Value: rm.OnHit},
			{Value: rm.Cancel},
		}

		grid := models.FrameDataGrid{
			Headers: standardHeaders,
			Rows: []models.FrameDataRow{
				{Cells: cells},
			},
		}

		var notes []string
		if strings.TrimSpace(rm.Comments) != "" {
			notes = append(notes, strings.TrimSpace(rm.Comments))
		}

		move := models.Move{
			Name:			name,
			Input: 			rm.Input,
			FrameDataGrids: []models.FrameDataGrid{grid},
			Notes:			notes,
		}

		parsedMoves = append(parsedMoves, move)
	}

	return parsedMoves, nil
}
