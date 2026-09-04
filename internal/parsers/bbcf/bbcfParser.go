package bbcf

import (
	"bytes"
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

func BBCFCharPageParser(data []byte) ([]models.Move, error) {

	// Use goquery to create a readable datapoint from page HTML data
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Initialise a slice of Move structs
	var moves []models.Move
	
	// Find all attack containers in the created doc
	doc.Find(".attack-container").Each(func(i int, container *goquery.Selection) {
		paragraphs := []string{}
		// Initialise move struct maps/slices for elements
		move := models.Move{
			FrameDataGrids: []models.FrameDataGrid{},
		}

		// Assign move names
		heading := container.Prev()
		if heading.Is("p") {
		    heading = heading.Prev()
		}
		
		move.Name = strings.Join(strings.Fields(heading.Find("h3").Text()), " ")
		
		// Assign move inputs
		inputNode := container.Prev()
		inputBadge := inputNode.Find(".input-badge")
		
		if inputBadge.Length() > 0 {
		    move.Input = strings.Join(strings.Fields(inputBadge.Text()), " ")
		} else {
		    move.Input = strings.Join(
		        strings.Fields(inputNode.Find("h3").Text()),
		        " ",
		    )
		}
		
		
		// Find each frameDataGrid in container
		container.Find(".frameDataGrid").Each(func(_ int, grid *goquery.Selection) {
			var headers []string

			// Find each header in the grid, append to headers slice
			grid.Find(".frameDataGridHeader").First().ChildrenFiltered("div").Each(
			    func(_ int, cell *goquery.Selection) {
			        visible := cell.Clone()
			        visible.Find(".tooltiptext").Remove()
			
			        headers = append(
			            headers,
			            strings.Join(strings.Fields(visible.Text()), " "),
			        )
			    },
			)

			// Initialise a FrameDataGrid struct, assign headers to element and initialise rows
			currentGrid := models.FrameDataGrid{
				Headers: headers,
				Rows:    []models.FrameDataRow{},
			}

			// find each row in the grid
			grid.Find(".frameDataGridRow").Each(func(_ int, row *goquery.Selection) {
					// Initialise a framerow struct with a slice of Cell structs
					frameRow := models.FrameDataRow{
						Cells: []models.Cell{},
					}

					row.ChildrenFiltered("div").Each(
					func(_ int, cell *goquery.Selection) {
						frameRow.Cells = append(
							frameRow.Cells, 
							parseCell(cell),
						)
					},
				)
	
				
				// Append framerows to the currentGrid.Rows slice		
				currentGrid.Rows = append(currentGrid.Rows, frameRow)
			},
		)
		// Append the completed grid to move.FrameDataGrids
		move.FrameDataGrids = append(move.FrameDataGrids, currentGrid)

		})

		for gi, grid := range move.FrameDataGrids {
		    for ri, row := range grid.Rows {
		        for ci, cell := range row.Cells {
		            fmt.Printf("%s grid=%d row=%d cell=%d: %+v\n",
		                move.Name, gi, ri, ci, cell)
		        }
		    }
		}

		// Find the full move body, assign to a created paragraphs slice. 
		container.Find(".attack-info-body > p").Each(func(_ int, paragraph *goquery.Selection) {
			text := strings.TrimSpace(paragraph.Text())
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		})

		
		// Join the strings in the created paragraphs slice, assign to move.Description
		move.Description = strings.Join(paragraphs, "\n")

		// Append the created struct to the Moves Slice
		moves = append(moves, move)
	})

	return moves, nil
}

// Helper to parse goquery data
func parseCell(cell *goquery.Selection) models.Cell {
	toolTip := strings.TrimSpace(cell.Find(".tooltiptext").Text())

	visible := cell.Clone()
	visible.Find(".tooltiptext").Remove()

	return models.Cell{
		Value:	strings.TrimSpace(visible.Text()),
		Tooltip: toolTip,
	}
}


