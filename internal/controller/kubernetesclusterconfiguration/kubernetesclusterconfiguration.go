package kubernetesclusterconfiguration

import (
	"context"
	"time"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusterconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusterconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters,verbs=get;list;watch;create;update;patch;delete

type KubernetesClusterConfigurationReconciler struct {
	client.Client
	opts KubernetesClusterConfigurationReconcilerOptions
}

type KubernetesClusterConfigurationReconcilerOptions struct {
	// PollingInterval is the interval at which the controller will requeue the reconciliation request
	// when the KubernetesClusterConfiguration is in a non-terminal phase.
	// If not set, it defaults to 10 seconds.
	PollingInterval time.Duration
}

func NewKubernetesClusterConfigurationReconciler(client client.Client, opts KubernetesClusterConfigurationReconcilerOptions) *KubernetesClusterConfigurationReconciler {
	if opts.PollingInterval == 0 {
		opts.PollingInterval = 10 * time.Second
	}
	return &KubernetesClusterConfigurationReconciler{
		Client: client,
		opts:   opts,
	}
}

func (r *KubernetesClusterConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling KubernetesClusterConfiguration")
	return ctrl.Result{}, nil
}

func (r *KubernetesClusterConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("kubernetesclusterconfiguration-controller").
		For(&v1alpha1.KubernetesClusterConfiguration{}).
		Complete(r)
}
