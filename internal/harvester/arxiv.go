package harvester

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

type ArxivFetcher struct{}

type ArxivFeed struct {
	Entry ArxivEntry `xml:"entry"`
}

type ArxivEntry struct {
	Title   string   `xml:"title"`
	Summary string   `xml:"summary"`
	Author  []string `xml:"author>name"`
	Link    []struct {
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	} `xml:"link"`
}

func (f *ArxivFetcher) Fetch(arxivID string) (string, string, error) {
	// 1. Fetch Metadata via Arxiv API
	apiURL := fmt.Sprintf("http://export.arxiv.org/api/query?id_list=%s", arxivID)
	body, _, err := utils.GetWithRetry(apiURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch metadata: %v", err)
	}

	var feed ArxivFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return "", "", fmt.Errorf("failed to decode XML: %v", err)
	}

	entry := feed.Entry
	if entry.Title == "" {
		return "", "", fmt.Errorf("paper not found or empty title")
	}

	// 2. Attempt to fetch and parse PDF (optional but nice)
	pdfText := ""
	pdfURL := ""
	for _, link := range entry.Link {
		if link.Type == "application/pdf" {
			pdfURL = link.Href
			break
		}
	}

	if pdfURL != "" {
		// Replace http with https if needed
		pdfURL = strings.Replace(pdfURL, "http://", "https://", 1)
		pdfText, _ = f.extractPDFText(pdfURL)
	}

	// 3. Build Markdown
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", strings.TrimSpace(entry.Title)))
	sb.WriteString(fmt.Sprintf("**Authors:** %s\n\n", strings.Join(entry.Author, ", ")))
	sb.WriteString("## Abstract\n\n")
	sb.WriteString(strings.TrimSpace(entry.Summary))
	sb.WriteString("\n\n---\n\n")

	if pdfText != "" {
		sb.WriteString("## PDF Content Snippet / Full Text\n\n")
		sb.WriteString(pdfText)
	} else {
		sb.WriteString("> *Note: Full PDF text extraction was not available or failed.*\n")
	}

	return entry.Title, sb.String(), nil
}

func (f *ArxivFetcher) extractPDFText(pdfURL string) (string, error) {
	body, _, err := utils.GetWithRetry(pdfURL, nil)
	if err != nil {
		return "", err
	}

	r, err := pdf.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
