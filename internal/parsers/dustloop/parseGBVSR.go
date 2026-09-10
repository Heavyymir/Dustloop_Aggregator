package dustloop

import(
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)


// Parser for Granblue Fantasy Versus: Rising HTML
func ParseGBVSR(data []byte) ([]models.Move, error) {
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

		// Handle dual inputs: Technical (.toggle-group-a) and simple (toggle-group-b)
		prevNode := container.Prev()

		// Check toggle groups first. If present, Group A represents the technical inputs
		techBadge := prevNode.Find(".toggle-group-a .input-badge")
		simpleBadge := prevNode.Find(".toggle-group-b .input-badge")

		if techBadge.Length() > 0 {
			techText := strings.Join(strings.Fields(techBadge.Text()), " ")
			if simpleBadge.Length() > 0 {
				simpleText := strings.Join(strings.Fields(simpleBadge.Text()), " ")
				move.Input = strings.TrimSpace(techText + " (Simple: " + simpleText + ")")
			} else {
				move.Input = techText
			}
		} else {
			// Fallback for normals or moves without a simple input toggles
			inputBadge := prevNode.Find(".input-badge")
			if inputBadge.Length() > 0 {
				move.Input = strings.Join(strings.Fields(inputBadge.Text()), " ")
			}
		}

		// Extract move name
		heading := container.PrevAll().Filter("div.mw-heading4, h4, div.mw-heading3, h3").First()
		if heading.Length() > 0 {
			rawName := heading.Find("h4, h3").Text()
			if rawName == "" {
				rawName = heading.Text()
			}
			move.Name = strings.Join(strings.Fields(rawName), " ")
		}

		// Assign move name as move input if move name matches any of the below conditions
		if move.Name == "" || strings.HasSuffix(move.Name, "Normals") || strings.HasSuffix(move.Name, "Attacks") {
			if move.Input != "" {
				move.Name = move.Input
			}
		}

		// Extract frame data grids
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

		// extract move description
		infoBody := container.Find(".attack-info-body").Clone()
		infoBody.Find(".frameChart, details, script, style, .details-classic-container").Remove()
		infoBody.Find(".toggle-group-b").Remove()
		infoBody.Find(".tooltiptext, sup.reference").Remove()

		infoBody.Find("a").Each(func(_ int, a *goquery.Selection) {
			a.ReplaceWithHtml(" " + a.Text() + " ")
		})

		var paragraphs []string
		infoBody.Find("p").Each(func(_ int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		})
		
	//If no <p> tags were found, check if there is plain text inside infoBody
		if len(paragraphs) == 0 {
		    // Remove ul elements so notes aren't duplicated as description
		    clone := infoBody.Clone()
		    clone.Find("ul").Remove()
		    rawText := strings.TrimSpace(clone.Text())
		    if rawText != "" {
		        paragraphs = append(paragraphs, rawText)
		    }
		}

		move.Description = strings.Join(paragraphs, "\n")

		// Extract any wiki notes/bullet points
		infoBody.Find("> ul > li").Each(func(_ int, li *goquery.Selection) {
			text := strings.TrimSpace(li.Text())
			if text != "" {
				move.Notes = append(move.Notes, text)
			}
		})

		moves = append(moves, move)

	})
	return moves, nil
}

