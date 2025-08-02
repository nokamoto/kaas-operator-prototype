package v1alpha1test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestPipelineReconciler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite v1alpha1")
}

var (
	testEnv         *envtest.Environment
	k8sClient       client.Client
	pollingInterval = 1 * time.Second
)

var _ = BeforeSuite(func() {
	By("setting up test environment")
	log.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	fromRoot := func(elem ...string) string {
		return filepath.Join(append([]string{"..", "..", ".."}, elem...)...)
	}

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			fromRoot("config", "crd"),
		},
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())

	scheme := runtime.NewScheme()
	fns := []func(*runtime.Scheme) error{
		v1alpha1.SchemeBuilder.AddToScheme,
		clientgoscheme.AddToScheme,
	}
	for _, fn := range fns {
		err = fn(scheme)
		Expect(err).NotTo(HaveOccurred())
	}

	k8sClient, err = client.New(cfg, client.Options{
		Scheme: scheme,
	})
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func updateStatusPL(ctx context.Context, pl *v1alpha1.Pipeline, phase v1alpha1.PipelinePhase) {
	pl.Status.Phase = phase
	err := k8sClient.Status().Update(ctx, pl)
	Expect(err).NotTo(HaveOccurred(), "failed to update Pipeline status")
}

func updateStatusKC(ctx context.Context, kc *v1alpha1.KubernetesCluster, phase v1alpha1.KubernetesClusterPhase) {
	kc.Status.Phase = phase
	err := k8sClient.Status().Update(ctx, kc)
	Expect(err).NotTo(HaveOccurred(), "failed to update KubernetesCluster status")
}

func updateStatusKCC(ctx context.Context, kcc *v1alpha1.KubernetesClusterConfiguration, phase v1alpha1.KubernetesClusterConfigurationPhase) {
	kcc.Status.Phase = phase
	err := k8sClient.Status().Update(ctx, kcc)
	Expect(err).NotTo(HaveOccurred(), "failed to update KubernetesClusterConfiguration status")
}
