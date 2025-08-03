package cli

import (
	"github.com/nokamoto/kaas-operator-prototype/internal/apiclient"
	"github.com/nokamoto/kaas-operator-prototype/internal/cli/cluster"
	"github.com/nokamoto/kaas-operator-prototype/internal/cli/logrunningoperation"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var baseURL string
	cmd := &cobra.Command{
		Use:   "kcli",
		Short: "Kubernetes as a Service CLI",
	}
	cmd.PersistentFlags().StringVar(&baseURL, "url", apiclient.DefaultBaseURL, "API endpoint URL")

	r := apiclient.New(func() string {
		return baseURL
	})
	cmd.AddCommand(
		cluster.New(r),
		logrunningoperation.New(r),
	)
	return cmd
}
