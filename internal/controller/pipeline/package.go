package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=pipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters,verbs=get;list;watch;create;update;patch;delete

type PipelineReconcilerOptions struct {
	// PollingInterval is the interval at which the controller will requeue the reconciliation request
	// when the Pipeline is in a non-terminal phase.
	//
	// If not set, it defaults to 10 seconds.
	PollingInterval time.Duration
}

type reconciler struct {
	client.Client
	opts PipelineReconcilerOptions
}

func newReconciler(c client.Client, opts PipelineReconcilerOptions) reconciler {
	if opts.PollingInterval == 0 {
		opts.PollingInterval = 10 * time.Second
	}
	return reconciler{
		Client: c,
		opts:   opts,
	}
}

func (r *reconciler) updateStatus(ctx context.Context, pipeline *v1alpha1.Pipeline, phase v1alpha1.PipelinePhase, cond *metav1.Condition) error {
	logger := log.FromContext(ctx)
	now := metav1.Now()
	pipeline.Status.Phase = phase
	pipeline.Status.LastSyncedTime = now
	if cond != nil {
		cond.LastTransitionTime = now
		pipeline.Status.Conditions = append(pipeline.Status.Conditions, *cond)
	}
	if err := r.Status().Update(ctx, pipeline); err != nil {
		return fmt.Errorf("failed to update Pipeline status: %w", err)
	}
	logger.Info("Updated Pipeline status", "phase", phase, "lastSyncedTime", now)
	return nil
}
