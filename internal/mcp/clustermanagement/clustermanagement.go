package clustermanagement

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
)

type runtime interface {
	ClusterService() v1alpha1connect.ClusterServiceClient
	LongRunningOperationService() v1alpha1connect.LongRunningOperationServiceClient
}

// Tool is an interface for a tool that can be registered with the MCP server.
type Tool[T any] interface {
	New() *mcp.Tool
	Handler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[T]) (*mcp.CallToolResultFor[any], error)
}
