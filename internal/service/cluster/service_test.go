package cluster

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	typev1alpha1 "github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	apiv1alpha1 "github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterService_CreateCluster(t *testing.T) {
	testPipelineName := "test-cluster"
	testClusterName := "test-kubernetescluster"
	type testcase struct {
		name string
		req  *apiv1alpha1.CreateClusterRequest
		mock func(*Mockclient, *Mocknamegen)
		want *apiv1alpha1.LongRunningOperation
		code connect.Code
	}
	tests := []testcase{
		{
			name: "ok if pipeline creation succeeds",
			req:  &apiv1alpha1.CreateClusterRequest{},
			mock: func(client *Mockclient, namegen *Mocknamegen) {
				gomock.InOrder(
					namegen.EXPECT().New("cluster-create").Return(testPipelineName),
					namegen.EXPECT().New("kubernetescluster").Return(testClusterName),
					client.EXPECT().CreatePipeline(gomock.Any(), &typev1alpha1.Pipeline{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testPipelineName,
							Namespace: "default",
						},
						Spec: typev1alpha1.PipelineSpec{
							Cluster: typev1alpha1.PipelineClusterSpec{
								Name: testClusterName,
							},
						},
					}).Return(nil),
				)
			},
			want: &apiv1alpha1.LongRunningOperation{
				Name: testPipelineName,
			},
		},
		{
			name: "unavailable if pipeline creation fails",
			req:  &apiv1alpha1.CreateClusterRequest{},
			mock: func(client *Mockclient, namegen *Mocknamegen) {
				gomock.InOrder(
					namegen.EXPECT().New(gomock.Any()).Return(testPipelineName),
					namegen.EXPECT().New(gomock.Any()).Return(testClusterName),
					client.EXPECT().CreatePipeline(gomock.Any(), gomock.Any()).Return(errors.New("failed to create pipeline")),
				)
			},
			code: connect.CodeUnavailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := NewMockclient(ctrl)
			namegen := NewMocknamegen(ctrl)
			if tt.mock != nil {
				tt.mock(client, namegen)
			}
			service := New(client, namegen)
			res, err := service.CreateCluster(context.TODO(), connect.NewRequest(tt.req))
			if err != nil {
				if connect.CodeOf(err) != tt.code {
					t.Errorf("CreateCluster() error = %v, wantCode %v", connect.CodeOf(err), tt.code)
				}
				return
			}
			got := res.Msg
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("CreateCluster() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
