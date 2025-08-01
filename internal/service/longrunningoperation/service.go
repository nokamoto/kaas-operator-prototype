package longrunningoperation

import "github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"

type LongRunningOperationService struct {
	v1alpha1connect.UnimplementedLongRunningOperationServiceHandler
}
