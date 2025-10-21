# MCP WebCrawler - 构建与测试指南

## 快速构建

```bash
# 安装依赖
go mod download

# 构建二进制文件
mkdir -p bin
go build -o bin/mcp-webcrawler ./cmd/server

# 验证二进制文件
./bin/mcp-webcrawler -h
```

## 运行测试

### 所有测试
```bash
go test -v ./...
```

### 特定测试

#### 单元测试（不需要网络）
```bash
go test -v -short ./internal/crawl -run TestNewCrawler
```

#### 集成测试（需要网络）
```bash
go test -v -short ./... -run TestCrawl
```

#### 带覆盖率
go test -v -cover ./...
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 测试输出示例

```
=== RUN   TestCrawlExampleCom
    examples_test.go:22: 页面标题: Example Domain
    examples_test.go:23: 状态码: 200
    examples_test.go:24: 内容长度: 513
    examples_test.go:25: 执行时间: 0.51 秒
    examples_test.go:26: 找到的链接: 1
    examples_test.go:27: 内容预览: Example Domain...
--- PASS: TestCrawlExampleCom (0.52s)

=== RUN   TestCrawlResponseFormat
    examples_test.go:53: 响应 JSON：
        {
          "title": "Example Domain",
          "description": "",
          "content": "Example Domain\nThis domain is for use in documentation examples without needing permission. Avoid use in operations.",
          "links": [
            "https://iana.org/domains/example"
          ],
          "fetched_at": "2025-10-21T04:09:25Z"
        }
--- PASS: TestCrawlResponseFormat (0.18s)
PASS
ok  	mcp-webcrawler	0.699s
```

## 启动服务器

### 标准输入/输出模式（用于 MCP 客户端）
```bash
./bin/mcp-webcrawler --timeout 20
```

服务器将通过标准输入/输出连接并等待 MCP 客户端请求。

### 带日志输出
```bash
./bin/mcp-webcrawler --timeout 20 2>&1 | tee server.log
```

## 调试

### 启用详细输出
```bash
DEBUG=1 go test -v ./...
```

### 测试特定 URL
创建测试文件 `test_url.go`:
```go
package main

import (
	"context"
	"fmt"
	"time"

	"mcp-webcrawler/internal/crawl"
)

func main() {
	// 创建爬虫实例
	crawler := crawl.NewCrawler(15 * time.Second)
	// 测试网页获取
	page, err := crawler.FetchPage(context.Background(), "https://example.com")
	if err != nil {
		fmt.Printf("错误: %v", err)
		return
	}
	fmt.Printf("标题: %s\n", page.Title)
	fmt.Printf("内容长度: %d\n", len(page.Content))
	fmt.Printf("链接: %v\n", page.Links)
}
```

然后：
```bash
go run test_url.go
```

## 项目文件

| 文件 | 用途 |
|------|----------|
| `cmd/server/main.go` | 服务器入口点 |
| `internal/crawl/types.go` | 数据结构 |
| `internal/crawl/crawler.go` | 网页获取逻辑 |
| `internal/crawl/crawler_test.go` | 单元测试 |
| `internal/mcp/server.go` | MCP 服务器实现 |
| `examples_test.go` | 集成测试 |
| `go.mod` / `go.sum` | 依赖 |
| `bin/mcp-webcrawler` | 编译的二进制文件 |

## 依赖项

- `github.com/modelcontextprotocol/go-sdk` - MCP 协议
- `github.com/PuerkitoBio/goquery` - HTML 解析
- `github.com/google/jsonschema-go` - JSON 模式验证
- `golang.org/x/net` - 网络工具

## 构建统计数据

- 二进制大小: 约30 MB（静态链接）
- 构建时间: <1 秒
- 支持操作系统: Linux, macOS, Windows
- Go 版本: 1.21+

## 故障排除

### 不能构建二进制文件
```bash
# 清理依赖
go clean -cache
go mod tidy
go mod download

# 重新构建
go build -o bin/mcp-webcrawler ./cmd/server
```

### 测试超时
- 确保你有网络连接
- 增加超时: `go test -timeout 30s ./...`
- 使用 `-short` 标志跳过网络测试

### 导入错误
```bash
go mod tidy
go get -u ./...
```
