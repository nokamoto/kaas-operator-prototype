package pipeline

import (
	"context"
	"fmt"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PipelineReconciler is responsible for running cluster creation pipelines.
// If the pipeline is in running phase, it will create a KubernetesCluster resource and wait for it to be in the running phase.
type PipelineReconciler struct {
	reconciler
}

func NewPipelineReconciler(client client.Client, opts PipelineReconcilerOptions) *PipelineReconciler {
	return &PipelineReconciler{
		reconciler: reconciler{
			Client: client,
			opts:   opts,
		},
	}
}

func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Pipeline")
	// Fetch the Pipeline instance
	pipeline := &v1alpha1.Pipeline{}
	if err := r.Get(ctx, req.NamespacedName, pipeline); err != nil {
		logger.Error(err, "unable to fetch Pipeline")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if pipeline.Status.Phase != v1alpha1.PipelinePhaseRunning {
		logger.Info("Pipeline is not running. No action required.", "phase", pipeline.Status.Phase)
		return ctrl.Result{}, nil
	}
	// Check whether the KubernetesCluster name is set
	if pipeline.Spec.Cluster.Name == "" {
		if err := r.updateStatus(ctx, pipeline, v1alpha1.PipelinePhaseFailed, &metav1.Condition{
			Type:    string(v1alpha1.PipelineConditionTypeFailed),
			Status:  metav1.ConditionFalse,
			Reason:  "ValidationFailed",
			Message: "KubernetesCluster name is not set in the Pipeline spec.",
		}); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update Pipeline status: %w", err)
		}
		logger.Info("KubernetesCluster name is not set in the Pipeline spec. Failing the Pipeline.")
		return ctrl.Result{}, nil
	}
	logger.Info("KubernetesCluster name is set", "name", pipeline.Spec.Cluster.Name)
	// Check if the KubernetesCluster resource exists
	var kubernetesCluster v1alpha1.KubernetesCluster
	if err := r.Get(ctx, client.ObjectKey{Name: pipeline.Spec.Cluster.Name, Namespace: req.Namespace}, &kubernetesCluster); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "failed to get KubernetesCluster")
			return ctrl.Result{}, fmt.Errorf("failed to get KubernetesCluster: %w", err)
		}
		// KubernetesCluster does not exist, create it
		kubernetesCluster = v1alpha1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pipeline.Spec.Cluster.Name,
				Namespace: req.Namespace,
				Annotations: map[string]string{
					v1alpha1.KubernetesClusterAnnotationDisplayName: pipeline.Spec.Cluster.DisplayName,
					v1alpha1.KubernetesClusterAnnotationDescription: pipeline.Spec.Cluster.Description,
				},
			},
			Spec: v1alpha1.KubernetesClusterSpec{},
		}
		if err := r.Create(ctx, &kubernetesCluster); err != nil {
			logger.Error(err, "failed to create KubernetesCluster")
			return ctrl.Result{}, fmt.Errorf("failed to create KubernetesCluster: %w", err)
		}
		// immediately requeue to poll the status of the KubernetesCluster
		return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	// KubernetesCluster exists, check if it is in running phase
	if kubernetesCluster.Status.Phase != v1alpha1.KubernetesClusterPhaseRunning {
		logger.Info("KubernetesCluster is not running. Waiting for it to be ready.", "phase", kubernetesCluster.Status.Phase)
		return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	logger.Info("KubernetesCluster is running", "name", kubernetesCluster.Name)
	// Update the Pipeline status to indicate that the KubernetesCluster is running
	if err := r.updateStatus(ctx, pipeline, v1alpha1.PipelinePhaseSucceeded, &metav1.Condition{
		Type:    string(v1alpha1.PipelineConditionTypeReady),
		Status:  metav1.ConditionTrue,
		Reason:  "KubernetesClusterRunning",
		Message: "KubernetesCluster is running and Pipeline has succeeded.",
	}); err != nil {
		logger.Error(err, "failed to update Pipeline status")
		return ctrl.Result{}, fmt.Errorf("failed to update Pipeline status: %w", err)
	}
	logger.Info("Pipeline has succeeded")
	return ctrl.Result{}, nil
}

func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("pipeline-controller").
		For(&v1alpha1.Pipeline{}).
		Complete(r)
}
