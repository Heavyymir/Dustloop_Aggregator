package dustloop

import(
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)

// Helper for Cell Parsing for game parsers
func parseCell(cell *goquery.Selection) models.Cell {
	tooltip := strings.TrimSpace(cell.Find(".tooltiptext").Text())

	visible := cell.Clone()
	visible.Find(".tooltiptext").Remove()

	return models.Cell{
		Value:   strings.TrimSpace(visible.Text()),
		Tooltip: tooltip,
	}
}
