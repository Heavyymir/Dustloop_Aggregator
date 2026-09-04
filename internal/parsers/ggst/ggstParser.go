package ggst

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

func GGSTCharPageParser(data []byte) ([]models.Move, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var moves []models.Move

	doc.Find(".attack-container").Each(func(_ int, container *goquery.Selection) {
		move := models.Move{
			FrameDataGrids: []models.FrameDataGrid{},
			Notes:          []string{},
		}

		// Gather move name
		heading := container.PrevAll().Filter("div.mw-heading3, h3").First()
		if heading.Length() > 0 {
			rawName := heading.Find("h3").Text()
			if rawName == "" {
				rawName = heading.Text()
			}
			move.Name = strings.Join(strings.Fields(rawName), " ")
		}

		// Gather input notation
		prevNode := container.Prev()
		inputBadge := prevNode.Find(".input-badge")
		if inputBadge.Length() > 0 {
			move.Input = strings.Join(strings.Fields(inputBadge.Text()), " ")
		} else {
			move.Input = move.Name
		}

		// Extract Grids (handles single-row, multi-row, notes grids)
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
				Headers: headers,
				Rows:    []models.FrameDataRow{},
			}

			gridNode.Find(".frameDataGridRow").Each(func(_ int, rowNode *goquery.Selection) {
				frameRow := models.FrameDataRow{
					Cells: []models.Cell{},
				}

				rowNode.ChildrenFiltered("div").Each(func(_ int, cellNode *goquery.Selection) {
					frameRow.Cells = append(frameRow.Cells, parseCell(cellNode))
				})

				grid.Rows = append(grid.Rows, frameRow)
			})

			move.FrameDataGrids = append(move.FrameDataGrids, grid)
		})

		// Extract Description
		var paragraphs []string
		container.Find(".attack-info-body > p").Each(func(_ int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		})
		move.Description = strings.Join(paragraphs, "\n")

		// Extract Notes / Bullet points
		container.Find(".attack-info-body > ul > li").Each(func(_ int, li *goquery.Selection) {
			text := strings.TrimSpace(li.Text())
			if text != "" {
				move.Notes = append(move.Notes, text)
			}
		})

		moves = append(moves, move)
	})

	return moves, nil
}

func parseCell(cell *goquery.Selection) models.Cell {
	tooltip := strings.TrimSpace(cell.Find(".tooltiptext").Text())

	visible := cell.Clone()
	visible.Find(".tooltiptext").Remove()

	return models.Cell{
		Value:   strings.TrimSpace(visible.Text()),
		Tooltip: tooltip,
	}
}
