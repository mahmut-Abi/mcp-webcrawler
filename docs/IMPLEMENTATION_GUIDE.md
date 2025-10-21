# MCP WebCrawler - 实现指南

## 概述

本指南說明了 MCP WebCrawler 服务器的实现方式以及如何扩展它。

## 架构

```
┌───────────────────────────────────────────────┐
│      MCP 客户端 (Claude, GPT, 等)│
└─────────────┬───────────────────────────┘
                   │ stdio/HTTP
                   ╵
┌───────────────────────────────────────────────┐
│         MCP 服务器 (mcp-webcrawler)│
│  ┌───────────────────────────────────┐   │
│  │  工具处理: crawl_page                │   │
│  │  - 输入验证                         │   │
│  │  - URL 解析                      │   │
│  │  - 超时管理                       │   │
│  └─────────────┬─────────────────────┘   │
│                   │                                 │
│  ┌─────────────╵─────────────────────┐   │
│  │  爬虫 (internal/crawl)│   │
│  │  - HTTP 请求                            │   │
│  │  - HTML 解析 (goquery)│   │
│  │  - 内容提取                       │   │
│  │  - 链接/图片提取                   │   │
│  └─────────────┬─────────────────────┘   │
│                   │                                 │
└─────────────┬─────────────────────┘
                    │
                    ╵
            ┌───────────┌
            │   HTTP GET   │
            │  目标 URL  │
            └───────────┘
```

## 核心组件

### 1. 数据类型 (internal/crawl/types.go)

#### PageContent
代表完整的提取页面数据:
```go
type PageContent struct {
    URL            string        // 原始 URL
    Title          string        // 页面标题
    Description    string        // 元描述
    Content        string        // 提取的文本内容
    Links          []string      // 找到的所有链接
    Images         []string      // 找到的所有图片
    StatusCode     int           // HTTP 状态码
    FetchedAt      time.Time     // 获取时间戳
    ContentLength  int           // 响应大小
    ElapsedSeconds float64       // 请求时间
}
```

#### CrawlRequest/CrawlResponse

工具调用的标准请求/响应格式。

### 2. 爬虫 (internal/crawl/crawler.go)

#### NewCrawler()
使用配置的超时创建网页爬虫实例:
```go
// 创建爬虫实例並设置超时时间
crawler := crawl.NewCrawler(15 * time.Second)
```

#### FetchPage(ctx, url)
主要函数，执行以下操作:
1. 创建 HTTP 请求并设置正常头
2. 设置 User-Agent 以不被阻止
3. 获取页面内容
4. 使用 goquery 解析 HTML
5. 提取主要内容
6. 收集链接和图片
7. 返回结构化的 PageContent

#### 内容提取函数

**extractMainContent()**
- 移除噪声 (script, style, nav, footer 标签)
- 提取 article/main/content 部分
- 需要时提取身体段落和标题
- 返回去重的文本

**extractLinks()**
- 找到所有 `<a href>` 元素
- 去重 URL
- 返回链接列表

**extractImages()**
- 找到所有 `<img src>` 元素
- 返回图片 URL

### 3. MCP 服务器 (internal/mcp/server.go)

#### 输入/输出类型

```go
type CrawlPageInput struct {
    URL              string // 必需
    RenderJS         bool   // 可选
    MaxContentLength int    // 可选
}

type CrawlPageOutput struct {
    Title       string
    Description string
    Content     string
    Links       []string
    FetchedAt   string
}
```

#### 工具注册

服务器使用 go-sdk 的类型化工具 API:
```go
// 模式会一自动从结构体标签生成
inputSchema, _ := jsonschema.For[CrawlPageInput](nil)
outputSchema, _ := jsonschema.For[CrawlPageOutput](nil)

// 注册处理程序
```

## 扩展指南

### 添加 JavaScript 渲染

要添加 JavaScript 支持:

1. 安装 go-rod:
```bash
go get github.com/go-rod/rod
```

2. 添加浏览器选项:
```go
type CrawlPageInput struct {
    URL      string
    RenderJS bool // 启用 JS 渲染
}
```

3. 在 FetchPage() 中实现:
```go
if renderJS {
    // 使用 rod 启动浏览器
    // 执行 JavaScript
    // 获取渲染后的 HTML
}
```

### 缓存

添加结果缓存:

1. 使用 sync.Map 或 redis:
```go
type MCPServer struct {
    server *mcp.Server
    crawler *crawl.Crawler
    cache map[string]*PageContent // 添加缓存
}
```

2. 跳过获取前检查缓存:
```go
if cached, exists := ms.cache[url]; exists {
    return cached
}
```

### 速率限制

添加基于时间的调歡:
```go
import "golang.org/x/time/rate"

type Crawler struct {
    client *http.Client
    limiter *rate.Limiter
}

func (c *Crawler) FetchPage(ctx context.Context, url string) {
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    // 获取...
}
```

## 测试

### 单元测试

孤立测试各个函数:
```go
func TestNewCrawler(t *testing.T) {
    crawler := NewCrawler(10 * time.Second)
    if crawler == nil {
        t.Fatal("爬虫不应为 nil")
    }
}
```

### 集成测试

端到端测试功能:
```go
func TestCrawlExampleCom(t *testing.T) {
    crawler := NewCrawler(15 * time.Second)
    page, err := crawler.FetchPage(context.Background(), "https://example.com")
    // 断言...
}
```

### 伪造测试

不用网络测试:
```go
type MockCrawler struct {
    Pages map[string]*PageContent
}

func (m *MockCrawler) FetchPage(ctx context.Context, url string) (*PageContent, error) {
    if p, ok := m.Pages[url]; ok {
        return p, nil
    }
    return nil, errors.New("未找到")
}
```

## 性能上的考虑

### 并发请求

MCP 框架自动处理并发。多个工具调用可以并行进行。

### 内存使用

- 每个请求将整个 HTML 加载到内存
- 大页面 (>10MB) 可能會引起问题
- 考虑未来的流传输优化

### 超时管理

- 默认每个请求 15 秒
- 可通过标志配置: `--timeout 30`
- 防止资源耗尽

