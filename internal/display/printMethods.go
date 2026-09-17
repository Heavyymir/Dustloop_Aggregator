package display

import(
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
)


func PrintGrid(grid models.FrameDataGrid) {
	if len(grid.Headers) == 0 {
		return
	}

	// Detect columns for colour coding
	advantageCols := make(map[int]bool)
	for idx, h := range grid.Headers {
		norm := strings.ToLower(strings.TrimSpace(h))
		if strings.Contains(norm, "block") ||
			strings.Contains(norm, "whiff") ||
			strings.Contains(norm, "adv") ||
			norm == "on hit" || 
			norm == "on-hit" ||
 			strings.HasPrefix(norm, "hit adv") {
			advantageCols[idx] = true
		}
	}

	// Calculate the maximum width required for a column
	widths := make([]int, len(grid.Headers))

	// Check the length of headers
	for i, header := range grid.Headers {
		widths[i] = utf8.RuneCountInString(header)
	}

	// Check length of cell values across all rows
	for _, row := range grid.Rows {
		for i, cell := range row.Cells {
			if i < len(widths) {
				cellLen := utf8.RuneCountInString(cell.Value)
				if cellLen > widths[i] {
					widths[i] = cellLen	
				}
			}
		}
	}

	// Add spaces to pad the printed grid
	for i := range widths {
		widths[i] += 2
	}

	// Print headers
	for i, header := range grid.Headers {
		fmt.Print(padRight(header, widths[i]))
	}
	fmt.Println()
	
	// Print a seperator line under the headers
	totalWidth := 0
	for _, w := range widths {
		totalWidth += w
	}

	fmt.Println(strings.Repeat("-", totalWidth))

	// Print data rows
	for _, row := range grid.Rows {
		for i, cell := range row.Cells {
			if i < len(widths) {
				padded := padRight(cell.Value, widths[i])

				if advantageCols[i] {
					fmt.Print(colouriseAdvantage(cell.Value) + strings.Repeat(" ", widths[i]-utf8.RuneCountInString(cell.Value)))
				} else {
					fmt.Print(padded)
				}
			}
		}
		fmt.Println()
	}
}


// padRight pads a string with a trailing spaces to match targetWidth
func padRight(s string, targetWidth int) string {
	length := utf8.RuneCountInString(s)
	if length >= targetWidth {
		return s
	}
	return s + strings.Repeat(" ", targetWidth - length)
}


