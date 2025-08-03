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

// TypedClient provides typed access to Kubernetes resources.
//
// Get methods return the resource type directly, or ErrResourceNotFound if the resource does not exist.
type TypedClient struct {
	pl  objectClient[*v1alpha1.Pipeline]
	kc  objectClient[*v1alpha1.KubernetesCluster]
	kcc objectClient[*v1alpha1.KubernetesClusterConfiguration]
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
	return &TypedClient{
		pl: objectClient[*v1alpha1.Pipeline]{
			client: c,
			typ:    "Pipeline",
		},
		kc: objectClient[*v1alpha1.KubernetesCluster]{
			client: c,
			typ:    "KubernetesCluster",
		},
		kcc: objectClient[*v1alpha1.KubernetesClusterConfiguration]{
			client: c,
			typ:    "KubernetesClusterConfiguration",
		},
	}, nil
}

// CreatePipeline creates a new Pipeline resource in the Kubernetes cluster.
func (c *TypedClient) CreatePipeline(ctx context.Context, pipeline *v1alpha1.Pipeline) error {
	return c.pl.create(ctx, pipeline)
}

// GetPipeline retrieves a Pipeline resource by its name and namespace.
func (c *TypedClient) GetPipeline(ctx context.Context, name, namespace string) (*v1alpha1.Pipeline, error) {
	return c.pl.get(ctx, name, namespace)
}

// GetKubernetesCluster retrieves a KubernetesCluster resource by its name and namespace.
func (c *TypedClient) GetKubernetesCluster(ctx context.Context, name, namespace string) (*v1alpha1.KubernetesCluster, error) {
	return c.kc.get(ctx, name, namespace)
}

// GetKubernetesClusterConfiguration retrieves a KubernetesClusterConfiguration resource by its name and namespace.
func (c *TypedClient) GetKubernetesClusterConfiguration(ctx context.Context, name, namespace string) (*v1alpha1.KubernetesClusterConfiguration, error) {
	return c.kcc.get(ctx, name, namespace)
}
