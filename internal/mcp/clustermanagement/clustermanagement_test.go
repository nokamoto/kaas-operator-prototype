package clustermanagement

import (
	mockv1alpha1 "github.com/nokamoto/kaas-operator-prototype/internal/mock/mock_v1alpha1connect"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
)

type mockRuntime struct {
	c   *mockv1alpha1.MockClusterServiceClient
	lro *mockv1alpha1.MockLongRunningOperationServiceClient
}

func (m *mockRuntime) ClusterService() v1alpha1connect.ClusterServiceClient {
	return m.c
}

func (m *mockRuntime) LongRunningOperationService() v1alpha1connect.LongRunningOperationServiceClient {
	return m.lro
}
