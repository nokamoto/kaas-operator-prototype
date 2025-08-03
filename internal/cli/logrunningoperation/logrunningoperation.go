package logrunningoperation

import (
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
	"github.com/spf13/cobra"
)

type runtime interface {
	LongRunningOperationService() v1alpha1connect.LongRunningOperationServiceClient
}

func New(r runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logrunningoperation",
		Short:   "Manage long-running operations",
		Aliases: []string{"operation", "lro"},
	}
	cmd.AddCommand(
		newGet(r),
	)
	return cmd
}
