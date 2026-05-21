package harvester

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-shiori/go-readability"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

type WebFetcher struct{}

func (f *WebFetcher) Fetch(targetURL string) (string, string, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %v", err)
	}

	body, _, err := utils.GetWithRetry(targetURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch URL: %v", err)
	}

	// 1. Extract main content using readability
	article, err := readability.FromReader(bytes.NewReader(body), parsedURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse readability: %v", err)
	}

	// 2. Convert HTML to Markdown
	// Note: go-readability provides article.Content which is the HTML of the main body
	markdown, err := htmltomarkdown.ConvertString(article.Content)
	if err != nil {
		// Fallback to plain text if markdown conversion fails
		return article.Title, article.TextContent, nil
	}

	return article.Title, markdown, nil
}
