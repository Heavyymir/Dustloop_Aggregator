package mizuumi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
	"github.com/PuerkitoBio/goquery"
)

// Function to parse fetched Uni2 character data html bodies
func Uni2CharDataParser(data []byte) ([]models.Move, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	var moves []models.Move

	// Find each move inside wrapped .movedata-container
	doc.Find(".movedata-container").Each(func(_ int, container *goquery.Selection) {
		move := models.Move{
			FrameDataGrids: []models.FrameDataGrid{},
			Notes:          []string{},
		}

		// Find and assign move input or heading name (from preceding h5 or mw-heading5 element)
		heading := container.PrevAll().Filter("div.mw-heading5, h5").First()
		if heading.Length() > 0 {
			move.Input = strings.TrimSpace(heading.Text())
		}

		// Find the full move title (from the <big> element in the image container)
		bigTitle := strings.TrimSpace(container.Find(".movedata-flex-image-container big").First().Text())
		if bigTitle != "" {
			move.Name = bigTitle
		} else if move.Input != "" {
			move.Name = move.Input
		} else {
			move.Name = "Unknown Move"
		}

		// Extract the frame data grids from base tables and movedata table
		container.Find("table.movedata-flex-framedata-table").Each(func(_ int, table *goquery.Selection) {
			var headers []string

			// Extract table headers
			table.Find("tr").First().Find("th").Each(func(_ int, th *goquery.Selection) {
				headerText := strings.Join(strings.Fields(th.Text()), " ")
				headers = append(headers, headerText)
			})

			if len(headers) == 0 {
				return
			}

			grid := models.FrameDataGrid{
				Headers: headers,
				Rows:    []models.FrameDataRow{},
			}

			// Extract Data Rows
			table.Find("tbody > tr").Each(func(i int, tr *goquery.Selection) {
				// Skip header row if it only contains <th> element and no <td>
				if tr.Find("td").Length() == 0 {
					return
				}

				row := models.FrameDataRow{
					Cells: []models.Cell{},
				}

				// Some rows start with a <th> element for a move version (e.g., "B + C", "[B] + [C]")
				tr.ChildrenFiltered("th, td").Each(func(_ int, cellNode *goquery.Selection) {
					tooltip, _ := cellNode.Find("[title]").Attr("title")
					val := strings.Join(strings.Fields(cellNode.Text()), " ")

					row.Cells = append(row.Cells, models.Cell{
						Value:   val,
						Tooltip: strings.TrimSpace(tooltip),
					})
				})

				if len(row.Cells) > 0 {
					grid.Rows = append(grid.Rows, row)
				}
			})

			if len(grid.Rows) > 0 {
				move.FrameDataGrids = append(move.FrameDataGrids, grid)
			}
		})

		
		// Gather move description and notes
		var descriptions []string
		container.Find(".movedata-flex-description").Each(func(_ int, descNode *goquery.Selection) {
			// Extract move bullet points as notes
			descNode.Find("ul > li").Each(func(_ int, li *goquery.Selection) {
				noteText := strings.TrimSpace(li.Text())
				if noteText != "" {
					move.Notes = append(move.Notes, noteText)
				}
			})

			// Clone and remove <ul> elements to get description text
			descClone := descNode.Clone()
			descClone.Find("ul").Remove()
			descClone.Find(".mw-headline, h6, h5, b, strong").Each(func(_ int, label *goquery.Selection) {
				// If the element is just a move label (eg. "2a :" or "5[C] :")	strip it
				txt := strings.TrimSpace(label.Text())
				if strings.HasSuffix(txt, ":") || txt == move.Input || txt == move.Name {
					label.Remove()
				}
			})
			
			// Clean and filter individual lines of the description text
			rawText := descClone.Text()
			var cleanLines []string
			for _, line := range strings.Split(rawText, "\n") {
				trimmed := strings.TrimSpace(line)
				// Filter out empty lines or orphan labels like "2A :" or "2A"
				if trimmed == "" || strings.HasSuffix(trimmed, ":") || trimmed == move.Input {
					continue
				}
				cleanLines = append(cleanLines, trimmed)
			}
		
			if len(cleanLines) > 0 {
				descriptions = append(descriptions, strings.Join(cleanLines, "\n"))
			}
		})

		move.Description = strings.Join(descriptions, "\n")

		// Only append if func found framedata or valid move information
		if len(move.FrameDataGrids) > 0 || move.Name != "Unknown Move" {
			moves = append(moves, move)
		}
	}) 

	if len(moves) == 0 {
		return nil, fmt.Errorf("no moves found in character HTML")
	}
	return moves, nil
}
