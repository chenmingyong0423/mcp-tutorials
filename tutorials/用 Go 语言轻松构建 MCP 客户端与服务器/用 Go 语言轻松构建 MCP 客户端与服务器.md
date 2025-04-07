# 前言
模型上下文协议（`Model Context Protocol`，简称 `MCP`）是一种开放标准，旨在标准化大型语言模型（`LLM`）与外部数据源和工具之间的交互方式。随着 `MCP` 越来越受欢迎，`Go MCP` 库应运而生。本文将介绍如何在 `Go` 语言里面轻松构建 `MCP` 客户端和服务器。

准备好了吗？准备一杯你最喜欢的咖啡或茶，随着本文一探究竟吧。

# mcp-go
要构建 `MCP` 客户端和服务器，我们需要使用 `mcp-go` 库。

`mcp-go` 是 `Go` 语言实现的 `Model Context Protocol`（`MCP`）库，通过这个库可以实现 `LLM` 应用与外部数据源和工具之间的无缝集成。

## 主要特点

- 快速：高级接口意味着更少的代码和更快的开发速度

- 简单：使用极少的样板代码构建 `MCP` 服务器

- 完整：`MCP Go` 旨在提供 `MCP` 核心规范的完整实现


# 安装 MCP 库

在 `Go` 项目根目录下，执行以下命令：

```bash
go get github.com/mark3labs/mcp-go
```

# 构建 MCP 服务器
接下来，我们使用 `mcp-go` 提供的 `server` 模块，构建一个通过 `stidio` 方式连接的 `MCP` 服务器。

## 创建 server 对象

```go
s := server.NewMCPServer(
    "Server Demo",
    "1.0.0",
)
```

创建 `server` 对象时，我们可以指定 **服务器名**，**版本号** 等参数。

## 添加工具（tools）

以下是一个示例，用于创建并注册一个简单的计算器工具：

```go
calculatorTool := mcp.NewTool("calculate",
    mcp.WithDescription("执行基本的算术运算"),
    mcp.WithString("operation",
        mcp.Required(),
        mcp.Description("要执行的算术运算类型"),
        mcp.Enum("add", "subtract", "multiply", "divide"), // 保持英文
    ),
    mcp.WithNumber("x",
        mcp.Required(),
        mcp.Description("第一个数字"),
    ),
    mcp.WithNumber("y",
        mcp.Required(),
        mcp.Description("第二个数字"),
    ),
)

s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    op := request.Params.Arguments["operation"].(string)
    x := request.Params.Arguments["x"].(float64)
    y := request.Params.Arguments["y"].(float64)

    var result float64
    switch op {
    case "add":
        result = x + y
    case "subtract":
        result = x - y
    case "multiply":
        result = x * y
    case "divide":
        if y == 0 {
            return nil, errors.New("不允许除以零")
        }
        result = x / y
    }

    return mcp.FormatNumberResult(result), nil
})
```

**添加工具的步骤如下：**

- 创建工具对象
使用 `mcp.NewTool` 创建一个工具实例。  
  - 第一个参数是工具名称（必须），例如 `"calculate"`。
  - 其余参数通过函数选项（`functional options`）方式传入，例如：
    - `mcp.WithDescription(...)` 添加工具描述；
    - `mcp.WithString(...)` 或 `mcp.WithNumber(...)` 定义参数及其规则（如是否必填、参数说明、枚举限制等）。

- 注册工具到服务器
通过 `s.AddTool` 方法将工具注册到 `MCP` 服务中。
  - 第一个参数是上一步创建的工具对象；
  - 第二个参数是该工具的处理函数（handler），用于实现工具的具体逻辑，如参数解析、运算执行、返回结果等。


## 添加资源（Resources）

下面的示例展示了如何创建并注册一个静态资源，用于读取并提供 `README.md` 文件的内容。

```go
resource := mcp.NewResource(
    "docs://readme",
    "项目说明文档",
    mcp.WithResourceDescription("项目的 README 文件"),
    mcp.WithMIMEType("text/markdown"),
)

s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
    content, err := os.ReadFile("README.md")
    if err != nil {
        return nil, err
    }

    return []mcp.ResourceContents{
        mcp.TextResourceContents{
            URI:      "docs://readme",
            MIMEType: "text/markdown",
            Text:     string(content),
        },
    }, nil
})
```

**添加资源的步骤如下：**

- 创建资源对象  
使用 `mcp.NewResource` 函数创建资源实例。
  - 第一个参数为资源 URI，用于标识资源；
  - 第二个参数为资源名称；
  - 通过函数选项补充更多信息，例如：
    - `mcp.WithResourceDescription(...)` 设置资源描述；
    - `mcp.WithMIMEType(...)` 指定资源的 MIME 类型。

- 注册资源处理函数  
使用 `s.AddResource` 将资源对象注册到服务器，并提供一个处理函数：
  - 该处理函数会在资源被访问时执行；
  - 返回值是资源内容的数组（例如读取本地文件内容并封装为 `TextResourceContents`）。


## 添加提示词（Prompts）

以下示例展示了如何创建并添加一个带参数的简单提示词，用于生成个性化的问候语。

```go
s.AddPrompt(mcp.NewPrompt("greeting",
    mcp.WithPromptDescription("一个友好的问候提示"),
    mcp.WithArgument("name",
        mcp.ArgumentDescription("要问候的人的名字"),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    name := request.Params.Arguments["name"]
    if name == "" {
        name = "朋友"
    }

    return mcp.NewGetPromptResult(
        "友好的问候",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleAssistant,
                mcp.NewTextContent(fmt.Sprintf("你好，%s！今天有什么可以帮您的吗？", name)),
            ),
        },
    ), nil
})
```

---

**添加提示词的步骤如下：**

- 创建提示词对象  
通过 `mcp.NewPrompt` 创建一个提示词定义。
  - 第一个参数是提示词名称；
  - 可通过 `mcp.WithPromptDescription(...)` 添加描述；
  - 使用 `mcp.WithArgument(...)` 定义参数及其说明（如提示词中需要动态插值的内容）。

- 注册提示词处理函数  
使用 `s.AddPrompt` 将提示词对象注册到服务器，并提供对应的处理逻辑函数：
  - 函数接收用户输入参数；
  - 返回一个结构化的提示词响应（如构造一个带有用户名字的问候消息）。

## 启动基于 `stdio` 传输类型的服务器

```go
// 启动基于 stdio 的服务器
if err := server.ServeStdio(s); err != nil {
    fmt.Printf("Server error: %v\n", err)
}
```

使用 `server.ServeStdio` 方法可以启动一个基于标准输入/输出（`stdio`）的 `MCP` 服务器。

这种方式适用于本地集成与命令行工具。

## 启动基于 `sse`（Server-Sent Events）传输类型的服务器

如果需要通过 `HTTP` 的方式提供服务，支持服务端推送数据，可以使用 `SS`E（`Server-Sent Events`）传输模式。


```go
s := server.NewMCPServer(
    "My Server", // Server 名称
    "1.0.0",     // 版本号
)

// 创建基于 SSE 的服务器实例
sseServer := server.NewSSEServer(s)

// 启动服务器，监听指定端口（如 :8080）
err := sseServer.Start(":8080")
if err != nil {
    panic(err)
}

```

与 `stdio` 不同，`sse` 模式基于 `HTTP` 协议，更适合 `Web` 应用中的长连接场景，支持服务端推送数据。  

## 完整的 stdio 代码示例
```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Server Demo",
		"1.0.0",
	)

	// 添加工具
	{
		calculatorTool := mcp.NewTool("calculate",
			mcp.WithDescription("执行基本的算术运算"),
			mcp.WithString("operation",
				mcp.Required(),
				mcp.Description("要执行的算术运算类型"),
				mcp.Enum("add", "subtract", "multiply", "divide"), // 保持英文
			),
			mcp.WithNumber("x",
				mcp.Required(),
				mcp.Description("第一个数字"),
			),
			mcp.WithNumber("y",
				mcp.Required(),
				mcp.Description("第二个数字"),
			),
		)

		s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			op := request.Params.Arguments["operation"].(string)
			x := request.Params.Arguments["x"].(float64)
			y := request.Params.Arguments["y"].(float64)

			var result float64
			switch op {
			case "add":
				result = x + y
			case "subtract":
				result = x - y
			case "multiply":
				result = x * y
			case "divide":
				if y == 0 {
					return nil, errors.New("不允许除以零")
				}
				result = x / y
			}

			return mcp.FormatNumberResult(result), nil
		})
	}

	// 添加资源
	{
		// 静态资源示例 - 暴露一个 README 文件
		resource := mcp.NewResource(
			"docs://readme",
			"项目说明文档",
			mcp.WithResourceDescription("项目的 README 文件"),
			mcp.WithMIMEType("text/markdown"),
		)

		// 添加资源及其处理函数
		s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			content, err := os.ReadFile("README.md")
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "docs://readme",
					MIMEType: "text/markdown",
					Text:     string(content),
				},
			}, nil
		})
	}

	// 添加提示词
	{
		// 简单问候提示
		s.AddPrompt(mcp.NewPrompt("greeting",
			mcp.WithPromptDescription("一个友好的问候提示"),
			mcp.WithArgument("name",
				mcp.ArgumentDescription("要问候的人的名字"),
			),
		), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			name := request.Params.Arguments["name"]
			if name == "" {
				name = "朋友"
			}

			return mcp.NewGetPromptResult(
				"友好的问候",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(
						mcp.RoleAssistant,
						mcp.NewTextContent(fmt.Sprintf("你好，%s！今天有什么可以帮您的吗？", name)),
					),
				},
			), nil
		})
	}

	// 启动基于 stdio 的服务器
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}

```

# 构建 MCP 客户端

接下来，我们使用 `mcp-go` 提供的 `client` 模块，构建一个通过 `stdio` 方式连接到前面打包好的 `MCP` 服务器的客户端。

该客户端将展示以下功能：

- 初始化客户端并连接服务器  
- 获取提示词、资源、工具列表  
- 调用远程工具（tool）

## 创建 MCP 客户端

```go
mcpClient, err := client.NewStdioMCPClient(
    "./client/server", // 服务器可执行文件路径
    []string{},        // 启动参数（如果有）
)
if err != nil {
    panic(err)
}
defer mcpClient.Close()
```

通过 `client.NewStdioMCPClient` 方法可以创建一个基于 `stdio` 传输的客户端，并连接到指定的 `MCP` 服务器可执行文件。

## 初始化客户端连接

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

initRequest := mcp.InitializeRequest{}
initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
initRequest.Params.ClientInfo = mcp.Implementation{
    Name:    "Client Demo",
    Version: "1.0.0",
}

initResult, err := mcpClient.Initialize(ctx, initRequest)
if err != nil {
    panic(err)
}
fmt.Printf("初始化成功，服务器信息: %s %s\n", initResult.ServerInfo.Name, initResult.ServerInfo.Version)
```

初始化操作通过 `Initialize` 方法完成，需指定协议版本及客户端信息。

---

## 获取提示词（Prompts）列表

```go
promptsRequest := mcp.ListPromptsRequest{}
prompts, err := mcpClient.ListPrompts(ctx, promptsRequest)
if err != nil {
    panic(err)
}
for _, prompt := range prompts.Prompts {
    fmt.Printf("- %s: %s\n", prompt.Name, prompt.Description)
    fmt.Println("参数:", prompt.Arguments)
}
```

客户端可以使用 `ListPrompts` 获取服务器上定义的所有提示词，包括名称、描述和参数结构。

## 获取资源（Resources）列表

```go
resourcesRequest := mcp.ListResourcesRequest{}
resources, err := mcpClient.ListResources(ctx, resourcesRequest)
if err != nil {
    panic(err)
}
for _, resource := range resources.Resources {
    fmt.Printf("- uri: %s, name: %s, description: %s, MIME类型: %s\n",
        resource.URI, resource.Name, resource.Description, resource.MIMEType)
}
```

通过 `ListResources` 方法，客户端可以查看服务器上可用的静态或动态资源信息。

## 获取工具（Tools）列表

```go
toolsRequest := mcp.ListToolsRequest{}
tools, err := mcpClient.ListTools(ctx, toolsRequest)
if err != nil {
    panic(err)
}
for _, tool := range tools.Tools {
    fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
    fmt.Println("参数:", tool.InputSchema.Properties)
}
```

通过 `ListTools`，客户端可以获取所有注册的工具信息，方便用户交互式选择或自动生成表单调用。

## 调用工具（Tool）

```go
toolRequest := mcp.CallToolRequest{
    Request: mcp.Request{
        Method: "tools/call",
    },
}
toolRequest.Params.Name = "calculate"
toolRequest.Params.Arguments = map[string]any{
    "operation": "add",
    "x":         1,
    "y":         1,
}

result, err := mcpClient.CallTool(ctx, toolRequest)
if err != nil {
    panic(err)
}
fmt.Println("调用工具结果:", result.Content[0].(mcp.TextContent).Text)
```

通过构造 `CallToolRequest`，客户端可以向 `MCP` 服务器发起工具调用请求，并获取返回的结构化结果。  

在此示例中，我们调用了服务器端注册的 `calculate` 工具，实现 `1 + 1` 运算。

## 完整代码示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {

	// 创建一个基于 stdio 的MCP客户端
	mcpClient, err := client.NewStdioMCPClient(
		"./client/server",
		[]string{},
	)
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("初始化 mcp 客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Client Demo",
		Version: "1.0.0",
	}

	// 初始化MCP客户端并连接到服务器
	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"\n初始化成功，服务器信息: %s %s\n\n",
		initResult.ServerInfo.Name,
		initResult.ServerInfo.Version,
	)

	// 从服务器获取提示词列表
	fmt.Println("提示词列表:")
	promptsRequest := mcp.ListPromptsRequest{}
	prompts, err := mcpClient.ListPrompts(ctx, promptsRequest)
	if err != nil {
		panic(err)
	}
	for _, prompt := range prompts.Prompts {
		fmt.Printf("- %s: %s\n", prompt.Name, prompt.Description)
		fmt.Println("参数:", prompt.Arguments)
	}

	// 从服务器获取资源列表
	fmt.Println()
	fmt.Println("资源列表:")
	resourcesRequest := mcp.ListResourcesRequest{}
	resources, err := mcpClient.ListResources(ctx, resourcesRequest)
	if err != nil {
		panic(err)
	}
	for _, resource := range resources.Resources {
		fmt.Printf("- uri: %s, name: %s, description: %s, MIME类型: %s\n", resource.URI, resource.Name, resource.Description, resource.MIMEType)
	}

	// 从服务器获取工具列表
	fmt.Println()
	fmt.Println("可用工具列表:")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := mcpClient.ListTools(ctx, toolsRequest)
	if err != nil {
		panic(err)
	}

	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
		fmt.Println("参数:", tool.InputSchema.Properties)
	}
	fmt.Println()

	// 调用工具
	fmt.Println("调用工具: calculate")
	toolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	toolRequest.Params.Name = "calculate"
	toolRequest.Params.Arguments = map[string]any{
		"operation": "add",
		"x":         1,
		"y":         1,
	}
	// Call the tool
	result, err := mcpClient.CallTool(ctx, toolRequest)
	if err != nil {
		panic(err)
	}
	fmt.Println("调用工具结果:", result.Content[0].(mcp.TextContent).Text)
}

```

# 小结

本文介绍了如何使用 `mcp-go` 构建一个完整的 `MCP` 应用，包括服务端和客户端两部分。

- 服务端支持注册工具（Tool）、资源（Resource）和提示词（Prompt），并可通过 `stdio` 或 `sse` 模式对外提供服务；
- 客户端通过 `stdio` 连接服务器，支持初始化、列出服务内容、调用远程工具等操作。

