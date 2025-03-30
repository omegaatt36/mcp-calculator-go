package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shopspring/decimal"
)

func main() {
	s := server.NewMCPServer(
		"Calculator Go",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString(ArgumentOp.String(),
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum(
				OperationAdd.String(),
				OperationSub.String(),
				OperationMul.String(),
				OperationDiv.String(),
			),
		),
		mcp.WithNumber(ArgumentX.String(),
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber(ArgumentY.String(),
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	calculatorPrompt := mcp.NewPrompt("calc",
		mcp.WithPromptDescription("Perform basic arithmetic operations"),
		mcp.WithArgument(ArgumentOp.String(),
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("The operation to perform (add, subtract, multiply, divide)"),
		),

		mcp.WithArgument(ArgumentX.String(),
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("First number"),
		),
		mcp.WithArgument(ArgumentY.String(),
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Second number"),
		),
	)

	s.AddTool(calculatorTool, calculateToolHandler)
	s.AddPrompt(calculatorPrompt, calculatePromptHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

func calculatePromptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	op, ok := request.Params.Arguments[ArgumentOp.String()]
	if !ok {
		return mcp.NewGetPromptResult(fmt.Sprintf("invalid operation type, expected string, got %T", request.Params.Arguments[ArgumentOp.String()]), nil), nil
	}

	operation, err := ParseOperation(op)
	if err != nil {
		return mcp.NewGetPromptResult(err.Error(), nil), nil
	}

	paramX, ok := request.Params.Arguments[ArgumentX.String()]
	if !ok {
		return mcp.NewGetPromptResult(fmt.Sprintf("invalid x value type, expected float64, got %T", request.Params.Arguments[ArgumentX.String()]), nil), nil
	}

	paramY, ok := request.Params.Arguments[ArgumentY.String()]
	if !ok {
		return mcp.NewGetPromptResult(fmt.Sprintf("invalid y value type, expected float64, got %T", request.Params.Arguments[ArgumentY.String()]), nil), nil
	}

	x, err := decimal.NewFromString(paramX)
	if err != nil {
		return mcp.NewGetPromptResult(err.Error(), nil), nil
	}

	y, err := decimal.NewFromString(paramY)
	if err != nil {
		return mcp.NewGetPromptResult(err.Error(), nil), nil
	}

	result, err := calculate(operation, x, y)
	if err != nil {
		return mcp.NewGetPromptResult(err.Error(), nil), nil
	}

	return mcp.NewGetPromptResult("result", []mcp.PromptMessage{
		{
			Role:    mcp.RoleUser,
			Content: mcp.NewTextContent(fmt.Sprintf("user want to calculate '%s %s %s'", paramX, operation, paramY)),
		},
		{
			Role:    mcp.RoleAssistant,
			Content: mcp.NewTextContent(fmt.Sprintf("result is %s", result.String())),
		},
	}), nil
}

func calculateToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, ok := request.Params.Arguments[ArgumentOp.String()].(string)
	if !ok {
		return mcp.NewToolResultText(fmt.Sprintf("invalid operation type, expected string, got %T", request.Params.Arguments[ArgumentOp.String()])), nil
	}

	operation, err := ParseOperation(op)
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}

	x, ok := request.Params.Arguments[ArgumentX.String()].(float64)
	if !ok {
		return mcp.NewToolResultText(fmt.Sprintf("invalid x value type, expected float64, got %T", request.Params.Arguments[ArgumentX.String()])), nil
	}

	y, ok := request.Params.Arguments[ArgumentY.String()].(float64)
	if !ok {
		return mcp.NewToolResultText(fmt.Sprintf("invalid y value type, expected float64, got %T", request.Params.Arguments[ArgumentY.String()])), nil
	}

	result, err := calculate(operation, decimal.NewFromFloat(x), decimal.NewFromFloat(y))
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}

	return mcp.NewToolResultText(result.String()), nil
}

// calculate performs the arithmetic operation specified by op on x and y.
func calculate(op Operation, x, y decimal.Decimal) (decimal.Decimal, error) {
	switch op {
	case OperationAdd:
		return x.Add(y), nil
	case OperationSub:
		return x.Sub(y), nil
	case OperationMul:
		return x.Mul(y), nil
	case OperationDiv:
		if y.IsZero() {
			return decimal.Zero, fmt.Errorf("cannot divide by zero")
		}
		return x.Div(y), nil
	default:
		return decimal.Zero, fmt.Errorf("invalid operation: %v", op)
	}
}
