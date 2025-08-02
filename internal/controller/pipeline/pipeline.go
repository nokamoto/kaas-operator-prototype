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
		reconciler: newReconciler(client, opts),
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

	// Steps for creating a KubernetesCluster
	end, res, err := r.forKubernetesCluster(ctx, req, pipeline)
	if !end || err != nil {
		return res, err
	}

	// Steps for creating a KubernetesClusterConfiguration
	end, res, err = r.forKubernetesClusterConfiguration(ctx, req, pipeline)
	if !end || err != nil {
		return res, err
	}

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

func (r *PipelineReconciler) forKubernetesCluster(ctx context.Context, req ctrl.Request, pipeline *v1alpha1.Pipeline) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// Check whether the KubernetesCluster name is set
	if pipeline.Spec.Cluster.Name == "" {
		if err := r.updateStatus(ctx, pipeline, v1alpha1.PipelinePhaseFailed, &metav1.Condition{
			Type:    string(v1alpha1.PipelineConditionTypeFailed),
			Status:  metav1.ConditionFalse,
			Reason:  "ValidationFailed",
			Message: "KubernetesCluster name is not set in the Pipeline spec.",
		}); err != nil {
			return false, ctrl.Result{}, fmt.Errorf("failed to update Pipeline status: %w", err)
		}
		logger.Info("KubernetesCluster name is not set in the Pipeline spec. Failing the Pipeline.")
		return false, ctrl.Result{}, nil
	}
	logger.Info("KubernetesCluster name is set", "name", pipeline.Spec.Cluster.Name)
	// Check if the KubernetesCluster resource exists
	var kubernetesCluster v1alpha1.KubernetesCluster
	if err := r.Get(ctx, client.ObjectKey{Name: pipeline.Spec.Cluster.Name, Namespace: req.Namespace}, &kubernetesCluster); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "failed to get KubernetesCluster")
			return false, ctrl.Result{}, fmt.Errorf("failed to get KubernetesCluster: %w", err)
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
			return false, ctrl.Result{}, fmt.Errorf("failed to create KubernetesCluster: %w", err)
		}
		// immediately requeue to poll the status of the KubernetesCluster
		return false, ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	// KubernetesCluster exists, check if it is in running phase
	if kubernetesCluster.Status.Phase != v1alpha1.KubernetesClusterPhaseRunning {
		logger.Info("KubernetesCluster is not running. Waiting for it to be ready.", "phase", kubernetesCluster.Status.Phase)
		return false, ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	logger.Info("KubernetesCluster is running", "name", kubernetesCluster.Name)
	return true, ctrl.Result{}, nil
}

func (r *PipelineReconciler) forKubernetesClusterConfiguration(ctx context.Context, req ctrl.Request, pipeline *v1alpha1.Pipeline) (bool, ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// Check if the KubernetesClusterConfiguration is already created
	name := pipeline.Spec.Cluster.Name
	var kubernetesClusterConfiguration v1alpha1.KubernetesClusterConfiguration
	if err := r.Get(ctx, client.ObjectKey{Name: name, Namespace: req.Namespace}, &kubernetesClusterConfiguration); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "failed to get KubernetesClusterConfiguration")
			return false, ctrl.Result{}, fmt.Errorf("failed to get KubernetesClusterConfiguration: %w", err)
		}
		// KubernetesClusterConfiguration does not exist, create it
		kubernetesClusterConfiguration = v1alpha1.KubernetesClusterConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: req.Namespace,
			},
			Spec: v1alpha1.KubernetesClusterConfigurationSpec{
				Owner: v1alpha1.KubernetesClusterConfigurationSpecOwner{
					Name: name,
				},
			},
		}
		if err := r.Create(ctx, &kubernetesClusterConfiguration); err != nil {
			logger.Error(err, "failed to create KubernetesClusterConfiguration")
			return false, ctrl.Result{}, fmt.Errorf("failed to create KubernetesClusterConfiguration: %w", err)
		}
		logger.Info("KubernetesClusterConfiguration created", "name", kubernetesClusterConfiguration.Name)
		return false, ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	// KubernetesClusterConfiguration exists, check if it is in running phase
	if kubernetesClusterConfiguration.Status.Phase != v1alpha1.KubernetesClusterConfigurationPhaseRunning {
		logger.Info("KubernetesClusterConfiguration is not running. Waiting for it to be ready.", "phase", kubernetesClusterConfiguration.Status.Phase)
		return false, ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	}
	logger.Info("KubernetesClusterConfiguration is ready", "name", kubernetesClusterConfiguration.Name, "phase", kubernetesClusterConfiguration.Status.Phase)
	return true, ctrl.Result{}, nil
}

func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("pipeline-controller").
		For(&v1alpha1.Pipeline{}).
		Complete(r)
}
