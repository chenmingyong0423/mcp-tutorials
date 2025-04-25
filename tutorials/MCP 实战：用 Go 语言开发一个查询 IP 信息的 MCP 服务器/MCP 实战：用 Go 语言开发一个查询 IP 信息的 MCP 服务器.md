# 📝 前言

随着 `MCP` 的快速普及和广泛应用，`MCP` 服务器也层出不穷。大多数开发者使用的 `MCP` 服务器开发库是官方提供的 `typescript-sdk`，而作为 `Go` 开发者，我们也可以借助优秀的第三方库去开发 `MCP` 服务器，例如 `ThinkInAIXYZ/go-mcp`。

本文将详细介绍如何在 `Go` 语言中使用 `go-mcp` 库来开发一个查询 `IP` 信息的 `MCP` 服务器。

准备好了吗？准备一杯你最喜欢的咖啡或茶，随着本文一探究竟吧。

![image](https://github.com/user-attachments/assets/fef15469-89d4-4d0a-ae08-143ec0a31396)


# 🌐 mcp-ip-geo 服务器

`mcp-ip-geo` 是一个用于查询 `IP` 信息的 `MCP` 服务器，项目已开源，仓库地址：[mcp-ip-geo](https://github.com/chenmingyong0423/mcp-ip-geo)。

## 📁 目录结构说明

    ├─cmd
    │  └─mcp-ip-geo
    └─internal
        ├─domain
        ├─server
        ├─service
        └─tools

*   `cmd/mcp-ip-geo`：应用的启动入口目录，包含如 `main.go` 启动文件。
*   `internal/domain`：定义项目中的核心数据结构，例如 `IP API` 的响应体等。
*   `internal/server`：`MCP` 服务器的核心逻辑实现。
*   `internal/service`：对接第三方服务的实现，如调用 `IP` 解析 `API`。
*   `internal/tools`：`MCP` 工具的具体实现，支持灵活扩展和注册。

# 🔍 查询 IP 信息功能实现

代码位于 `service` 包中，通过 `ip-api.com` 提供的接口获取 `IP` 地理位置信息，具体实现如下：

```go
package service

import (
    "context"
    "fmt"
    "github.com/chenmingyong0423/mcp-ip-geo/internal/domain"
    "net/http"
    "time"

    httpchain "github.com/chenmingyong0423/go-http-chain"
)

func NewIpApiService() *IpApiService {
    return &IpApiService{
       host: "http://ip-api.com",
       client: httpchain.NewWithClient(&http.Client{
          Timeout: time.Second * 10,
       }),
    }
}

type IIpApiService interface {
    GetLocation(ctx context.Context, ip string) (*domain.IpApiResponse, error)
}

var _ IIpApiService = (*IpApiService)(nil)

type IpApiService struct {
    host   string
    client *httpchain.Client
}

func (s *IpApiService) GetLocation(ctx context.Context, ip string) (*domain.IpApiResponse, error) {
    var resp domain.IpApiResponse
    err := s.client.Get(fmt.Sprintf("%s/json/%s", s.host, ip)).DoAndParse(ctx, &resp)
    if err != nil {
       return nil, err
    }
    return &resp, nil
}

```

**代码解释：**

*   服务初始化（`NewIpApiService`）
    *   创建一个新的 `IpApiService` 实例。
    *   设置了 `API` 地址为 `http://ip-api.com`。
    *   使用 `httpchain` 封装的 `HTTP` 客户端，设置请求超时时间为 10 秒。

*   接口定义（`IIpApiService`）
    *   定义了服务对外暴露的功能：`GetLocation` 方法，用于获取 `IP` 地理位置信息。  
    *   使用接口有助于后续做依赖注入、`mock` 测试等。  
    *   `var _ IIpApiService = (*IpApiService)(nil)` 这行代码用于编译时检查，确保 `IpApiService` 实现了 `IIpApiService` 接口。

*   结构体定义（`IpApiService`） 
    *   包含两个字段：
    *   `host`：`API` 的基础地址。
    *   `client`：封装的 `HTTP` 客户端，类型为 `*httpchain.Client`。

*   核心方法实现（`GetLocation`）
    *   根据传入的 `IP` 构造请求地址：`http://ip-api.com/json/{ip}`。  
    *   使用 `httpchain` 库发起 `GET` 请求，并将结果解析到 `domain.IpApiResponse` 结构体中。  

# 🧰 工具实现

## 🔧 工具管理

代码位于 `tools` 包中，用于管理工具，具体实现如下：

```go
package tools

import (
    "github.com/ThinkInAIXYZ/go-mcp/protocol"
    "github.com/ThinkInAIXYZ/go-mcp/server"
)

type ToolFunc func() (tool *protocol.Tool, toolHandler server.ToolHandlerFunc)

func GetToolFuncList() []ToolFunc {
    return []ToolFunc{
       SingleIpParser,
    }
}

```

代码解释：

*   **`ToolFunc` 类型定义**

    *   定义了一个函数类型 `ToolFunc`，返回两个值：

        *   `*protocol.Tool`：工具的元信息；
        *   `server.ToolHandlerFunc`：该工具的处理逻辑函数。

    *   用这种方式可以将 **工具的定义** 与 **工具的执行逻辑** 一并管理，后续在定义工具时都可以通过实现该函数签名进行表示。

*   **`GetToolFuncList` 函数**

    *   返回一个 `ToolFunc` 列表。
    *   当前只注册了一个工具：`SingleIpParser`，但这种结构易于扩展，后续只需往列表中添加新的工具函数即可。
    *   通过集中注册，应用在初始化时可以统一加载所有工具。

## 🌍 查询单个 IP 信息工具的实现

代码位于 `tools` 包中，用于查询单个 `IP` 信息，具体实现如下：

```go
package tools

import (
    "context"
    "encoding/json"
    "github.com/ThinkInAIXYZ/go-mcp/protocol"
    "github.com/ThinkInAIXYZ/go-mcp/server"
    "github.com/chenmingyong0423/mcp-ip-geo/internal/service"
)

var singleIpParserTool *protocol.Tool

type ipRequest struct {
    Ip string `json:"ip"`
}

func init() {
    var err error
    singleIpParserTool, err = protocol.NewTool("ip-details", "a tool that provides IP geolocation information", ipRequest{})
    if err != nil {
       panic(err)
    }
}

func SingleIpParser() (*protocol.Tool, server.ToolHandlerFunc) {
    ipApiService := service.NewIpApiService()

    return singleIpParserTool, func(toolRequest *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
       var req ipRequest
       if err := protocol.VerifyAndUnmarshal(toolRequest.RawArguments, &req); err != nil {
          return nil, err
       }
       resp, err := ipApiService.GetLocation(context.Background(), req.Ip)
       if err != nil {
          return nil, err
       }

       marshal, err := json.Marshal(resp)
       if err != nil {
          return nil, err
       }

       return &protocol.CallToolResult{
          Content: []protocol.Content{
             protocol.TextContent{
                Type: "text",
                Text: string(marshal),
             },
          },
       }, nil
    }
}

```

**代码解释：**

*   **全局变量声明**

    *   `singleIpParserTool`：存储工具元信息的协议工具对象
    *   `ipRequest`：定义工具输入参数结构体，包含 `ip` 字符串字段

*   **初始化函数（`init`）**

    *   在包加载时通过 `protocol.NewTool` 创建工具元信息
    *   指定工具标识符 `ip-details`，描述信息和输入参数结构体 `ipRequest{}`
    *   错误处理采用 `panic`，确保工具元信息必须正确初始化

*   **工具注册函数（`SingleIpParser`）**

    *   创建 `IpApiService` 服务实例用于 `IP` 定位查询

    *   返回两个值：

        *   预定义的 `singleIpParserTool` 元信息对象
        *   工具处理函数

*   **工具处理函数**

    *   **参数验证与解析**

        *   调用 `protocol.VerifyAndUnmarshal` 验证请求参数有效性
        *   将原始参数反序列化到 `ipRequest` 结构体

    *   **服务调用**

        *   使用 `ipApiService.GetLocation` 获取 `IP` 地理位置信息

    *   **结果处理**

        *   将服务响应结果序列化为 `JSON` 字符串并包装为 `protocol.CallToolResult` 结构体返回

# 🚀 服务器的创建与启动

代码位于 `server` 包中，用于初始化服务并启动服务端，具体实现如下：

```go
package server

import (
    "github.com/ThinkInAIXYZ/go-mcp/server"
    "github.com/ThinkInAIXYZ/go-mcp/transport"
    "github.com/chenmingyong0423/mcp-ip-geo/internal/tools"
)

func Run(address string) error {
    var err error

    var ts transport.ServerTransport
    if address == "" {
       ts = transport.NewStdioServerTransport()
    } else {
       ts, err = transport.NewSSEServerTransport(address)
       if err != nil {
          return err
       }
    }

    s, err := server.NewServer(ts)
    if err != nil {
       return err
    }

    toolFuncList := tools.GetToolFuncList()
    for _, tool := range toolFuncList {
       s.RegisterTool(tool())
    }

    return s.Run()
}

```

**代码解释：**

*   **传输层初始化**

    *   根据 `address` 参数判断运行模式：

        *   **空地址模式**：使用 `NewStdioServerTransport` 创建标准输入输出传输，适用于命令行工具等场景。
        *   **指定地址模式**：使用 `NewSSEServerTransport` 创建 `SSE` (`Server-Sent Events`) 传输，适用于 `HTTP` 长连接服务。

*   **服务实例化**

    *   使用 `server.NewServer` 方法创建服务实例，注入配置好的传输层对象 `ts`。

*   **工具注册**

    *   调用 `tools.GetToolFuncList` 获取所有预定义的工具函数列表。

    *   遍历工具列表，通过 `s.RegisterTool(tool())` 注册每个工具：

        *   `tool()` 执行后返回元信息 `*protocol.Tool` 和处理函数 `ToolHandlerFunc`。

*   **服务启动**

    *   调用 `s.Run()` 启动服务，开始监听请求。

# 🧩 主程序入口实现

代码位于 `main` 包中，作为程序启动入口，具体实现如下：

```go
package main

import (
    "flag"
    "github.com/chenmingyong0423/mcp-ip-geo/internal/server"
)

func main() {
    addr := flag.String("address", "", "The host and port to run the sse server")
    flag.Parse()

    if err := server.Run(*addr); err != nil {
       panic(err)
    }
}

```

**代码解释：**

*   **命令行参数解析**

    *   定义 `address` 参数：

        *   参数名称：`-address`
        *   默认值：空字符串
        *   描述：指定 `SSE` 服务运行的地址和端口

    *   调用 `flag.Parse()` 解析命令行参数

*   **服务启动**

    *   调用 `server.Run(*addr)` 启动服务
    *   将解析后的 `address` 参数值传递给服务启动函数

# ⚙️ 从源码构建


## 🛠 本地构建

### 💻 使用 Go 命令

[](https://github.com/chenmingyong0423/mcp-ip-geo/blob/main/README-zh_CN.md#%E4%BD%BF%E7%94%A8-go-%E5%91%BD%E4%BB%A4)

```
# 在类 Unix 系统（Linux/macOS）上
go build -o mcp-ip-geo ./cmd/mcp-ip-geo

# 在 Windows 上
go build -o mcp-ip-geo.exe .\cmd\mcp-ip-geo
```

### 🐋 使用 Docker

[](https://github.com/chenmingyong0423/mcp-ip-geo/blob/main/README-zh_CN.md#%E4%BD%BF%E7%94%A8-docker)

1.  构建 Docker 镜像：

    ```
    docker build -t mcp-ip-geo-server .
    ```

1.  运行容器：

    ```
    docker run -d --name mcp-ip-geo-server -p 8000:8000 mcp-ip-geo-server
    ```

## 📦 安装预编译版本

[](https://github.com/chenmingyong0423/mcp-ip-geo/blob/main/README-zh_CN.md#%E5%AE%89%E8%A3%85%E9%A2%84%E7%BC%96%E8%AF%91%E7%89%88%E6%9C%AC)

使用 `Go` 安装最新版本的服务：

```
go install github.com/chenmingyong0423/mcp-ip-geo/cmd/mcp-ip-geo@latest
```

# 🧩 MCP 集成

[](https://github.com/chenmingyong0423/mcp-ip-geo/blob/main/README-zh_CN.md#mcp-%E9%9B%86%E6%88%90)

你可以通过以下两种方式集成 `mcp-ip-geo` 服务：

## 🖥 可执行文件集成（本地运行）

[](https://github.com/chenmingyong0423/mcp-ip-geo/blob/main/README-zh_CN.md#-%E5%8F%AF%E6%89%A7%E8%A1%8C%E6%96%87%E4%BB%B6%E9%9B%86%E6%88%90%E6%9C%AC%E5%9C%B0%E8%BF%90%E8%A1%8C)

```
{
  "mcpServers": {
    "mcp-ip-geo": {
      "command": "/path/to/mcp-ip-geo"
    }
  }
}
```

## 🌐 HTTP 接口集成（连接到已运行的服务）

[](https://github.com/chenmingyong0423/mcp-ip-geo/blob/main/README-zh_CN.md#-http-%E6%8E%A5%E5%8F%A3%E9%9B%86%E6%88%90%E8%BF%9E%E6%8E%A5%E5%88%B0%E5%B7%B2%E8%BF%90%E8%A1%8C%E7%9A%84%E6%9C%8D%E5%8A%A1)

```
{
  "mcpServers": {
    "mcp-ip-geo": {
      "url": "http://host:port/sse"
    }
  }
}
```

# 👀 效果演示

![企业微信截图_17455550209491](https://github.com/user-attachments/assets/d2f7f0ac-2a97-43e2-9dd0-341bd209c234)


# ✅ 小结

本文将详细介绍 `mcp-ip-geo` —— 一个用于查询 `IP` 信息的 `MCP` 服务器的实现细节。该服务器目前支持两种数据传输方式：`stdio` 和 `SSE（Server-Sent Events）`。未来还计划支持 `Streamable HTTP` 传输方式，并持续扩展更多实用的工具（`tools`）模块。
