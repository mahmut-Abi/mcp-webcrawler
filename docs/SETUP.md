# MCP WebCrawler - 设置指南

## 快速开始

### 前置要求
- Go 1.21 或更高版本
- git

### 1. 克隆或初始化仓库

```bash
cd mcp-webcrawler
```

### 2. 安装依赖

```bash
go mod download
go mod tidy
```

### 3. 构建项目

```bash
mkdir -p bin
go build -o bin/mcp-webcrawler ./cmd/server
```

### 4. 运行测试

```bash
# 运行所有测试
go test -v ./...

# 使用网络访问运行（用于实时测试）
go test -v -count=1 ./...

# 运行特定测试
go test -v -run TestCrawlExampleCom ./...
```

### 5. 启动服务器

```bash
./bin/mcp-webcrawler --timeout 15
```

## 故障排除

### 构建问题

**问题**: `go: downloading ...` 挂起
- **解决方案**: 确保网络连接正常，尝试 `go mod cache clean`

**问题**: 缺失依赖
- **解决方案**: 运行 `go mod tidy` 清理依赖

### 测试失败

**问题**: 测试中的网络超时
- **解决方案**: 增加超时时间或检查网络连接

**问题**: 无法连接外部 URL
- **解决方案**: 检查防火墙设置和 DNS 解析

## 项目结构

```
mcp-webcrawler/
├── cmd/
│   └─ server/
│       └── main.go              # 入口点
├── internal/
│   ├── crawl/
│   │   ├── types.go            # 数据类型
│   │   ├── crawler.go          # 核心爬虫
│   │   └── crawler_test.go     # 爬虫测试
│   └── mcp/
│       └── server.go           # MCP 服务器
├── go.mod                 # 模块定义
├── go.sum                 # 依赖校验和
├── examples_test.go       # 集成测试
├── DEVELOPMENT.md         # 开发指南
├── SETUP.md               # 此文件
├── README.md              # 项目概览
└── bin/
    └── mcp-webcrawler         # 编译的二进制文件
```

## 测试

### 单元测试

位于 `internal/crawl/crawler_test.go` 的测试：
- `TestNewCrawler`: 验证爬虫实例化
- `TestFetchPage`: 测试从 example.com 的实际页面获取

### 集成测试

位于 `examples_test.go` 的测试：
- `TestCrawlExampleCom`: 完整的页面爬虫示例
- `TestCrawlResponseFormat`: 响应 JSON 格式化

### 运行测试

```bash
# 运行所有测试并显示详细输出
go test -v ./...

# 运行测试并显示覆盖率报告
go test -v -cover ./...

# 运行特定测试
go test -v -run TestFetchPage ./internal/crawl

# 运行简短模式的测试（输出较少）
go test -short ./...
```

## 开发工作流程

### 添加新功能

1. 在适当的包中创建新文件
2. 先写测试（TDD 方法）
3. 实现功能
4. 运行测试: `go test -v ./...`
5. 格式化代码: `go fmt ./...`
6. 检查问题: `go vet ./...`

### 调试

```bash
# 启用详细日志
go test -v -run TestName

# 使用调试器运行
dlv test ./cmd/server
```

## 环境变量

当前不需要任何环境变量。服务器接受命令行标志：

```bash
./bin/mcp-webcrawler --timeout 20  # 自定义超时
```

## 常用命令参考

```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 构建二进制文件
go build -o bin/mcp-webcrawler ./cmd/server

# 运行服务器
./bin/mcp-webcrawler

# 运行测试
go test -v ./...

# 格式化代码
go fmt ./...

# 分析代码
go vet ./...

# 检查安全问题
go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
```

## 后续步骤

1. 启动服务器: `./bin/mcp-webcrawler`
2. 连接 MCP 客户端
3. 使用 `crawl_page` 工具获取网页内容
4. 根据需要扩展其他工具
