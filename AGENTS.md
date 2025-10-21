# 代码库指南

## 项目概述

**mcp-webcrawler** 是用 Go 编写的模型上下文协议（MCP）服务器，为大语言模型提供网页爬取和内容提取功能。它使 LLM 能够获取和解析网页，提取标题、描述、主要内容和链接等结构化信息。

---

## 项目结构

```
mcp-webcrawler/
├── cmd/
│   └── server/
│       └── main.go              # 服务器入口点
├── internal/
│   ├── crawl/
│   │   ├── types.go            # 数据结构定义
│   │   ├── crawler.go          # 网页爬取和 HTML 解析
│   │   └── crawler_test.go     # 单元测试
│   └── mcp/
│       └── server.go           # MCP 协议实现
├── go.mod                       # 模块依赖定义
├── go.sum                       # 依赖校验和
├── examples_test.go             # 集成测试
├── README.md                    # 项目愿景和设计
├── DEVELOPMENT.md               # 架构和实现细节
├── SETUP.md                     # 设置和故障排除
└── bin/                         # 编译二进制文件（生成）
```

---

## 构建和开发命令

**安装依赖：**
```bash
go mod download && go mod tidy
```

**构建二进制文件：**
```bash
mkdir -p bin && go build -o bin/mcp-webcrawler ./cmd/server
```

**运行服务器：**
```bash
./bin/mcp-webcrawler --timeout 15
```
选项：`--timeout` 设置请求超时时间（秒），默认为 15。

**代码格式化和检查：**
```bash
go fmt ./...
go vet ./...
```

**清理构建缓存：**
```bash
go clean -cache && go clean -testcache
```

---

## 测试

**运行所有测试：**
```bash
go test -v ./...
```

**运行带覆盖率的测试：**
```bash
go test -v -cover ./...
```

**运行特定测试：**
```bash
go test -v -run TestCrawlExampleCom ./...
```

**测试分类：**
- **单元测试：** `internal/crawl/crawler_test.go`（NewCrawler、FetchPage）
- **集成测试：** `examples_test.go`（对 example.com 的实时网络测试）

测试需要网络访问。覆盖率目标：核心包 70% 以上。

---

## 代码风格和命名规范

**重要提示：所有代码必须使用中文注释！**

- **语言版本：** Go 1.21+
- **缩进：** 制表符（Go 标准）
- **包：** 保持专注和精简（`crawl` 用于爬取逻辑，`mcp` 用于协议）
- **函数名：** 导出函数使用 PascalCase，非导出使用 camelCase
- **注释语言：** **所有代码注释必须用中文编写**
- **结构体名：** 清晰描述用途（PageContent、CrawlRequest、CrawlResponse）
- **错误处理：** 始终显式处理；提供描述性消息
- **导入：** 先导入标准库，再导入外部包

**必须执行 `go fmt` 和 `go vet` 后再提交。**

---

## 核心组件

**爬虫（internal/crawl/crawler.go）：**
- `NewCrawler(timeout)`：使用请求超时时间初始化
- `FetchPage(ctx, url)`：下载并解析单个网页
- 返回 `PageContent`，包含标题、描述、内容、链接、图片、状态码等

**MCP 服务器（internal/mcp/server.go）：**
- 实现 MCP 协议工具注册
- 注册 `crawl_page` 工具供 LLM 使用
- 将爬虫响应格式化为 LLM 友好的 JSON

**数据类型（internal/crawl/types.go）：**
- `PageContent`：完整页面数据（URL、标题、描述、内容、链接、图片、元数据）
- `CrawlRequest`：请求参数（URL、渲染 JS、最大长度、格式、超时）
- `CrawlResponse`：结构化 LLM 输出（标题、描述、内容、链接、抓取时间）

---

## MCP 工具：crawl_page

**输入参数：**
- `url`（字符串，必需）：要爬取的目标 URL
- `render_js`（布尔值，可选）：是否渲染 JavaScript 内容
- `max_content_length`（整数，可选）：限制响应大小（字节）

**输出格式：**
```json
{
  "title": "页面标题",
  "description": "元描述文本",
  "content": "提取的主要内容...",
  "links": ["http://example.com/path1", "http://example.com/path2"],
  "fetched_at": "2024-01-01T12:00:00Z"
}
```

---

## 添加新功能

1. **先写测试**（TDD）：在 `*_test.go` 文件中添加测试用例
2. **实现功能：** 在相应包中添加代码，**使用中文注释**
3. **运行测试：** `go test -v ./...`
4. **格式化代码：** `go fmt ./...` 和 `go vet ./...`
5. **更新文档：** 如有重大变更，更新 DEVELOPMENT.md

未来可考虑的工具：`crawl_site`（递归爬取）、`extract_content`（自定义选择器）、`summarize_page`。

---

## 依赖项

- `github.com/modelcontextprotocol/go-sdk`：MCP 协议实现
- `github.com/PuerkitoBio/goquery`：HTML 解析和 CSS 选择器
- `github.com/go-rod/rod`：浏览器自动化（用于 JavaScript 渲染）
- `golang.org/x/net`：网络工具库

保持依赖最少且最新。变更后总是运行 `go mod tidy`。

---

## 安全和最佳实践

- **URL 验证：** 拒绝内网 IP（127.0.0.1、10.x.x.x、172.x.x.x、192.168.x.x）、file:// URL 和 localhost
- **超时控制：** 始终设置上下文超时以防止请求挂起
- **错误处理：** 记录错误但不向 LLM 响应中泄露敏感信息
- **速率限制：** 在请求间实现延迟以尊重服务器资源
- **用户代理：** 使用描述性头标识 MCP 爬虫
- **内容限制：** 强制执行 max_content_length 防止内存问题

---

## 调试提示

**详细测试输出：**
```bash
go test -v -run TestName ./...
```

**验证依赖：**
```bash
go mod verify
```

**性能分析（内存或 CPU）：**
```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof -v ./...
go tool pprof cpu.prof
```

**运行竞态条件检测：**
```bash
go test -race ./...
```

---

## 向大模型的要求

### 代码规范要求

**所有代码贡献必须遵守以下规则：**

1. **使用中文注释：** 所有代码注释必须用中文编写，清晰解释逻辑和意图
2. **使用中文回答：** 当被问及代码问题或提供解释时，必须用中文回答
3. **代码示例中文化：** 在代码示例中包含中文注释和文档字符串
4. **文档说明：** 提交的任何文档变更应使用中文编写

**示例：**
```go
// 创建新的爬虫实例，设置超时时间为给定的时长
func NewCrawler(timeout time.Duration) *Crawler {
    return &Crawler{
        timeout: timeout,
        // 初始化 HTTP 客户端
        client: &http.Client{
            Timeout: timeout,
        },
    }
}

// 获取指定 URL 的网页内容并解析
func (c *Crawler) FetchPage(ctx context.Context, url string) (*PageContent, error) {
    // 验证 URL 格式
    // 发送 HTTP 请求
    // 解析 HTML 并提取内容
    // 返回结构化数据
}
```

### 测试要求

**完成任何功能实现前必须：**

1. **编写测试用例** — 为新功能添加相应的单元测试或集成测试
2. **运行所有测试** — 确保所有测试通过（包括现有测试）：
   ```bash
   go test -v ./...
   ```
3. **验证覆盖率** — 核心功能覆盖率应达到 70% 以上
4. **检查代码质量** — 执行格式化和静态检查：
   ```bash
   go fmt ./...
   go vet ./...
   ```

**任何测试失败或代码检查不通过的提交将被拒绝。**

### Git Commit 规范

**每次完成任务必须进行 Git 提交，要求如下：**

1. **必须提交** — 每个完整的功能实现必须通过 `git commit` 进行提交
2. **中文 Commit 信息** — 所有 commit 信息必须用中文编写
3. **Commit 信息格式**
   ```
   <类型>: <简短描述>
   
   <详细说明（可选）>
   ```

4. **Commit 类型（中文）**
   - `功能`: 新增功能或工具
   - `修复`: 修复 bug 或问题
   - `改进`: 代码优化、性能提升
   - `测试`: 添加或更新测试
   - `文档`: 文档更新或添加
   - `重构`: 代码重构或架构调整
   - `依赖`: 依赖库更新或管理

5. **Commit 信息示例**
   ```
   功能: 添加 crawl_page 工具的 JavaScript 渲染支持
   
   - 集成 go-rod 库进行浏览器自动化
   - 添加 render_js 参数以控制是否执行 JavaScript
   - 完整的单元测试和集成测试覆盖
   - 更新文档说明新功能用法
   ```

   ```
   修复: 修复 URL 验证绕过内网 IP 的安全漏洞
   
   - 改进内网 IP 检测逻辑
   - 添加对私有地址范围的全面检查
   - 添加相应的单元测试验证修复
   ```

   ```
   测试: 提高爬虫模块的测试覆盖率到 75%
   
   - 为 FetchPage 函数添加 10 个新测试用例
   - 测试超时处理、错误处理和边界情况
   - 所有测试通过，覆盖率达到要求
   ```

   ```
   文档: 更新 DEVELOPMENT.md 中的组件说明
   
   - 补充爬虫组件的详细实现说明
   - 添加常见问题和解决方案
   - 完善 MCP 工具接口文档
   ```

6. **Commit 前检查清单**
   - [ ] 所有测试都通过（`go test -v ./...`）
   - [ ] 代码已格式化（`go fmt ./...`）
   - [ ] 没有代码质量问题（`go vet ./...`）
   - [ ] 所有代码注释为中文
   - [ ] 新增功能有相应测试
   - [ ] 文档已更新（如需要）
   - [ ] Commit 信息为中文，格式正确

7. **提交命令**
   ```bash
   # 添加所有变更
   git add -A
   
   # 提交，使用中文信息
   git commit -m "功能: 新增功能描述"
   
   # 或使用多行详细信息
   git commit -m "功能: 简短描述" -m "详细说明"
   ```

### 示例工作流程

**完整的任务完成流程：**

```bash
# 1. 实现新功能（添加中文注释）
# 编辑相关文件，添加实现代码

# 2. 编写或更新测试
go test -v ./...

# 3. 格式化和检查
go fmt ./...
go vet ./...

# 4. 验证所有测试通过
go test -v -cover ./...

# 5. 查看变更
git status
git diff

# 6. 提交代码（中文 commit 信息）
git add -A
git commit -m "功能: 添加新的网页爬取功能"

# 7. 推送到远程（如适用）
git push origin main
```

---

## 资源

- [DEVELOPMENT.md](DEVELOPMENT.md) — 详细架构和组件说明
- [SETUP.md](SETUP.md) — 设置指南和故障排除
- [README.md](README.md) — 项目愿景和设计规范
- [MCP 协议规范](https://modelcontextprotocol.io/) — 官方 MCP 文档
- [Go 文档](https://go.dev/doc/) — Go 标准库和工具参考
