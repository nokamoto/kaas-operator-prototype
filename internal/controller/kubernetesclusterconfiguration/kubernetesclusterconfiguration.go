package kubernetesclusterconfiguration

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
	kcc := &v1alpha1.KubernetesClusterConfiguration{}
	if err := r.Get(ctx, req.NamespacedName, kcc); err != nil {
		logger.Error(err, "unable to fetch KubernetesClusterConfiguration")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger = logger.WithValues("phase", kcc.Status.Phase)
	switch kcc.Status.Phase {
	case v1alpha1.KubernetesClusterConfigurationPhaseCreating:
		// Create the KubernetesClusterConfigurationConfigMap if it does not exist
		logger.Info("Createing KubernetesClusterConfigurationConfigMap if it does not exist")
		name := kcc.Name
		kccm := &v1alpha1.KubernetesClusterConfigurationConfigMap{}
		if err := r.Get(ctx, client.ObjectKey{Namespace: kcc.Namespace, Name: name}, kccm); err != nil {
			if client.IgnoreNotFound(err) != nil {
				logger.Error(err, "failed to get KubernetesClusterConfigurationConfigMap")
				return ctrl.Result{}, err
			}
			// Create the KubernetesClusterConfigurationConfigMap with owner reference
			kccm = &v1alpha1.KubernetesClusterConfigurationConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: kcc.Namespace,
					Name:      name,
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(kcc, v1alpha1.KubernetesClusterConfigurationGVK),
					},
				},
				Spec: v1alpha1.KubernetesClusterConfigurationConfigMapSpec{
					Name: "testdata",
				},
			}
			logger.Info("KubernetesClusterConfigurationConfigMap does not exist, creating it")
			if err := r.Create(ctx, kccm); err != nil {
				logger.Error(err, "failed to create KubernetesClusterConfigurationConfigMap")
				return ctrl.Result{}, fmt.Errorf("failed to create KubernetesClusterConfigurationConfigMap: %w", err)
			}
			// Requeue to wait for the KubernetesClusterConfigurationConfigMap to be created
			return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
		}
		// Wait for the KubernetesClusterConfigurationConfigMap to be created
		logger.Info("KubernetesClusterConfigurationConfigMap is created, waiting for it to be in Running phase")
		if kccm.Status.Phase != v1alpha1.KubernetesClusterConfigurationPhaseRunning {
			logger.Info("KubernetesClusterConfigurationConfigMap is not in Running phase, requeuing")
			return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
		}
		// Update the KubernetesClusterConfiguration status to Running
		logger.Info("KubernetesClusterConfigurationConfigMap is in Running phase, updating KubernetesClusterConfiguration status")
		if err := r.status.Update(ctx, kcc, v1alpha1.KubernetesClusterConfigurationPhaseRunning, &metav1.Condition{
			Type:    string(v1alpha1.KubernetesClusterConfigurationConditionReady),
			Status:  metav1.ConditionTrue,
			Reason:  "KubernetesClusterConfigurationConfigMapCreated",
			Message: "KubernetesClusterConfigurationConfigMap is successfully created and ready to use",
		}); err != nil {
			logger.Error(err, "failed to update KubernetesClusterConfiguration status")
			return ctrl.Result{}, fmt.Errorf("failed to update KubernetesClusterConfiguration status: %w", err)
		}
		logger.Info("KubernetesClusterConfiguration status updated to Running")
		return ctrl.Result{}, nil

	case v1alpha1.KubernetesClusterConfigurationPhaseRunning:
		logger.Info("KubernetesClusterConfiguration is in Running phase, no action needed")
		return ctrl.Result{}, nil

	default:
		// If the phase is not recognized, set it to creating
		logger.Info("KubernetesClusterConfiguration phase is not recognized, setting it to Creating")
		if err := r.status.Update(ctx, kcc, v1alpha1.KubernetesClusterConfigurationPhaseCreating, &metav1.Condition{
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
}

func (r *KubernetesClusterConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("kubernetesclusterconfiguration-controller").
		For(&v1alpha1.KubernetesClusterConfiguration{}).
		Owns(&v1alpha1.KubernetesClusterConfigurationConfigMap{}).
		Complete(r)
}
