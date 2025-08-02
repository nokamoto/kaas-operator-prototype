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
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

var _ = Describe("PipelineQueueReconciler", func() {
	const testName = "test-pipeline"
	const testNamespace = "test-pipeline-queue-reconciler"

	namespacedName := types.NamespacedName{
		Name:      testName,
		Namespace: testNamespace,
	}

	now := metav1.Now()
	var reconciler PipelineQueueReconciler

	updateStatus := func(ctx context.Context, pipeline *v1alpha1.Pipeline, phase v1alpha1.PipelinePhase) {
		pipeline.Status.Phase = phase
		err := k8sClient.Status().Update(ctx, pipeline)
		Expect(err).NotTo(HaveOccurred(), "failed to update Pipeline status")
	}

	BeforeEach(func(ctx context.Context) {
		ns := &corev1.Namespace{}

		By("setting up test namespace")
		ns.Name = testNamespace
		err := k8sClient.Create(ctx, ns)
		Expect(client.IgnoreAlreadyExists(err)).NotTo(HaveOccurred())

		By("initializing the PipelineQueueReconciler")
		reconciler = PipelineQueueReconciler{
			Client: k8sClient,
			now: func() metav1.Time {
				return now
			},
		}
	})

	AfterEach(func(ctx context.Context) {
		By("cleaning up the test namespace")
		err := k8sClient.DeleteAllOf(ctx, &v1alpha1.Pipeline{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should set pending phase if a Pipeline resource is created", func(ctx context.Context) {
		By("creating a test Pipeline resource")
		err := k8sClient.Create(ctx, &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("reconciling the test Pipeline resource to set it to pending phase")
		res, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.IsZero()).To(BeTrue())

		By("verifying the Pipeline resource is in pending phase")
		var got v1alpha1.Pipeline
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.PipelinePhasePending))
	})

	It("should not set running phase if other Pipeline is running", func(ctx context.Context) {
		var got v1alpha1.Pipeline
		By("creating another Pipeline resource in running phase")
		otherPipeline := &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-pipeline",
				Namespace: testNamespace,
			},
		}
		err := k8sClient.Create(ctx, otherPipeline)
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      otherPipeline.Name,
			Namespace: otherPipeline.Namespace,
		}, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhaseRunning)

		By("creating a test Pipeline resource and setting it to pending phase")
		err = k8sClient.Create(ctx, &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("reconciling the test Pipeline resource to check if it remains in pending phase")
		res, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the Pipeline resource is still in pending phase")
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.PipelinePhasePending))
	})

	It("should not set running phase if no running pipelines exist but another pending pipeline is first in the queue", func(ctx context.Context) {
		var got v1alpha1.Pipeline
		By("creating another Pipeline resource in pending phase")
		otherPipeline := &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-pipeline",
				Namespace: testNamespace,
			},
		}
		err := k8sClient.Create(ctx, otherPipeline)
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      otherPipeline.Name,
			Namespace: otherPipeline.Namespace,
		}, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("creating a test Pipeline resource and setting it to pending phase")
		err = k8sClient.Create(ctx, &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("reconciling the test Pipeline resource to check if it remains in pending phase")
		res, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the Pipeline resource is still in pending phase")
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.PipelinePhasePending))
	})

	It("should set running phase if the Pipeline is first in the queue", func(ctx context.Context) {
		var got v1alpha1.Pipeline
		By("creating a test Pipeline resource and setting it to pending phase")
		err := k8sClient.Create(ctx, &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("creating another Pipeline resource in pending phase")
		otherPipeline := &v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "z-other-pipeline", // Ensure the name is lexicographically after the test pipeline in case of exact same creation time
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, otherPipeline)
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      otherPipeline.Name,
			Namespace: otherPipeline.Namespace,
		}, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("reconciling the test Pipeline resource to check if it is set to running phase")
		res, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the Pipeline resource is in running phase")
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.PipelinePhaseRunning))
	})
})
