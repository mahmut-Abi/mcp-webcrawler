package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mcp-webcrawler/internal/crawl"
	"mcp-webcrawler/internal/mcp"
)

func main() {
	var (
		timeout = flag.Int("timeout", 15, "Request timeout in seconds")
	)
	flag.Parse()

	log.Println("MCP WebCrawler Server Started")
	log.Printf("Timeout: %d seconds", *timeout)

	// 创建爬取器
	crawler := crawl.NewCrawler(time.Duration(*timeout) * time.Second)

	// 创建 MCP 服务器
	mcpServer, err := mcp.NewMCPServer(crawler)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// 创建带取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理优雅的關閉
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	// 启动服务器
	if err := mcpServer.Serve(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
