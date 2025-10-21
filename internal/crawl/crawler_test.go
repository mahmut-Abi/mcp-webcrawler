package crawl

import (
	"context"
	"testing"
	"time"
)

func TestNewCrawler(t *testing.T) {
	crawler := NewCrawler(10 * time.Second)
	if crawler == nil {
		t.Fatal("Crawler should not be nil")
	}
}

func TestFetchPage(t *testing.T) {
	crawler := NewCrawler(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 使用 example.com 进行测试
	page, err := crawler.FetchPage(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	if page.URL != "https://example.com" {
		t.Errorf("Expected URL https://example.com, got %s", page.URL)
	}

	if page.Title == "" {
		t.Error("Expected title, got empty string")
	}

	if page.Content == "" {
		t.Error("Expected content, got empty string")
	}
}
