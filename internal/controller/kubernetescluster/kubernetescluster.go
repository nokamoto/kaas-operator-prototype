package kubernetescluster

import (
	"context"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters/status,verbs=get;update;patch

type KubernetesClusterReconciler struct {
	client.Client
	now func() metav1.Time
}

func NewKubernetesClusterReconciler(client client.Client) *KubernetesClusterReconciler {
	return &KubernetesClusterReconciler{
		Client: client,
		now:    metav1.Now,
	}
}

func (r *KubernetesClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling KubernetesCluster")
	return ctrl.Result{}, nil
}

func (r *KubernetesClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("kubernetescluster-controller").
		For(&v1alpha1.KubernetesCluster{}).
		Complete(r)
}
