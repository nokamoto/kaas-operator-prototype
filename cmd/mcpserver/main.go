package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nokamoto/kaas-operator-prototype/internal/apiclient"
	"github.com/nokamoto/kaas-operator-prototype/internal/mcp/clustermanagement"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kaas-operator-prototype",
		Version: "v0.0.0",
	}, nil)

	r := apiclient.New(func() string {
		baseURL := os.Getenv("KAAAS_OPERATOR_PROTOTYPE_API_URL")
		if baseURL == "" {
			baseURL = apiclient.DefaultBaseURL
		}
		return baseURL
	})

	// Register the Cluster Management tool
	createTool := clustermanagement.New(r)
	mcp.AddTool(server, createTool.Tool(), createTool.Handler)

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}
