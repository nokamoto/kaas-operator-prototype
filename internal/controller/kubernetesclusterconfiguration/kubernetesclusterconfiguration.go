package kubernetesclusterconfiguration

import (
	"context"
	"time"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/boilerplate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusterconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusterconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters,verbs=get;list;watch;create;update;patch;delete

type KubernetesClusterConfigurationReconciler struct {
	client.Client
	status *boilerplate.StatusUpdater[*v1alpha1.KubernetesClusterConfiguration, v1alpha1.KubernetesClusterConfigurationPhase]
	opts   KubernetesClusterConfigurationReconcilerOptions
}

type KubernetesClusterConfigurationReconcilerOptions struct {
	// PollingInterval is the interval at which the controller will requeue the reconciliation request
	// when the KubernetesClusterConfiguration is in a non-terminal phase.
	PollingInterval time.Duration
}

func NewKubernetesClusterConfigurationReconciler(client client.Client, opts KubernetesClusterConfigurationReconcilerOptions) *KubernetesClusterConfigurationReconciler {
	return &KubernetesClusterConfigurationReconciler{
		Client: client,
		status: boilerplate.NewStatusUpdater[*v1alpha1.KubernetesClusterConfiguration, v1alpha1.KubernetesClusterConfigurationPhase](client),
		opts:   opts,
	}
}

func (r *KubernetesClusterConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling KubernetesClusterConfiguration")
	// Fetch the KubernetesClusterConfiguration instance
	kubernetesClusterConfiguration := &v1alpha1.KubernetesClusterConfiguration{}
	if err := r.Get(ctx, req.NamespacedName, kubernetesClusterConfiguration); err != nil {
		logger.Error(err, "unable to fetch KubernetesClusterConfiguration")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger = logger.WithValues("phase", kubernetesClusterConfiguration.Status.Phase)
	switch kubernetesClusterConfiguration.Status.Phase {
	case v1alpha1.KubernetesClusterConfigurationPhaseCreating:
	case v1alpha1.KubernetesClusterConfigurationPhaseRunning:
	default:
		// If the phase is not recognized, set it to creating
		logger.Info("KubernetesClusterConfiguration phase is not recognized, setting it to Creating")
		if err := r.status.Update(ctx, kubernetesClusterConfiguration, v1alpha1.KubernetesClusterConfigurationPhaseCreating, &metav1.Condition{
			Type:    string(v1alpha1.KubernetesClusterConfigurationConditionReady),
			Status:  metav1.ConditionTrue,
			Reason:  "KubernetesClusterConfigurationInitializing",
			Message: "KubernetesClusterConfiguration is initializing",
		}); err != nil {
			logger.Error(err, "failed to update KubernetesClusterConfiguration status")
			return ctrl.Result{}, err
		}
		// Requeue soon to process the Creating phase
		return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	return ctrl.Result{}, nil
}

func (r *KubernetesClusterConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("kubernetesclusterconfiguration-controller").
		For(&v1alpha1.KubernetesClusterConfiguration{}).
		Complete(r)
}
