package cluster

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/nokamoto/kaas-operator-prototype/internal/cli/encode"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1"
	"github.com/spf13/cobra"
)

func newCreate(r runtime) *cobra.Command {
	var displayName, description string
	var out encode.Encoder
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := r.ClusterService()
			res, err := service.CreateCluster(cmd.Context(), connect.NewRequest(&v1alpha1.CreateClusterRequest{
				Cluster: &v1alpha1.Cluster{
					DisplayName: displayName,
					Description: description,
				},
			}))
			if err != nil {
				return fmt.Errorf("failed to create cluster: %w", err)
			}
			out.Print(cmd, res.Msg)
			return nil
		},
	}
	cmd.Flags().StringVar(&displayName, "display-name", "", "Display name for the cluster")
	cmd.Flags().StringVar(&description, "description", "", "Description for the cluster")
	out.VarP(cmd)
	return cmd
}
