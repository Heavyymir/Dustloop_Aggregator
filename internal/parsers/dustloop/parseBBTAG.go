package dustloop

import(
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Parser to extract move framedata and descriptions from BBTag character page
func ParseBBTAG(data []byte) ([]models.Move, error) {
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

		// Extract move name
		heading := container.PrevAll().Filter("div.mw-heading3, h3, h4").First()
		if heading.Length() > 0 {
			rawName := heading.Find("h3, h4")
			if rawName.Length() == 0 {
				rawName = heading
			}
			move.Name = strings.Join(strings.Fields(rawName.Text()), " ")
		}

		// Input Notation
		inputBadge := container.PrevAll().Filter("p").First().Find(".input-badge")
		if inputBadge.Length() > 0 {
			move.Input = strings.Join(strings.Fields(inputBadge.Text()), " ")
		} else {
			// Input notation for normals
			move.Input = move.Name
		}

		// Extract framedata grids
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

			gridNode.Find(".frameDataGridRow").Each(func(_ int, rowNode *goquery.Selection) {
				frameRow := models.FrameDataRow{
					Cells:	[]models.Cell{},
				}

				rowNode.ChildrenFiltered("div").Each(func(_ int, cellNode *goquery.Selection) {
					frameRow.Cells = append(frameRow.Cells, parseCell(cellNode))
				})

				grid.Rows = append(grid.Rows, frameRow)
			})

			move.FrameDataGrids = append(move.FrameDataGrids, grid)
		})

		//Extract move descriptions
		infoBody := container.Find(".attack-info-body").Clone()
		infoBody.Find("script, style, link, details, .frameChart, .tooltiptext, sup.reference").Remove()

		var paragraphs []string
		infoBody.Find("p").Each(func(_ int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		})

		// Fallback for desctiptions not using the <p> tag
		if len(paragraphs) == 0 {
			clone := infoBody.Clone()
			clone.Find("ul").Remove()
			rawText := strings.TrimSpace(clone.Text())
			if rawText != "" {
				paragraphs = append(paragraphs, rawText)
			}
		}

		move.Description = strings.Join(paragraphs, "\n")

		// Extract notes / bullet points
		infoBody.Find("ul > li").Each(func(_ int, li *goquery.Selection) {
			text := strings.TrimSpace(li.Text())
			if text != "" {
				move.Notes = append(move.Notes, text)
			}
		})	

		moves = append(moves, move)
	})

	return moves, nil
}
