package api

import(
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Heavyymir/Dustloop_Aggregator/catalog"
)


// FetchMediaWikiPage queries MediaWiki based domains' api.php
// and returns the raw rendered HTML bytes of the page body
func (c *Client) FetchMediaWikiPage(apiURL, pageTitle string) ([]byte, error) {

	// Safetly construct the query paramaters
	params := url.Values{}
	params.Set("action", "parse")
	params.Set("page", pageTitle)
	params.Set("format", "json")
	params.Set("prop", "text")
	params.Set("disablelimitreport", "1")
	params.Set("disableeditsection", "1")

	// Construct the full URL by encoding the paramaters
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Check the internal cache for the request first
	if data, ok := c.Cache.Get(fullURL); ok && len(data) > 0 {
		return data, nil
	}

	// Build the HTTP request
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create mediawiki request: %w", err)
	}

	// Set user-agent headers to comply with MediaWiki requirements
	req.Header.Set("User-Agent", "CharDataAggregator/1.0 (https://github.com/Heavyymir/CharData_Aggregator)")
	req.Header.Set("Accept", "application/json")

	// Catch the http response
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do mediawiki request: %w", err)
	}

	defer res.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read mediawiki response: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("mediawiki request failed: status=%w", err)
	}

	// Initalise a MediaWikiParseResponse struct to hold the unmarshalled JSON data from the API request
	var parsedResponse MediaWikiParseResponse

	// Unmarshal the HTML into the Response struct 
	if err := json.Unmarshal(bodyBytes, &parsedResponse); err != nil {
		return nil, fmt.Errorf("unmarshal mediawiki json: %w", err)
	}
	
	if parsedResponse.Error != nil {
		return nil, fmt.Errorf("mediawiki api error (%s): %s", parsedResponse.Error.Code, parsedResponse.Error.Info)
	}

	// Assign the content in the response struct to a []byte slice for parser functions
	htmlContent := []byte(parsedResponse.Parse.Text.HTML)
	if len(htmlContent) == 0 {
		return nil, fmt.Errorf("no content returned for page: %s", pageTitle)
	}

	c.Cache.Add(fullURL, htmlContent)

	return htmlContent, nil
}


// Resolve page title builds the MediaWiki page title using catalog paths
func ResolvePageTitle(game catalog.Game, characterName string) string {
	formattedName := url.PathEscape(characterName)
	return fmt.Sprintf(game.CharacterPath, formattedName)
}
