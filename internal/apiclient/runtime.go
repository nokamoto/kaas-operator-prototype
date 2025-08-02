package apiclient

import (
	"net/http"

	"github.com/nokamoto/kaas-operator-prototype/pkg/api/proto/v1alpha1/v1alpha1connect"
)

// DefaultBaseURL is the default base URL for the API client.
// The default value is set to "http://localhost:8080", which is commonly used for local development.
const DefaultBaseURL = "http://localhost:8080"

// Runtime provides methods to create API clients for the kaas-operator-prototype.
// It uses a lazy evaluation for the base URL to allow dynamic configuration.
type Runtime struct {
	lazyBaseURL func() string
}

// New creates a new Runtime instance with a lazy evaluation function for the base URL.
func New(lazyBaseURL func() string) *Runtime {
	return &Runtime{lazyBaseURL: lazyBaseURL}
}

func (r *Runtime) httpClient() *http.Client {
	return &http.Client{}
}

func (r *Runtime) ClusterService() v1alpha1connect.ClusterServiceClient {
	return v1alpha1connect.NewClusterServiceClient(
		r.httpClient(),
		r.lazyBaseURL(),
	)
}
