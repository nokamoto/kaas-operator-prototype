package kubernetescluster

import (
	"context"
	"fmt"
	"time"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/boilerplate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters/status,verbs=get;update;patch

type KubernetesClusterReconciler struct {
	client.Client
	status *boilerplate.StatusUpdater[*v1alpha1.KubernetesCluster, v1alpha1.KubernetesClusterPhase]
	opts   KubernetesClusterReconcilerOptions
}

type KubernetesClusterReconcilerOptions struct {
	// PollingInterval is the interval at which the controller will requeue the reconciliation request
	// when the KubernetesCluster is in a non-terminal phase.
	PollingInterval time.Duration
}

func NewKubernetesClusterReconciler(client client.Client, opts KubernetesClusterReconcilerOptions) *KubernetesClusterReconciler {
	return &KubernetesClusterReconciler{
		Client: client,
		status: boilerplate.NewStatusUpdater[*v1alpha1.KubernetesCluster, v1alpha1.KubernetesClusterPhase](client),
		opts:   opts,
	}
}

func (r *KubernetesClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling KubernetesCluster")
	// Fetch the KubernetesCluster instance
	kubernetesCluster := &v1alpha1.KubernetesCluster{}
	if err := r.Get(ctx, req.NamespacedName, kubernetesCluster); err != nil {
		logger.Error(err, "unable to fetch KubernetesCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger = logger.WithValues("phase", kubernetesCluster.Status.Phase)

	switch kubernetesCluster.Status.Phase {
	case v1alpha1.KubernetesClusterPhaseCreating:
		// Simulate the creation of a Kubernetes cluster
		logger.Info("KubernetesCluster is successfully created, setting phase to Running")
		if err := r.status.Update(ctx, kubernetesCluster, v1alpha1.KubernetesClusterPhaseRunning, &metav1.Condition{
			Type:    string(v1alpha1.KubernetesClusterConditionReady),
			Status:  metav1.ConditionTrue,
			Reason:  "KubernetesClusterCreated",
			Message: "KubernetesCluster is successfully created and ready to use",
		}); err != nil {
			logger.Error(err, "failed to update KubernetesCluster status")
			return ctrl.Result{}, fmt.Errorf("failed to update KubernetesCluster status: %w", err)
		}
		logger.Info("KubernetesCluster status updated to Running")
		return ctrl.Result{}, nil

	case v1alpha1.KubernetesClusterPhaseRunning:
		logger.Info("unimplemented yet")
		return ctrl.Result{}, nil

	case v1alpha1.KubernetesClusterPhaseDeleting:
		logger.Info("unimplemented yet")
		return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil

	default:
		// If the phase is not recognized, we set it to Creating
		logger.Info("initializing kubernetes cluster, setting phase to Creating")
		if err := r.status.Update(ctx, kubernetesCluster, v1alpha1.KubernetesClusterPhaseCreating, &metav1.Condition{
			Type:    string(v1alpha1.KubernetesClusterConditionReady),
			Status:  metav1.ConditionTrue,
			Reason:  "KubernetesClusterInitializing",
			Message: "KubernetesCluster is being initialized",
		}); err != nil {
			logger.Error(err, "failed to update KubernetesCluster status")
			return ctrl.Result{}, fmt.Errorf("failed to update KubernetesCluster status: %w", err)
		}
		logger.Info("KubernetesCluster status updated to Creating")
		return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
}

func (r *KubernetesClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("kubernetescluster-controller").
		For(&v1alpha1.KubernetesCluster{}).
		Complete(r)
}
