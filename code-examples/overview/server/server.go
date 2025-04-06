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
