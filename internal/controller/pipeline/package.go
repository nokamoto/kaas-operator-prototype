package pipeline

import (
	"time"
)

// +kubebuilder:rbac:groups=nokamoto.github.com,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=pipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nokamoto.github.com,resources=kubernetesclusterconfigurations,verbs=get;list;watch;create;update;patch;delete

type PipelineReconcilerOptions struct {
	// PollingInterval is the interval at which the controller will requeue the reconciliation request
	// when the Pipeline is in a non-terminal phase.
	//
	// If not set, it defaults to 10 seconds.
	PollingInterval time.Duration
}
