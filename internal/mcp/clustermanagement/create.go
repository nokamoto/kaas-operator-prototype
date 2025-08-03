package clustermanagement

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
)

// CreateClusterRequest is the request for creating a KaaS cluster.
type CreateClusterRequest struct {
	DisplayName string `json:"display_name" jsonschema:"optional. The display name of the cluster."`
	Description string `json:"description" jsonschema:"optional. The description of the cluster."`
}

// CreateClusterTool is a tool for creating KaaS clusters.
type CreateClusterTool struct {
	r runtime
}

func NewCreateClusterTool(r runtime) Tool[CreateClusterRequest] {
	return &CreateClusterTool{r: r}
}

func (c *CreateClusterTool) New() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_kaas_cluster",
		Description: `Manage KaaS clusters.
This tool provides create operations for KaaS clusters.
Note:
Cluster creation typically takes several minutes to complete, so it is recommended to use the long-running operation tool to check the status of the cluster creation.`,
	}
}

func (c *CreateClusterTool) Handler(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[CreateClusterRequest],
) (*mcp.CallToolResultFor[any], error) {
	res, err := c.r.ClusterService().CreateCluster(ctx, connect.NewRequest(&v1alpha1.CreateClusterRequest{
		Cluster: &v1alpha1.Cluster{
			DisplayName: params.Arguments.DisplayName,
			Description: params.Arguments.Description,
		},
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Cluster creation successfully started at operation `%s`.", res.Msg.GetName()),
			},
		},
	}, nil
}
