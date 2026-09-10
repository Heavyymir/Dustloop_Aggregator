package commands

import (
	"fmt"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// helper to print parsed CharData to console in a readable format
func printMoveTable(moves []models.Move) {
    for _, move := range moves {
    	fmt.Println("---------------------------------------------")
    	if move.Input != "" && move.Input != move.Name {
        	fmt.Printf("\n%s %s\n", move.Name, move.Input)
		} else {
			fmt.Printf("%s\n", move.Name)
		}
        	
        fmt.Println("---------------------------------------------")

        for gridIndex, grid := range move.FrameDataGrids {
            fmt.Printf("Frame data grid %d:\n", gridIndex)

            for _, header := range grid.Headers {
                fmt.Printf("%-18s", header)
            }
            fmt.Println()

            for _, row := range grid.Rows {
                for _, cell := range row.Cells {
                    fmt.Printf("%-18s", cell.Value)
                }
                fmt.Println()
            }
        }

        if move.Description != "" {
            fmt.Printf("Description: %s\n", move.Description)
        }
    }
}
