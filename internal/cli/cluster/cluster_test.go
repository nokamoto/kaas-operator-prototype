package cluster

import (
	"bytes"
	"errors"
	"testing"

	"buf.build/go/protoyaml"
	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	mockv1alpha1 "github.com/nokamoto/kaas-operator-prototype/internal/mock/mock_v1alpha1connect"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

type mockRuntime struct {
	client *mockv1alpha1.MockClusterServiceClient
}

func (m *mockRuntime) ClusterService() v1alpha1connect.ClusterServiceClient {
	return m.client
}

type testcase struct {
	name    string
	args    []string
	mock    func(*mockv1alpha1.MockClusterServiceClient)
	want    proto.Message
	wantErr error
}

func TestNew_create(t *testing.T) {
	want := &v1alpha1.LongRunningOperation{
		Name: "operation-123",
	}
	testDisplayName := "test-cluster"
	testDescription := "test description"
	tests := []testcase{
		{
			name: "got long-running operation if create cluster successfully",
			args: []string{
				"create",
				"--display-name", testDisplayName,
				"--description", testDescription,
			},
			mock: func(m *mockv1alpha1.MockClusterServiceClient) {
				m.EXPECT().CreateCluster(gomock.Any(), connect.NewRequest(&v1alpha1.CreateClusterRequest{
					Cluster: &v1alpha1.Cluster{
						DisplayName: testDisplayName,
						Description: testDescription,
					},
				})).Return(connect.NewResponse(want), nil)
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mockv1alpha1.NewMockClusterServiceClient(ctrl)
			if tt.mock != nil {
				tt.mock(m)
			}

			cmd := New(&mockRuntime{
				client: m,
			})
			cmd.SetArgs(tt.args)

			var out bytes.Buffer
			cmd.SetOutput(&out)
			if err := cmd.Execute(); !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}

			var got v1alpha1.LongRunningOperation
			if err := protoyaml.Unmarshal(out.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal output: %v", err)
			}
			if diff := cmp.Diff(tt.want, &got, protocmp.Transform()); diff != "" {
				t.Errorf("New() got = %v, want %v, diff: %s", &got, tt.want, diff)
			}
		})
	}
}
