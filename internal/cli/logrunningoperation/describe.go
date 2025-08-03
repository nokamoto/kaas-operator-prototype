package logrunningoperation

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/nokamoto/kaas-operator-prototype/internal/cli/encode"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/spf13/cobra"
)

func newGet(r runtime) *cobra.Command {
	var out encode.Encoder
	cmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get a long-running operation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			service := r.LongRunningOperationService()
			res, err := service.GetOperation(cmd.Context(), connect.NewRequest(&v1alpha1.GetOperationRequest{
				Name: name,
			}))
			if err != nil {
				return fmt.Errorf("failed to describe long-running operation: %w", err)
			}
			out.Print(cmd, res.Msg)
			return nil
		},
	}
	out.VarP(cmd)
	return cmd
}
