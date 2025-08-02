package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TypedClient struct {
	client client.Client
}

func newDefaultRestConfig() (*rest.Config, error) {
	cfg, err := rest.InClusterConfig()
	if err == nil {
		return cfg, nil
	}
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	cfg, fallbackErr := config.ClientConfig()
	if fallbackErr != nil {
		return nil, fmt.Errorf("failed to create rest config: %w", errors.Join(err, fallbackErr))
	}
	return cfg, nil
}

// New creates a new TypedClient with a default REST configuration.
// It uses in-cluster configuration if available, otherwise it falls back to the kubeconfig file.
func New() (*TypedClient, error) {
	cfg, err := newDefaultRestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create default rest config: %w", err)
	}
	scheme := runtime.NewScheme()
	if err := v1alpha1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add v1alpha1 scheme: %w", err)
	}
	c, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &TypedClient{client: c}, nil
}

// CreatePipeline creates a new Pipeline resource in the Kubernetes cluster.
func (c *TypedClient) CreatePipeline(ctx context.Context, pipeline *v1alpha1.Pipeline) error {
	if err := c.client.Create(ctx, pipeline); err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}
	return nil
}
