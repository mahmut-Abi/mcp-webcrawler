package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"mcp-webcrawler/internal/crawl"
)

func TestCrawlExampleCom(t *testing.T) {
	crawler := crawl.NewCrawler(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, err := crawler.FetchPage(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	t.Logf("Page Title: %s", page.Title)
	t.Logf("Status Code: %d", page.StatusCode)
	t.Logf("Content Length: %d", page.ContentLength)
	t.Logf("Elapsed Time: %.2f seconds", page.ElapsedSeconds)
	t.Logf("Links Found: %d", len(page.Links))
	t.Logf("Content Preview: %.100s...", page.Content)
}

func TestCrawlResponseFormat(t *testing.T) {
	crawler := crawl.NewCrawler(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, err := crawler.FetchPage(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	resp := crawl.CrawlResponse{
		Title:       page.Title,
		Description: page.Description,
		Content:     page.Content,
		Links:       page.Links,
		FetchedAt:   page.FetchedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	t.Logf("Response JSON:\n%s", string(data))
}

// TestEnhancedMetaExtraction tests the enhanced metadata extraction
func TestEnhancedMetaExtraction(t *testing.T) {
	crawler := crawl.NewCrawler(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, err := crawler.FetchPage(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	// Verify that all expected fields are populated
	if page.Title == "" {
		t.Error("Title should not be empty")
	}

	if page.Description == "" {
		t.Logf("Description is empty (this is acceptable for some sites)")
	}

	if len(page.Headings) == 0 {
		t.Logf("No headings found (this is acceptable for some sites)")
	}

	if page.Language != "" {
		t.Logf("Language detected: %s", page.Language)
	}

	t.Logf("Enhanced Metadata:")
	t.Logf("- Keywords: %s", page.Keywords)
	t.Logf("- Author: %s", page.Author)
	t.Logf("- Image URL: %s", page.ImageURL)
	t.Logf("- Language: %s", page.Language)
	t.Logf("- Headings count: %d", len(page.Headings))
}

// TestURLValidation tests URL validation and security
func TestURLValidation(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		should bool
	}{
		{"valid https", "https://example.com", true},
		{"valid http", "http://example.com", true},
		{"localhost", "http://localhost:8080", false},
		{"127.0.0.1", "http://127.0.0.1", false},
		{"invalid protocol", "ftp://example.com", false},
	}

	crawler := crawl.NewCrawler(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, tc := range tests {
		_, err := crawler.FetchPage(ctx, tc.url)
		hasError := err != nil
		expectedError := !tc.should

		if hasError != expectedError {
			t.Errorf("%s: expected error=%v, got error=%v", tc.name, expectedError, hasError)
			if err != nil {
				t.Logf("Error: %v", err)
			}
		} else {
			t.Logf("%s: OK (error=%v)", tc.name, hasError)
		}
	}
}
