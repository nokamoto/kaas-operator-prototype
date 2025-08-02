package main

import (
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/boilerplate"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/kubernetescluster"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	boilerplate.V1alpha1Controller(
		func(mgr ctrl.Manager) error {
			opts := kubernetescluster.KubernetesClusterReconcilerOptions{}
			r := kubernetescluster.NewKubernetesClusterReconciler(mgr.GetClient(), opts)
			return r.SetupWithManager(mgr)
		},
	)
}
