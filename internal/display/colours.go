package display

import (
	"strings"
)

const (
	colourReset 	= "\033[0m"
	colourRed 		= "\033[31m"
	colourYellow	= "\033[33m"
	colourGreen		= "\033[32m"
)

// Colourise advantage colours positive frames green, neutral yellow and negative red for table and grid prints
func colouriseAdvantage(adv string) string {
	adv = strings.TrimSpace(adv)
	if adv == "" {
		return adv
	}

	upper := strings.ToUpper(adv)
	// handle positive advantage (i.e. +1)
	if strings.HasPrefix(adv, "+") || strings.Contains(upper, "KD") {
		return colourGreen + adv + colourReset
	}
	
	// handle negative advantage (i.e. -1)
	if strings.HasPrefix(adv, "-") {
		return colourRed + adv + colourReset
	}

	// Handle neutral advantage (0)
	if adv == "0" || adv == "±0" {
		return colourYellow + adv + colourReset
	}

	return adv
}
