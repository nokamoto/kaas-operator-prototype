package cli

import (
	"net/http"

	"github.com/nokamoto/kaas-operator-prototype/internal/cli/cluster"
	"github.com/nokamoto/kaas-operator-prototype/internal/cli/logrunningoperation"
	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
	"github.com/spf13/cobra"
)

type runtime struct {
	baseURL string
}

func (r *runtime) httpClient() *http.Client {
	return &http.Client{}
}

func (r *runtime) ClusterService() v1alpha1connect.ClusterServiceClient {
	return v1alpha1connect.NewClusterServiceClient(
		r.httpClient(),
		r.baseURL,
	)
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kcli",
		Short: "Kubernetes as a Service CLI",
	}
	r := &runtime{}
	cmd.PersistentFlags().StringVar(&r.baseURL, "url", "http://localhost:8080", "API endpoint URL")

	cmd.AddCommand(
		cluster.New(r),
		logrunningoperation.New(),
	)
	return cmd
}
