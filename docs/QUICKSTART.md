# MCP WebCrawler - 快速开始指南

## 5 分钟设置

### 1. 前置要求
- Go 1.21 或更高版本
- Git
- 网络连接（用于初始依赖下载）

### 2. 克隆並设置

```bash
cd mcp-webcrawler
go mod download
mkdir -p bin
```

### 3. 构建

```bash
go build -o bin/mcp-webcrawler ./cmd/server
```

### 4. 验证构建

```bash
./bin/mcp-webcrawler -h
```

预有输出：
```
Usage of ./bin/mcp-webcrawler:
  -timeout int
    	请求超时（秒）（默认 15）
```

### 5. 运行测试

```bash
go test -v -short ./...
```

预有输出（应诠出 PASS）：
```
=== RUN   TestCrawlExampleCom
--- PASS: TestCrawlExampleCom (0.52s)
=== RUN   TestCrawlResponseFormat
--- PASS: TestCrawlResponseFormat (0.18s)
```

## 运行服务器

### 启动服务器

```bash
./bin/mcp-webcrawler --timeout 20
```

服务器已准备好通过 stdio 接受 MCP 客户端连接。

### 与 Claude 或 GPT 一起使用

添加到你的 MCP 配置：

```json
{
  "name": "mcp-webcrawler",
  "command": "/path/to/mcp-webcrawler",
  "args": ["--timeout", "20"]
}
```

### 以程序方式使用

创建一个简单的测试脚本：

```bash
#!/bin/bash

# 在后台启动服务器
./bin/mcp-webcrawler &
SERVER_PID=$!
sleep 1

# 发送请求
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "crawl_page",
    "arguments": {
      "url": "https://example.com"
    }
  }
}' | nc localhost 9000

# 清理
kill $SERVER_PID
```

## 测试功能

### 单元测试

```bash
go test -v ./internal/crawl -run TestNewCrawler
```

测试爬虫创始化。

### 集成测试

```bash
go test -v -short ./...
```

测试实际的网页爬取（需要网络）。

### 响应格式测试

```bash
go test -v -short ./... -run TestCrawlResponseFormat
```

验证 JSON 响应结构。

## 后续步骤

1. 检查 BUILD_AND_TEST_GUIDE.md 了解测试详情
2. 查看 API_REFERENCE.md 了解工具文档
3. 调查 IMPLEMENTATION_GUIDE.md 了解架构
4. 查看 examples_test.go 了解使用示例

## 常见问题

### 构建失败输出 "未找到捷径"

```bash
go mod tidy
go get -u ./...
go build -o bin/mcp-webcrawler ./cmd/server
```

### 测试超时

```bash
# 使用简短测试（不测网络）
go test -short ./internal/crawl

# 或增加超时时间
go test -timeout 30s ./...
```

### 服务器不响应

1. 验证它是否正在运行：
   ```bash
   ps aux | grep mcp-webcrawler
   ```

2. 检查日志：
   ```bash
   ./bin/mcp-webcrawler 2>&1 | tee server.log
   ```

3. 使用简单请求测试：
   ```bash
   echo '{}' | ./bin/mcp-webcrawler
   ```

## 项目结构

```
mcp-webcrawler/
├── bin/
│   └── mcp-webcrawler          # 编译的二进制文件
├── cmd/server/
│   └── main.go               # 入口点
├── internal/
│   ├── crawl/
│   │   ├── types.go
│   │   ├── crawler.go
│   │   └── crawler_test.go
│   └── mcp/
│       └── server.go
├── go.mod / go.sum
├── QUICKSTART.md             # 此文件
├── SETUP.md                 # 详细设置
├── BUILD_AND_TEST_GUIDE.md  # 构建与测试详情
├── API_REFERENCE.md         # API 文档
├── IMPLEMENTATION_GUIDE.md  # 实现细节
└── README.md                # 项目概览
```

## 包含内容

- ✅ 需要 goquery 解析器的网页爬取
- ✅ MCP 服务器框架集成
- ✅ 带有类型化模式的工具注册
- ✅ 全面的测试套件
- ✅ 完整的文档
- ✅ 错误处理和超时
- ✅ 内容提取（文本、链接、图片）
- ✅ JSON 响应格式化

## 后续步骤

1. 添加 JavaScript 渲染支持
2. 实现递归网站爬取
3. 添加缓存层（Redis）
4. 实现速率限制
5. 添加身份验证支持
6. 支持自定义提取器
7. 指标和监控

## 获取帮助

1. 检查 BUILD_AND_TEST_GUIDE.md 了解测试详情
2. 查看 API_REFERENCE.md 了解工具文档
3. 调查 IMPLEMENTATION_GUIDE.md 了解架构
4. 查看 examples_test.go 了解使用示例

## 架构

简单的分层架构：

```
MCP 客户端
    ↑↓ stdio/HTTP
MCP 服务器 (mcp/server.go)
    ↑↓ 工具调用
爬虫 (crawl/crawler.go)
    ↑↓ HTTP 请求
网页
```

## 性能

- 二进制大小: 11 MB
- 构建时间: <1 秒
- 典制响应时间: 0.5-2 秒/页
- 每个请求内存: 1-10 MB

## 许可证与附幸

构建于：
- [go-sdk](https://github.com/modelcontextprotocol/go-sdk) - MCP 协议
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML 解析

---

**准备开始了！** 运行 `./bin/mcp-webcrawler` 来启动。
