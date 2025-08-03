package clustermanagement

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	mockv1alpha1 "github.com/nokamoto/kaas-operator-prototype/internal/mock/mock_v1alpha1connect"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"go.uber.org/mock/gomock"
)

func TestCreateClusterTool_Handler(t *testing.T) {
	testDisplayName := "test-cluster"
	testDescription := "This is a test cluster."
	testOperationName := "test-operation"
	type testcase struct {
		name    string
		request CreateClusterRequest
		mock    func(*mockv1alpha1.MockClusterServiceClient)
		want    *mcp.CallToolResultFor[any]
		wantErr error
	}
	tests := []testcase{
		{
			name: "cluster creation successfully started",
			request: CreateClusterRequest{
				DisplayName: testDisplayName,
				Description: testDescription,
			},
			mock: func(m *mockv1alpha1.MockClusterServiceClient) {
				m.EXPECT().CreateCluster(gomock.Any(), connect.NewRequest(&v1alpha1.CreateClusterRequest{
					Cluster: &v1alpha1.Cluster{
						DisplayName: testDisplayName,
						Description: testDescription,
					},
				})).Return(connect.NewResponse(&v1alpha1.LongRunningOperation{
					Name: testOperationName,
				}), nil)
			},
			want: &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: "Cluster creation successfully started at operation `test-operation`.",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mockv1alpha1.NewMockClusterServiceClient(ctrl)
			if tt.mock != nil {
				tt.mock(m)
			}

			tool := CreateClusterTool{
				r: &mockRuntime{c: m},
			}

			params := &mcp.CallToolParamsFor[CreateClusterRequest]{
				Arguments: tt.request,
			}
			got, err := tool.Handler(context.TODO(), nil, params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateTool.Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("CreateTool.Handler() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
