package pipeline

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestPipelineReconciler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PipelineReconciler Suite")
}

var (
	testEnv   *envtest.Environment
	k8sClient client.Client
)

var _ = BeforeSuite(func() {
	By("setting up test environment")

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

var _ = Describe("PipelineReconciler", func() {
	const testNamespace = "test-pipeline-reconciler"

	BeforeEach(func(ctx context.Context) {
		By("setting up test namespace")
		ns := &corev1.Namespace{}
		ns.Name = testNamespace
		err := k8sClient.Create(ctx, ns)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func(ctx context.Context) {
		By("cleaning up test namespace")
		ns := &corev1.Namespace{}
		ns.Name = testNamespace
		err := k8sClient.Delete(ctx, ns)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should reconcile a Pipeline resource", func(ctx context.Context) {
		By("creating a Pipeline resource")
		err := k8sClient.Create(ctx, &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pipeline",
				Namespace: testNamespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Skip("not implemented yet")
	})
})
