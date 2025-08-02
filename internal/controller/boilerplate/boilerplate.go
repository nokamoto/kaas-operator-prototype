package boilerplate

import (
	"os"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// V1alpha1Controller sets up the v1alpha1 controller with the provided setup function.
// It initializes the scheme, sets up logging, creates a manager, and starts the controller.
// This function is intended to be used in the main function of the controller.
func V1alpha1Controller(setupWithManager ...func(ctrl.Manager) error) {
	logger := ctrl.Log.WithName("setup")

	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		logger.Error(err, "unable to add client-go scheme")
		os.Exit(1)
	}
	if err := v1alpha1.SchemeBuilder.AddToScheme(scheme); err != nil {
		logger.Error(err, "unable to add v1alpha1 scheme")
		os.Exit(1)
	}

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		logger.Error(err, "unable to start manager")
		os.Exit(1)
	}

	for _, setup := range setupWithManager {
		if err := setup(mgr); err != nil {
			logger.Error(err, "unable to create controller")
			os.Exit(1)
		}
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "problem running manager")
		os.Exit(1)
	}
}
