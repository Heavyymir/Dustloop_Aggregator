package discovery

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/models"
	"github.com/PuerkitoBio/goquery"
)

// Function to find characters dynamically using URLs using specific `href` HTML elements
func DiscoveredChars(data []byte, gameSlug string) ([]models.Character, error) {
	// Pages to be excluded from lookup
	excluded := map[string]bool{
		"FAQ":                true,
		"HUD":                true,
		"Mechanics":          true,
		"Damage":             true,
		"Frame_Data":         true,
		"Universal_Strategy": true,
		"Esoterica":          true,
		"Miscellaneous":      true,
		"Tier_Lists":         true,
		"Patch_Notes":        true,
		"Getting_Started":	  true,
		"Team_of_3":		  true,	
		
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("No characters found")
	}

	// Read the byte data and create a document for the function to scrape for character names
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Create a map to list duplicate hrefs
	seen := make(map[string]bool)

	var characters []models.Character
	// Find each `href` HTML element in the parsed document
	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		// Use the URL package to parse `href` elements into a usable URL struct
		parsed, err := url.Parse(href)
		if err != nil {
			return
		}

		// Check found `href` elements for the needed '/w/<game>/' format
		// If gameSlug is "UNI2" or "uni2", check if the link path contains the game path:
		prefix := "/w/" + gameSlug + "/"
		if strings.ToLower(gameSlug) == "uni2" {
		    prefix = "/w/Under_Night_In-Birth/UNI2/"
		}
		
		if !strings.HasPrefix(parsed.Path, prefix) {
		    return
		}
		
		slug := strings.TrimPrefix(parsed.Path, prefix)

		// Find character slugs using the parsed path and prefix
		slug = strings.TrimPrefix(parsed.Path, prefix)
		// Return if Character Slug is empty
		if slug == "" || strings.Contains(slug, "/") {
			return
		}

		// Check the created "link" go query for a text element to scrape the character name. Return if empty.
		name := strings.TrimSpace(link.Text())
		if name == "" {
			name = strings.ReplaceAll(slug, "_", " ")
		}

		// Exclude non-Character pages from lookup
		if excluded[slug] {
			return
		}

		// Check the bool value for slugs in seen. If seen
		if seen[slug] {
			return
		}

		// Assign a true value to the slug if not seen to avoid duplicating appends.
		seen[slug] = true	

		// Create and append a Character struct to Characters slice. Assign values that match elements.
		characters = append(characters, models.Character{
			Name: name,
			Slug: slug,
			URL:  parsed.Path,
		})
	})

	return characters, nil
}
