//go:generate mockgen -package longrunningoperation -destination mock_longrunningoperation_test.go . client
package longrunningoperation

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	typev1alpha1 "github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/domain"
	apiv1alpha1 "github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultNamespace = "default"

type client interface {
	GetPipeline(ctx context.Context, name, namespace string) (*typev1alpha1.Pipeline, error)
	GetKubernetesCluster(ctx context.Context, name, namespace string) (*typev1alpha1.KubernetesCluster, error)
	GetKubernetesClusterConfiguration(ctx context.Context, name, namespace string) (*typev1alpha1.KubernetesClusterConfiguration, error)
}

type LongRunningOperationService struct {
	v1alpha1connect.UnimplementedLongRunningOperationServiceHandler
	client client
}

func New(client client) *LongRunningOperationService {
	return &LongRunningOperationService{
		client: client,
	}
}

func (l *LongRunningOperationService) GetOperation(
	ctx context.Context,
	req *connect.Request[apiv1alpha1.GetOperationRequest],
) (*connect.Response[apiv1alpha1.LongRunningOperation], error) {
	// Retrieve the pipeline name from the request
	pipelineName := req.Msg.GetName()
	pipeline, err := l.client.GetPipeline(ctx, pipelineName, defaultNamespace)
	if errors.Is(err, domain.ErrResourceNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}
	// Set the metadata baed on the pipeline
	metadata := &apiv1alpha1.LongRunningOperation_Pipeline{
		Namespace: pipeline.Namespace,
		Spec: &apiv1alpha1.LongRunningOperation_Pipeline_Spec{
			Name:        pipeline.Spec.Cluster.Name,
			DisplayName: pipeline.Spec.Cluster.DisplayName,
			Description: pipeline.Spec.Cluster.Description,
		},
		Status: &apiv1alpha1.LongRunningOperation_Pipeline_Status{
			Phase:           string(pipeline.Status.Phase),
			LastSynchedTime: timestamppb.New(pipeline.Status.LastSyncedTime.Time),
		},
	}
	for _, cond := range pipeline.Status.Conditions {
		metadata.Status.Conditions = append(
			metadata.Status.Conditions,
			&apiv1alpha1.LongRunningOperation_Pipeline_Status_Condition{
				Message:            cond.Message,
				LastTransitionTime: timestamppb.New(cond.LastTransitionTime.Time),
			},
		)
	}
	m, err := anypb.New(metadata)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// Set the response based on the pipeline status
	var r *anypb.Any
	switch pipeline.Status.Phase {
	case typev1alpha1.PipelinePhaseSucceeded:
		// If the pipeline is succeeded, we can return Cluster as the response
		kc, err := l.client.GetKubernetesCluster(ctx, pipeline.Spec.Cluster.Name, defaultNamespace)
		if errors.Is(err, domain.ErrResourceNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if err != nil {
			return nil, connect.NewError(connect.CodeUnavailable, err)
		}
		kcc, err := l.client.GetKubernetesClusterConfiguration(ctx, pipeline.Spec.Cluster.Name, defaultNamespace)
		if errors.Is(err, domain.ErrResourceNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if err != nil {
			return nil, connect.NewError(connect.CodeUnavailable, err)
		}
		r, err = anypb.New(newCluster(kc, kcc))
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case typev1alpha1.PipelinePhaseFailed:
		// If the pipeline is failed, we can return an empty response
		r, err = anypb.New(&emptypb.Empty{})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	lro := &apiv1alpha1.LongRunningOperation{
		Name:     pipelineName,
		Metadata: m,
		Response: r,
		Done:     r != nil,
	}
	return connect.NewResponse(lro), nil
}

func newCluster(kc *typev1alpha1.KubernetesCluster, _ *typev1alpha1.KubernetesClusterConfiguration) *apiv1alpha1.Cluster {
	c := apiv1alpha1.Cluster{
		Name:        kc.Name,
		DisplayName: kc.Annotations[typev1alpha1.KubernetesClusterAnnotationDisplayName],
		Description: kc.Annotations[typev1alpha1.KubernetesClusterAnnotationDescription],
	}
	return &c
}
