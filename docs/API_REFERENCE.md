# MCP WebCrawler - API 参考

## 工具：crawl_page

获取并上解析网页为结构化文本和元数据。

### 请求

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "crawl_page",
    "arguments": {
      "url": "https://example.com",
      "render_js": false,
      "max_content_length": 10000
    }
  }
}
```

### 输入参数

| 参数 | 类型 | 是否必需 | 描述 | 默认 |
|-----------|------|----------|-------------|----------|
| url | 字符串 | 是 | 要爬取的 URL | - |
| render_js | 布尔值 | 否 | 是否渲染 JavaScript | false |
| max_content_length | 整数 | 否 | 最大内容长度（字符） | 10000 |

### 响应

#### 成功响应 (2xx)

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "title": "示例域名",
    "description": "示例域名。由 IANA 注册事残使用。",
    "content": "示例域名
此域名事残使用方案示例...",
    "links": [
      "https://www.iana.org/domains/example"
    ],
    "fetched_at": "2025-10-21T04:09:25Z"
  }
}
```

#### 错误响应 (4xx/5xx)

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "内部错误",
    "data": {
      "error": "获取页面失败: dial tcp: lookup example.invalid: 没有该主机"
    }
  }
}
```

### 输出字段

| 字段 | 类型 | 描述 |
|-------|------|-------------|
| title | 字符串 | 页面标题（来自 `<title>` 标签） |
| description | 字符串 | 元描述（来自 `<meta name=description>`） |
| content | 字符串 | 提取的主要文本内容 |
| links | 数组[字符串] | 页面上找到的所有 URL |
| fetched_at | 字符串 | 获取时的 ISO 8601 时间戳 |

### HTTP 状态码映射

该工具不直接返回 HTTP 错误。相反：
- **200 OK**: 正常成功
- **301/302 重定向**: 自动跟宏
- **404 未找到**: 返回错误，并在日志中记录状态码
- **403 禁止**: 返回错误
- **500 服务器错误**: 返回错误
- **连接超时**: 返回超时错误

### 错误情况

#### 无效的 URL
```
错误: URL 参数是必需的且必须为字符串
```

#### 网络无法达到
```
错误: 获取页面失败: dial tcp: lookup example.invalid: 没有该主机
```

#### 请求超时
```
错误: 上下文截限时间已超显
```

#### 无效的 HTML
爬虫正常处理无效的 HTML 并返回所有可以提取的内容。

## 示例

### 基本使用

```bash
curl -X POST http://localhost:3000/tools/call   -H "Content-Type: application/json"   -d '{
    "name": "crawl_page",
    "arguments": {
      "url": "https://example.com"
    }
  }'
```

### Python 客户端

```python
import json
import subprocess

server = subprocess.Popen(
    ['./bin/mcp-webcrawler'],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE
)

request = {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
        "name": "crawl_page",
        "arguments": {"url": "https://example.com"}
    }
}

server.stdin.write(json.dumps(request).encode() + b'\n')
response_line = server.stdout.readline()
response = json.loads(response_line)

print("Title:", response['result']['title'])
```

## 速率限制和调扖

服务器当前不是内置速率限制。实现客户端侧调昆：

```go
// 在请求之间添加延迟
 time.Sleep(time.Second)
crawler.FetchPage(ctx, url)
```

## 超时

通过命令行会标志配置：

```bash
./bin/mcp-webcrawler --timeout 30  # 30 秒
```

默认每个请求 15 秒。

## 内容提取策略

爬虫使用以下优先级顺序进行内容提取：

1. 找到 `<article>`、`<main>` 或 `[role=main]` 部分
2. 如果找不到，搜索常见的内容容器类
3. 提取所有 `<h1>-<h6>` 和 `<p>` 元素
4. 移除噪声: `<script>`、`<style>`、`<nav>`、`<footer>`
5. 去重和删除空白

## 链接提取

- 提取所有 `<a href>` 元素
- 保持绝对和相对 URL
- 去重 URL
- 不验证链接目标

## 图片提取

- 提取所有 `<img src>` 元素
- 按原样返回图片 URL（可能是相对路径）
- 有助于检测页面资源

## 性能笔记

- 典制响应时间: 0.5-2 秒/页
- 每个请求内存: 1-10 MB（取决于页面大小）
- 并发请求: 受 go-sdk 传输求限制
- 最大檄记页面大小: 50 MB

## 模式定义

### 输入模式

```json
{
  "type": "object",
  "properties": {
    "url": {
      "type": "string",
      "description": "要爬取的 URL"
    },
    "render_js": {
      "type": "boolean",
      "description": "是否渲染 JavaScript"
    },
    "max_content_length": {
      "type": "integer",
      "description": "最大内容长度（字符）"
    }
  },
  "required": ["url"],
  "additionalProperties": false
}
```

### 输出模式

```json
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "页面标题"
    },
    "description": {
      "type": "string",
      "description": "元描述"
    },
    "content": {
      "type": "string",
      "description": "提取的页面内容"
    },
    "links": {
      "type": "array",
      "items": {"type": "string"},
      "description": "页面上找到的链接"
    },
    "fetched_at": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 时间戳"
    }
  },
  "additionalProperties": false
}
```
