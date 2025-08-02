//go:generate mockgen -package cluster -destination mock_cluster_test.go . client,namegen
package cluster

import (
	"context"

	"connectrpc.com/connect"
	typev1alpha1 "github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	apiv1alpha1 "github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultNamespace = "default"

type client interface {
	CreatePipeline(ctx context.Context, pipeline *typev1alpha1.Pipeline) error
}

type namegen interface {
	New(format string, v ...any) string
}

type ClusterService struct {
	v1alpha1connect.UnimplementedClusterServiceHandler
	client  client
	namegen namegen
}

func New(client client, namegen namegen) *ClusterService {
	return &ClusterService{
		client:  client,
		namegen: namegen,
	}
}

// CreateCluster creates a pipeline resource to start a cluster creation operation.
// It returns a LongRunningOperation that can be used to track the progress of the operation.
func (c *ClusterService) CreateCluster(
	ctx context.Context,
	req *connect.Request[apiv1alpha1.CreateClusterRequest],
) (*connect.Response[apiv1alpha1.LongRunningOperation], error) {
	pipeline := &typev1alpha1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.namegen.New("cluster-create"),
			Namespace: defaultNamespace,
		},
	}
	if err := c.client.CreatePipeline(ctx, pipeline); err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}
	return connect.NewResponse(&apiv1alpha1.LongRunningOperation{
		Name: pipeline.Name,
	}), nil
}
