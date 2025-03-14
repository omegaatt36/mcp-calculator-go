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
		server.WithLogging(),
	)

	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
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

	s.AddTool(calculatorTool, calculateHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

func calculateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, ok := request.Params.Arguments["operation"].(string)
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("invalid operation type, expected string, got %T", request.Params.Arguments["operation"])), nil
	}

	operation, err := ParseOperation(op)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	x, ok := request.Params.Arguments[ArgumentX.String()].(float64)
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("invalid x value type, expected float64, got %T", request.Params.Arguments[ArgumentX.String()])), nil
	}

	y, ok := request.Params.Arguments[ArgumentY.String()].(float64)
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("invalid y value type, expected float64, got %T", request.Params.Arguments[ArgumentY.String()])), nil
	}

	result, err := calculate(operation, decimal.NewFromFloat(x), decimal.NewFromFloat(y))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
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
