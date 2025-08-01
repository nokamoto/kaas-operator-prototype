package pipeline

import (
	"context"

	"github.com/nokamoto/kaas-operator-prototype/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=pipelines/status,verbs=get;update;patch

type PipelineReconciler struct {
	client.Client
}

func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Pipeline", "namespace", req.Namespace, "name", req.Name)
	return ctrl.Result{}, nil
}

func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Pipeline{}).
		Complete(r)
}
