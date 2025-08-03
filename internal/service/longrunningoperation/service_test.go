package longrunningoperation

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	typev1alpha1 "github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/domain"
	apiv1alpha1 "github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLongRunningOperationService_GetPipeline(t *testing.T) {
	testPipelineName := "test-pipeline"
	type testcase struct {
		name string
		req  *apiv1alpha1.GetOperationRequest
		mock func(*Mockclient)
		want *apiv1alpha1.LongRunningOperation
		code connect.Code
	}
	must := func(v *anypb.Any, err error) *anypb.Any {
		return v
	}
	now := metav1.Now()
	tests := []testcase{
		{
			name: "ok if pipeline get succeeds",
			req:  &apiv1alpha1.GetOperationRequest{Name: testPipelineName},
			mock: func(client *Mockclient) {
				gomock.InOrder(
					client.EXPECT().GetPipeline(gomock.Any(), testPipelineName, "default").Return(&typev1alpha1.Pipeline{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testPipelineName,
							Namespace: "default",
						},
						Spec: typev1alpha1.PipelineSpec{
							Cluster: typev1alpha1.PipelineClusterSpec{
								Name:        "cluster1",
								DisplayName: "Cluster 1",
								Description: "desc",
							},
						},
						Status: typev1alpha1.PipelineStatus{
							Phase:          typev1alpha1.PipelinePhaseSucceeded,
							LastSyncedTime: now,
						},
					}, nil),
					client.EXPECT().GetKubernetesCluster(gomock.Any(), "cluster1", "default").Return(&typev1alpha1.KubernetesCluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "cluster1",
							Annotations: map[string]string{
								typev1alpha1.KubernetesClusterAnnotationDisplayName: "Cluster 1",
								typev1alpha1.KubernetesClusterAnnotationDescription: "desc",
							},
						},
					}, nil),
					client.EXPECT().GetKubernetesClusterConfiguration(gomock.Any(), "cluster1", "default").Return(&typev1alpha1.KubernetesClusterConfiguration{}, nil),
				)
			},
			want: &apiv1alpha1.LongRunningOperation{
				Name: testPipelineName,
				Done: true,
				Metadata: must(anypb.New(&apiv1alpha1.LongRunningOperation_Pipeline{
					Namespace: "default",
					Spec: &apiv1alpha1.LongRunningOperation_Pipeline_Spec{
						Name:        "cluster1",
						DisplayName: "Cluster 1",
						Description: "desc",
					},
					Status: &apiv1alpha1.LongRunningOperation_Pipeline_Status{
						Phase:           string(typev1alpha1.PipelinePhaseSucceeded),
						LastSynchedTime: timestamppb.New(now.Time),
					},
				})),
				Response: must(anypb.New(&apiv1alpha1.Cluster{
					Name:        "cluster1",
					DisplayName: "Cluster 1",
					Description: "desc",
				})),
			},
		},
		{
			name: "not found if pipeline does not exist",
			req:  &apiv1alpha1.GetOperationRequest{Name: testPipelineName},
			mock: func(client *Mockclient) {
				client.EXPECT().GetPipeline(gomock.Any(), testPipelineName, "default").Return(nil, domain.ErrResourceNotFound)
			},
			code: connect.CodeNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := NewMockclient(ctrl)
			if tt.mock != nil {
				tt.mock(client)
			}
			service := New(client)
			res, err := service.GetOperation(context.TODO(), connect.NewRequest(tt.req))
			if err != nil {
				if connect.CodeOf(err) != tt.code {
					t.Errorf("GetOperation() error = %v, wantCode %v", connect.CodeOf(err), tt.code)
				}
				return
			}
			got := res.Msg
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetOperation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
