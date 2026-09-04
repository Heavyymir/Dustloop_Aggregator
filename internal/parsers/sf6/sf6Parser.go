package sf6

import (
	"bytes"
	"strings"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

func SF6CharPageParser(data []byte) ([]models.Move, error) {
	// 
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	
	// Initialise a slice of Move structs
	var moves []models.Move

	doc.Find(".movedata-flex-framedata").Each(func(_ int, block *goquery.Selection) {
		move := models.Move{
			FrameData: make(map[string]models.Cell),
			Headers: []string{},
		}

		names := block.Find(".movedata-flex-framedata-name-item")
		if names.Length() > 0 {
			move.Input = strings.TrimSpace(names.Eq(0).Text())
		}

		if names.Length() > 1 {
			move.Name = strings.TrimSpace(names.Eq(1).Text())
		}

		headers := block.Find("table tr").First().Find("th")
		values := block.Find("table tr").Eq(1).Find("td")

		headers.Each(func(i int, header *goquery.Selection) {
			if i >= values.Length() {
				return
			}

			key := strings.TrimSpace(header.Text())
			move.Headers = append(move.Headers, key)
			move.FrameData[key] = models.Cell{
				Value: strings.TrimSpace(values.Eq(i).Text()),
			}
		})

		moves = append(moves, move)
	})
	return moves, nil 
}
