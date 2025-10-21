package crawl

import "time"

// PageContent 表示从网页中提取的内容
type PageContent struct {
	URL            string    `json:"url"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Keywords       string    `json:"keywords"`
	Author         string    `json:"author"`
	ImageURL       string    `json:"image_url"`
	Content        string    `json:"content"`
	Links          []string  `json:"links"`
	Images         []string  `json:"images"`
	Headings       []string  `json:"headings"`
	StatusCode     int       `json:"status_code"`
	FetchedAt      time.Time `json:"fetched_at"`
	ContentLength  int       `json:"content_length"`
	ElapsedSeconds float64   `json:"elapsed_seconds"`
	Language       string    `json:"language"`
}

// CrawlRequest 包含爬取操作的参数
type CrawlRequest struct {
	URL              string
	RenderJS         bool
	MaxContentLength int
	ReturnFormat     string // json 或 markdown
	FollowRedirects  bool
	TimeoutSeconds   int
}

// CrawlResponse 是爬取操作的响应
type CrawlResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Links       []string `json:"links"`
	FetchedAt   string   `json:"fetched_at"`
}
