package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-webcrawler/internal/crawl"
)

// CrawlPageInput 表示 crawl_page 工具的输入参数
type CrawlPageInput struct {
	URL              string `json:"url" jsonschema:"The URL to crawl"`
	RenderJS         bool   `json:"render_js,omitempty" jsonschema:"Whether to render JavaScript"`
	MaxContentLength int    `json:"max_content_length,omitempty" jsonschema:"Maximum content length in characters"`
}

// CrawlPageOutput 表示 crawl_page 工具的输出
type CrawlPageOutput struct {
	Title       string   `json:"title" jsonschema:"Page title"`
	Description string   `json:"description" jsonschema:"Meta description"`
	Content     string   `json:"content" jsonschema:"Extracted page content"`
	Links       []string `json:"links" jsonschema:"Links found on the page"`
	FetchedAt   string   `json:"fetched_at" jsonschema:"Timestamp when the page was fetched"`
}

// MCPServer 币包装 MCP 服务器并添加网页爬取工具
type MCPServer struct {
	server  *mcp.Server
	crawler *crawl.Crawler
}

// NewMCPServer 创建一个新的带有网页爬取功能的 MCP 服务器
func NewMCPServer(crawler *crawl.Crawler) (*MCPServer, error) {
	// 创建实现
	impl := &mcp.Implementation{
		Name:    "mcp-webcrawler",
		Version: "0.1.0",
	}

	// 创建服务器
	s := mcp.NewServer(impl, nil)

	ms := &MCPServer{
		server:  s,
		crawler: crawler,
	}

	// 注册工具
	if err := ms.registerTools(); err != nil {
		return nil, err
	}

	return ms, nil
}

func (ms *MCPServer) registerTools() error {
	// 为 crawl_page 工具创建输入模式
	inputSchema, err := jsonschema.For[CrawlPageInput](nil)
	if err != nil {
		return fmt.Errorf("failed to create input schema: %w", err)
	}

	// 创建输出模式
	outputSchema, err := jsonschema.For[CrawlPageOutput](nil)
	if err != nil {
		return fmt.Errorf("failed to create output schema: %w", err)
	}

	// 注册带有处理程序的工具
	mcp.AddTool(ms.server, &mcp.Tool{
		Name:         "crawl_page",
		Description:  "Fetch and parse a web page into structured text and metadata",
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
	}, ms.crawlPageHandler)

	log.Println("Tool 'crawl_page' registered successfully")
	return nil
}

func (ms *MCPServer) crawlPageHandler(ctx context.Context, req *mcp.CallToolRequest, input CrawlPageInput) (*mcp.CallToolResult, CrawlPageOutput, error) {
	if input.URL == "" {
		return nil, CrawlPageOutput{}, fmt.Errorf("URL is required")
	}

	// 获取页面
	page, err := ms.crawler.FetchPage(ctx, input.URL)
	if err != nil {
		return nil, CrawlPageOutput{}, fmt.Errorf("failed to fetch page: %w", err)
	}

	// 格式化响应
	output := CrawlPageOutput{
		Title:       page.Title,
		Description: page.Description,
		Content:     page.Content,
		Links:       page.Links,
		FetchedAt:   page.FetchedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return nil, output, nil
}

// Serve 使用 stdio 传输启动 MCP 服务器
func (ms *MCPServer) Serve(ctx context.Context) error {
	log.Println("Starting MCP WebCrawler server on stdio...")

	// 创建 stdio 传输
	stdioTransport := &mcp.StdioTransport{}

	// 将服务器连接到传输
	_, err := ms.server.Connect(ctx, stdioTransport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect server to transport: %w", err)
	}

	log.Println("MCP Server started successfully")
	return nil
}
