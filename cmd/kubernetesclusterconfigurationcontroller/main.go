package main

import (
	"time"

	"github.com/nokamoto/kaas-operator-prototype/internal/controller/boilerplate"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/kubernetesclusterconfiguration"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	boilerplate.V1alpha1Controller(
		func(m ctrl.Manager) error {
			opts := kubernetesclusterconfiguration.KubernetesClusterConfigurationReconcilerOptions{
				PollingInterval: 10 * time.Second,
			}
			r := kubernetesclusterconfiguration.NewKubernetesClusterConfigurationReconciler(m.GetClient(), opts)
			return r.SetupWithManager(m)
		},
	)
}
