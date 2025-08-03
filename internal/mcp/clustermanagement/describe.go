package clustermanagement

import (
	"context"
	"fmt"

	"buf.build/go/protoyaml"
	"connectrpc.com/connect"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
)

type DescribeLongRunningOperationRequest struct {
	Name string `json:"name" jsonschema:"required. The name of the long-running operation to describe."`
}

// DescribeLongRunningOperationTool is a tool for describing long-running operations.
type DescribeLongRunningOperationTool struct {
	r runtime
}

func NewDescribeLongRunningOperationTool(r runtime) Tool[DescribeLongRunningOperationRequest] {
	return &DescribeLongRunningOperationTool{r: r}
}

func (d *DescribeLongRunningOperationTool) New() *mcp.Tool {
	return &mcp.Tool{
		Name: "describe_kaas_long_running_operation",
		Description: `Manage KaaS clusters.
This tool provides operations to describe long-running operations for KaaS clusters.
It allows you to check the status and details of long-running operations that were initiated, such as cluster creation.`,
	}
}

func (d *DescribeLongRunningOperationTool) Handler(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[DescribeLongRunningOperationRequest],
) (*mcp.CallToolResultFor[any], error) {
	service := d.r.LongRunningOperationService()
	res, err := service.GetOperation(ctx, connect.NewRequest(&v1alpha1.GetOperationRequest{
		Name: params.Arguments.Name,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to describe long-running operation: %w", err)
	}
	yaml, err := protoyaml.Marshal(res.Msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(yaml),
			},
		},
	}, nil
}
