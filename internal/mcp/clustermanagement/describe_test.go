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
	"google.golang.org/protobuf/types/known/anypb"
)

func TestDescribeLongRunningOperationTool_Handler(t *testing.T) {
	testOperationName := "test-operation"
	type testcase struct {
		name    string
		request DescribeLongRunningOperationRequest
		mock    func(*mockv1alpha1.MockLongRunningOperationServiceClient)
		want    *mcp.CallToolResultFor[any]
		wantErr error
	}
	tests := []testcase{
		{
			name: "describe a succeeded operation",
			request: DescribeLongRunningOperationRequest{
				Name: testOperationName,
			},
			mock: func(m *mockv1alpha1.MockLongRunningOperationServiceClient) {
				pipeline := &v1alpha1.LongRunningOperation_Pipeline{
					Status: &v1alpha1.LongRunningOperation_Pipeline_Status{
						Phase: "ok",
					},
				}
				metadata, _ := anypb.New(pipeline)
				m.EXPECT().GetOperation(gomock.Any(), connect.NewRequest(&v1alpha1.GetOperationRequest{
					Name: testOperationName,
				})).Return(connect.NewResponse(&v1alpha1.LongRunningOperation{
					Name:     testOperationName,
					Done:     true,
					Metadata: metadata,
				}), nil)
			},
			want: &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: `Long-running operation ` + "`test-operation`" + ` is in progress.
Display Name: 
Description: 
Phase: ok
Last Synched Time: n/a
`,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			lro := mockv1alpha1.NewMockLongRunningOperationServiceClient(ctrl)
			if tt.mock != nil {
				tt.mock(lro)
			}

			tool := DescribeLongRunningOperationTool{
				r: &mockRuntime{lro: lro},
			}

			params := &mcp.CallToolParamsFor[DescribeLongRunningOperationRequest]{
				Arguments: tt.request,
			}
			got, err := tool.Handler(context.TODO(), nil, params)
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("DescribeLongRunningOperationTool.Handler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DescribeLongRunningOperationTool.Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("DescribeLongRunningOperationTool.Handler() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
