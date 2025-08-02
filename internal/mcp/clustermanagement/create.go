package clustermanagement

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
)

type runtime interface {
	ClusterService() v1alpha1connect.ClusterServiceClient
}

// ClusterCreateRequest is the request for creating a KaaS cluster.
type ClusterCreateRequest struct {
	DisplayName string `json:"display_name" jsonschema:"optional. The display name of the cluster."`
	Description string `json:"description" jsonschema:"optional. The description of the cluster."`
}

// CreateTool is a tool for creating KaaS clusters.
type CreateTool struct {
	r runtime
}

func New(r runtime) *CreateTool {
	return &CreateTool{r: r}
}

func (c *CreateTool) Tool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_kaas_cluster",
		Description: `Manage KaaS clusters.
This tool provides create operations for KaaS clusters.
Note:
Cluster creation typically takes serveral minutes to complete.`,
	}
}

func (c *CreateTool) Handler(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ClusterCreateRequest],
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
				Text: fmt.Sprintf("Cluster creation sucessfully started at operation `%s`.", res.Msg.GetName()),
			},
		},
	}, nil
}
