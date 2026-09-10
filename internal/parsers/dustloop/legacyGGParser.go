package dustloop

import(
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// HTML Parser for GGXRD-REV2 Character Pages
func ParseLegacyGG(data []byte) ([]models.Move, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var moves []models.Move

	doc.Find(".attack-container").Each(func(_ int, container *goquery.Selection) {
		move := models.Move{
			FrameDataGrids:		[]models.FrameDataGrid{},
			Notes:				[]string{},
		}

		// Find heading proceeding this attack container
		heading := container.PrevAll().Filter("div.mw-heading3, div.mw-heading4, h3, h4").First()
		moveName := strings.TrimSpace(heading.Text())

		// Check for input badge right after the heading or before the attack container
		inputBadge := container.PrevUntil("div.mw-heading3, div.mw-heading4, h3, h4").Filter("p").Find("input-badge")
		moveInput := strings.TrimSpace(inputBadge.Text())
		if moveInput == "" {
			moveInput = moveName  // Specifically for Normals (5P, 6K, JS etc.)
		}

		move.Name = moveName

		// Only assign Input if a badge is found and if it is distinct from the name
		if move.Name != "" && move.Input != moveName{
			move.Input = moveInput
		}

		// Extract grids (handles single, multi)
		container.Find(".frameDataGrid").Each(func(_ int, gridNode *goquery.Selection) {
			var headers []string
			gridNode.Find(".frameDataGridHeader").First().ChildrenFiltered("div").Each(
				func(_ int, cell *goquery.Selection) {
					visible := cell.Clone()
					visible.Find(".tooltiptext").Remove()
					headerName := strings.Join(strings.Fields(visible.Text()), " ")
					headers = append(headers, headerName)
				},
			)

			grid := models.FrameDataGrid{
				Headers:	headers,
				Rows:		[]models.FrameDataRow{},
			}

			// Extract Frame data rows
			gridNode.Find(".frameDataGridRow").Each(func(_ int, rowNode *goquery.Selection) {
				frameRow := models.FrameDataRow{
					Cells:		[]models.Cell{},
				}

				rowNode.ChildrenFiltered("div").Each(func(_ int, cellNode *goquery.Selection) {
					frameRow.Cells = append(frameRow.Cells, parseCell(cellNode))
				})

				grid.Rows = append(grid.Rows, frameRow)
			})
			
			move.FrameDataGrids = append(move.FrameDataGrids, grid)
		})     

		// Extract move descriptions
		var paragraphs []string
		container.Find(".attack-info-body > p").Each(func(_ int, p *goquery.Selection) {
			pClone := p.Clone()
			pClone.Find("style, script, .tmp, .tooltiptext, .cargo-hover, .mvd-hove, .ext-popups").Remove()

			text := strings.TrimSpace(pClone.Text())
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		})

		move.Description = strings.Join(paragraphs, "\n")

		// Extract Move notes
		container.Find(".attack-info-body > ul > li").Each(func(_ int, li *goquery.Selection) {
			text := strings.TrimSpace(li.Text())
			if text	!= "" {
				move.Notes = append(move.Notes, text)
			}
		})

		moves = append(moves, move)

		})

		return moves, nil
}
