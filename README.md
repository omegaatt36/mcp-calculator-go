# MCP Calculator Go

This is a simple calculator implemented in Go that conforms to the [Model Context Protocol (MCP)](https://www.anthropic.com/news/model-context-protocol). It demonstrates how to build an MCP server using the [`mcp-go`](https://github.com/mark3labs/mcp-go) library.

## Features

- Performs basic arithmetic operations: addition, subtraction, multiplication, and division.
- Uses `github.com/shopspring/decimal` for accurate decimal calculations.
- Exposes a single `calculate` tool with `operation`, `x`, and `y` parameters.
- Uses `go-enum` to generate enums and related helper methods.
- Provides a `Taskfile` for convenient build and dependency management.

## Building and Running
First you need to install the dependencies described in the task file.

```bash
go install github.com/abice/go-enum@latest
go generate ./...
```

or simply

```bash
task deps
task enum
```

To build the executable, run:

```bash
go build -o mcp-calculator-go
```

or

```
task build
```

## How to use

1. build mcp binary

2. configure

2-1. VS Code and Client:

```json
{
  "mcpServers": {
    "calculator-go": {
      "command": "/path/to/mcp-calculator-go",
      "args": [],
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

2-2. Zed

```json
"context_servers": {
  "calculator-go": {
    "command": {
      "path": "/path/to/mcp-calculator-go",
      "args": [],
      "env": {}
    },
    "settings": {}
  }
}
```
