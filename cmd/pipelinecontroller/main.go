package main

import (
	"time"

	"github.com/nokamoto/kaas-operator-prototype/internal/controller/boilerplate"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/pipeline"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	opts := pipeline.PipelineReconcilerOptions{
		PollingInterval: 10 * time.Second,
	}

	boilerplate.V1alpha1Controller(
		func(mgr ctrl.Manager) error {
			r := pipeline.NewPipelineQueueReconciler(mgr.GetClient(), opts)
			return r.SetupWithManager(mgr)
		},
		func(mgr ctrl.Manager) error {
			r := pipeline.NewPipelineReconciler(mgr.GetClient(), opts)
			return r.SetupWithManager(mgr)
		},
	)
}
