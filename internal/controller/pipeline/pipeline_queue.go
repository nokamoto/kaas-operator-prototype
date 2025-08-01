package pipeline

import (
	"context"
	"fmt"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/boilerplate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PipelineQueueReconciler reconciles a Pipeline object.
// This controller is responsible for managing the queue of pipelines in a Kubernetes cluster.
// It ensures that only one pipeline is running at a time within a namespace.
// If no pipelines are running, it will start the next one in the queue.
type PipelineQueueReconciler struct {
	client.Client
	opts   PipelineReconcilerOptions
	status *boilerplate.StatusUpdater[*v1alpha1.Pipeline, v1alpha1.PipelinePhase]
}

func NewPipelineQueueReconciler(client client.Client, opts PipelineReconcilerOptions) *PipelineQueueReconciler {
	return &PipelineQueueReconciler{
		Client: client,
		opts:   opts,
		status: boilerplate.NewStatusUpdater[*v1alpha1.Pipeline, v1alpha1.PipelinePhase](client),
	}
}

func (r *PipelineQueueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling Pipeline")
	// Fetch the Pipeline instance
	pipeline := &v1alpha1.Pipeline{}
	if err := r.Get(ctx, req.NamespacedName, pipeline); err != nil {
		logger.Error(err, "unable to fetch Pipeline")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	switch pipeline.Status.Phase {
	case v1alpha1.PipelinePhaseRunning:
		logger.Info("Pipeline is currently running. Waiting for it to complete.")
	case v1alpha1.PipelinePhaseFailed, v1alpha1.PipelinePhaseSucceeded:
		logger.Info("Pipeline has completed. No further action required.", "phase", pipeline.Status.Phase)
	case v1alpha1.PipelinePhasePending:
		logger.Info("Pipeline is pending. Check if it can be started.")
		if err := r.reconcile(ctx, pipeline); err != nil {
			logger.Error(err, "failed to reconcile Pipeline")
			return ctrl.Result{}, fmt.Errorf("failed to reconcile Pipeline: %w", err)
		}
		// requeue immediately to check the status again
		return ctrl.Result{RequeueAfter: r.opts.PollingInterval}, nil
	default:
		logger.Info("Pipeline is in an unknown phase. Set to Pending to start processing.")
		if err := r.status.Update(ctx, pipeline, v1alpha1.PipelinePhasePending, &metav1.Condition{
			Type:    string(v1alpha1.PipelineConditionTypeReady),
			Status:  metav1.ConditionTrue,
			Reason:  "PipelinePhasePending",
			Message: "Pipeline is now pending and waiting to be processed.",
		}); err != nil {
			logger.Error(err, "failed to update Pipeline status to Pending")
			return ctrl.Result{}, fmt.Errorf("failed to update Pipeline status: %w", err)
		}
	}
	return ctrl.Result{}, nil
}

func (r *PipelineQueueReconciler) reconcile(ctx context.Context, pipeline *v1alpha1.Pipeline) error {
	logger := log.FromContext(ctx)
	// list all pipelines in the namespace
	pipelineList := &v1alpha1.PipelineList{}
	if err := r.List(ctx, pipelineList, client.InNamespace(pipeline.Namespace)); err != nil {
		return fmt.Errorf("failed to list Pipelines: %w", err)
	}
	// check if the pipeline is first in the queue
	var waitList []*v1alpha1.Pipeline
	for _, p := range pipelineList.Items {
		logger.Info("Checking Pipeline in queue", "name", p.Name, "phase", p.Status.Phase, "creationTimestamp", p.CreationTimestamp)
		if p.Name == pipeline.Name {
			continue
		}
		switch p.Status.Phase {
		case v1alpha1.PipelinePhaseRunning:
			waitList = append(waitList, &p)
		case v1alpha1.PipelinePhasePending:
			if p.CreationTimestamp.Before(&pipeline.CreationTimestamp) {
				// if the pending pipeline was created before the current one, it should be processed first
				waitList = append(waitList, &p)
			}
			if p.CreationTimestamp.Equal(&pipeline.CreationTimestamp) {
				// if the pending pipeline was created at the same time, compare by name
				if p.Name < pipeline.Name {
					waitList = append(waitList, &p)
				}
			}
		}
	}
	if len(waitList) > 0 {
		logger.Info("There are pipelines waiting in the queue", "count", len(waitList))
		return nil
	}
	// if the current pipeline is the first in the queue, start processing it
	if err := r.status.Update(ctx, pipeline, v1alpha1.PipelinePhaseRunning, &metav1.Condition{
		Type:    string(v1alpha1.PipelineConditionTypeReady),
		Status:  metav1.ConditionTrue,
		Reason:  "PipelinePhaseRunning",
		Message: "Pipeline is now running.",
	}); err != nil {
		return fmt.Errorf("failed to update Pipeline status to running: %w", err)
	}
	logger.Info("Pipeline is now running")
	return nil
}

func (r *PipelineQueueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("pipeline-queue-controller").
		For(&v1alpha1.Pipeline{}).
		Complete(r)
}
