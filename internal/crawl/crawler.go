package crawl

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Crawler 处理网页拖取和解析
type Crawler struct {
	client  *http.Client
	timeout time.Duration
}

// NewCrawler 创建一个新的爬取器实例
func NewCrawler(timeout time.Duration) *Crawler {
	return &Crawler{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// FetchPage 下载并解析一个网页
func (c *Crawler) FetchPage(ctx context.Context, url string) (*PageContent, error) {
	start := time.Now()

	// Validate URL safety
	if err := validateURL(url); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置用户代理
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 执行请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// 读取内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// 提取内容
	title := doc.Find("title").First().Text()
	description, _ := doc.Find("meta[name=description]").Attr("content")
	// 提取附加元数据
	keywords, _ := doc.Find("meta[name=keywords]").Attr("content")
	author, _ := doc.Find("meta[name=author]").Attr("content")
	ogImage, _ := doc.Find("meta[property=og:image]").Attr("content")
	if ogImage == "" {
		ogImage, _ = doc.Find("meta[name=image]").Attr("content")
	}
	language, _ := doc.Find("html").Attr("lang")

	// 提取主要内容
	content := extractMainContent(doc)
	
	// 清洗内容
	if c.cleanContent {
		content = cleanContent(content)
	}

	// 提取链接
	links := extractLinks(doc, url)

	// 提取图片
	images := extractImages(doc, url)

	return &PageContent{
		URL:            url,
		Title:          title,
		Description:    description,
		Keywords:       keywords,
		Author:         author,
		ImageURL:       ogImage,
		Content:        content,
		Links:          links,
		Images:         images,
		Headings:       extractHeadings(doc),
		StatusCode:     resp.StatusCode,
		FetchedAt:      time.Now(),
		ContentLength:  len(body),
		Language:       language,
		ElapsedSeconds: time.Since(start).Seconds(),
	}, nil
}

// Add Headings field extraction for later use

func extractMainContent(doc *goquery.Document) string {
	var content strings.Builder

	// 移除脚本和样式元素
	doc.Find("script, style, nav, footer").Remove()

	// 提取段落和标题
	doc.Find("article, main, [role=main], .content, .main").Each(func(_ int, s *goquery.Selection) {
		s.Find("h1, h2, h3, h4, h5, h6, p").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				content.WriteString(text + "\n")
			}
		})
	})

	if content.Len() == 0 {
		doc.Find("body").Find("p, h1, h2, h3").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 10 {
				content.WriteString(text + "\n")
			}
		})
	}

	return strings.TrimSpace(content.String())
}

func extractLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	seen := make(map[string]bool)

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href != "" && !seen[href] {
			links = append(links, href)
			seen[href] = true
		}
	})

	return links
}

func extractImages(doc *goquery.Document, baseURL string) []string {
	var images []string
	seen := make(map[string]bool)

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src != "" && !seen[src] {
			images = append(images, src)
			seen[src] = true
		}
	})

	return images
}

// validateURL checks if URL is safe to access
func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check protocol
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported protocol: %s", parsedURL.Scheme)
	}

	// Check for localhost
	hostname := parsedURL.Hostname()
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1" {
		return fmt.Errorf("local host access is not allowed")
	}

	// Check for private IP
	if ip := net.ParseIP(hostname); ip != nil {
		if ip.IsPrivate() || ip.IsLoopback() {
			return fmt.Errorf("private IP access not allowed: %s", hostname)
		}
	}

	return nil
}

// extractHeadings extracts all headings from the page
func extractHeadings(doc *goquery.Document) []string {
	var headings []string

	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			headings = append(headings, text)
		}
	})

	return headings
}
