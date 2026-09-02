package api

// Move this file to internal fightdataapi package once testing of get request completed.
import (
	"fmt"
	"io"
	"net/http"
	"context"
	"strings"
	"time"
	//"os"
	//"path/filepath"

	"github.com/chromedp/chromedp"
)

// Get request to scrape for HTML data from wikis. Needs inputs to be updated to be selctable later.
func (c *Client) Fetch(url string) ([]byte, error) {

	// Check to see if the URL is present inside the Cache
	data, ok := c.Cache.Get(url)
	if ok && len(data) > 0 {
		return data, nil
	}

	//Basic request logic using input URL
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Read the HTTP Body. If empty, return an error.
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty response body from %s", url)
	}

	// Check res.StatusCode. If it is above 299, return an error with the status and body of the response.
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf(
		"response failed: status= %s body= %q", 
		res.Status, 
		string(data[:min(len(data), 500)]),
		)
	}

	// Add data from the request to the Cache.
	c.Cache.Add(url, data)

	return data, nil
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


// FetchHTMLHeadless launches a headless browser instance with the url, waits for cloudflare/js to settle
// and extracts the full outer HTML for parsing.
func (c *Client) FetchHTMLHeadless(url string, timeout time.Duration) ([]byte, error) {
	// Check cache to avoid slow browser launches
	if data, ok := c.Cache.Get(url); ok && len(data) > 0 {
		return data, nil
	}
	// Configure chrome options for the instance
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// "Headless" is the flag that hides the UI for the Chromium instance in the function
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("enable-automation", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	// Create the chromedp context
	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// Set the timeout context to prevent requests hanging
	ctx, cancelTimeout := context.WithTimeout(ctx, timeout)
	defer cancelTimeout()

	var htmlContent string

	// Run browser automation tasks
	err := chromedp.Run(ctx,
		// Navigate to target URL
		chromedp.Navigate(url),

		// Poll/Sleep until Cloudflare passes and the title changes
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 15; i++ {
				var currentTitle string
				if err := chromedp.Title(&currentTitle).Do(ctx); err == nil {
					if currentTitle != "" && currentTitle != "Just a moment..." {
						return nil
					}
				}
				time.Sleep(1 * time.Second)
			}
			return nil
		}),
	
		// Give dynamic scripts 1 extra second to render tables
		chromedp.Sleep(1 * time.Second),
	
		// Capture full rendered HTML
		chromedp.OuterHTML("html", &htmlContent),
		)

	if err != nil {
		return nil, fmt.Errorf("chromedp exectution failed for %s: %w", url, err)
	}


	// 5. Verify request didn't capture the challenge page
	if strings.Contains(htmlContent, "<title>Just a moment...</title>") {
		return nil, fmt.Errorf("cloudflare challenge was not passed for %s", url)
	}

	data := []byte(htmlContent)
	if len(data) == 0 {
		return nil, fmt.Errorf("empty html returned from headless session for %s", url)
	}

	// Cache the fetched html
	c.Cache.Add(url, data)

	return data, nil	
}
