package cluster

import (
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
	"github.com/spf13/cobra"
)

type runtime interface {
	ClusterService() v1alpha1connect.ClusterServiceClient
}

func New(r runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Manage Kubernetes clusters",
		Aliases: []string{"c"},
	}
	cmd.AddCommand(newCreate(r))
	return cmd
}
