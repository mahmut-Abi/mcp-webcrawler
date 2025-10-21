# MCP WebCrawler - 开发指南

## 项目概述

这是一个用 Go 编写的模型上下文协议（MCP）服务器实现，为 LLM 提供网页爬取功能。

## 架构

```
cmd/server/
  └── main.go           # 服务器入口点
internal/
  ├── crawl/
  │   ├── types.go      # 数据类型和结构
  │   ├── crawler.go    # 网页获取和解析
  │   └── crawler_test.go # 单元测试
  └── mcp/
      └── server.go     # MCP 服务器实现
```

## 核心组件

### 1. 爬虫（internal/crawl/crawler.go）
- **FetchPage**：下载并解析单个网页
- **extractMainContent**：提取主要文本内容
- **extractLinks**：找到页面上的所有链接
- **extractImages**：提取图片 URL

### 2. MCP 服务器（internal/mcp/server.go）
- 注册 MCP 工具
- 实现 `crawl_page` 工具
- 处理工具调用和响应格式化

### 3. 数据类型（internal/crawl/types.go）
- `PageContent`：带有元数据的完整页面数据
- `CrawlRequest`：请求参数
- `CrawlResponse`：为 LLM 提供的结构化响应

## 构建

```bash
go build -o bin/mcp-webcrawler ./cmd/server
```

## 运行测试

```bash
# 运行单元测试
go test ./internal/crawl -v

# 运行集成测试
go test . -v

# 运行带覆盖率的测试
go test -v -cover ./...
```

## 测试示例

### 直接测试爬虫
```go
// 创建爬虫实例并调用 FetchPage 函数
crawler := crawl.NewCrawler(15 * time.Second)
page, err := crawler.FetchPage(context.Background(), "https://example.com")
```

### 可用的测试命令
```bash
# 测试刬认爬取
go test -run TestCrawlExampleCom -v

# 测试响应格式化
go test -run TestCrawlResponseFormat -v

# 测试爬虫创建
go test -run TestNewCrawler -v
```

## 运行服务器

```bash
./bin/mcp-webcrawler --timeout 15
```

选项：
- `--timeout`：请求超时时间（秒）（默认：15）

## MCP 工具：crawl_page

### 输入参数
- `url`（字符串，必需）：要爬取的 URL
- `render_js`（布尔值，可选）：是否渲染 JavaScript
- `max_content_length`（整数，可选）：最大内容长度

### 响应格式
```json
{
  "title": "页面标题",
  "description": "元描述",
  "content": "提取的页面内容...",
  "links": ["http://..."],
  "fetched_at": "2024-01-01T12:00:00Z"
}
```

## 依赖项

- `github.com/modelcontextprotocol/go-sdk`：MCP 协议实现
- `github.com/PuerkitoBio/goquery`：HTML 解析
- `golang.org/x/net`：网络工具

## 未来改进

1. 添加 JavaScript 渲染支持
2. 实现 crawl_site 工具不递归爬取
3. 添加缓存機制
4. 支持自定义选择器进行内容提取
5. 添加速率限制和域名白名单
6. 实现伪异步爬取带任务跟踪
