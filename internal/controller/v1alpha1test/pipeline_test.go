package v1alpha1test

import (
	"context"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/pipeline"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("PipelineReconciler", func() {
	const testName = "test-pipeline"
	const testClusterName = "test-cluster"
	const testNamespace = "test-pipeline-reconciler"

	namespacedName := types.NamespacedName{
		Name:      testName,
		Namespace: testNamespace,
	}

	var pipelineReconciler *pipeline.PipelineReconciler

	BeforeEach(func(ctx context.Context) {
		ns := &corev1.Namespace{}

		By("setting up test namespace")
		ns.Name = testNamespace
		err := k8sClient.Create(ctx, ns)
		Expect(client.IgnoreAlreadyExists(err)).NotTo(HaveOccurred())

		By("initializing the PipelineReconciler")
		pipelineReconciler = pipeline.NewPipelineReconciler(k8sClient, pipeline.PipelineReconcilerOptions{
			PollingInterval: pollingInterval,
		})
	})

	AfterEach(func(ctx context.Context) {
		By("cleaning up the test namespace")
		err := k8sClient.DeleteAllOf(ctx, &v1alpha1.Pipeline{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &v1alpha1.KubernetesCluster{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &v1alpha1.KubernetesClusterConfiguration{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should set failed phase if a Pipeline resource is created without a cluster name", func(ctx context.Context) {
		got := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		By("creating a test Pipeline resource without a cluster name")
		err := k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatusPL(ctx, &got, v1alpha1.PipelinePhaseRunning)

		By("reconciling the Pipeline resource")
		res, err := pipelineReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.IsZero()).To(BeTrue())

		By("verifying the Pipeline resource is in failed phase")
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.PipelinePhaseFailed))
	})

	It("should create a KubernetesCluster resource if not exists", func(ctx context.Context) {
		displayName := "Test Cluster"
		description := "This is a test cluster"
		got := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
			Spec: v1alpha1.PipelineSpec{
				Cluster: v1alpha1.PipelineClusterSpec{
					Name:        testClusterName,
					DisplayName: displayName,
					Description: description,
				},
			},
		}
		By("creating a test Pipeline resource")
		err := k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatusPL(ctx, &got, v1alpha1.PipelinePhaseRunning)

		By("reconciling the Pipeline resource")
		res, err := pipelineReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the KubernetesCluster resource is created")
		var cluster v1alpha1.KubernetesCluster
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      testClusterName,
			Namespace: testNamespace,
		}, &cluster)
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster.ObjectMeta.Annotations[v1alpha1.KubernetesClusterAnnotationDisplayName]).To(Equal(displayName))
		Expect(cluster.ObjectMeta.Annotations[v1alpha1.KubernetesClusterAnnotationDescription]).To(Equal(description))
	})

	It("should create a KubernetesClusterConfiguration resource if not exists", func(ctx context.Context) {
		got := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
			Spec: v1alpha1.PipelineSpec{
				Cluster: v1alpha1.PipelineClusterSpec{
					Name: testClusterName,
				},
			},
		}

		By("creating a test Pipeline resource")
		err := k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatusPL(ctx, &got, v1alpha1.PipelinePhaseRunning)

		By("creating a KubernetesCluster resource in running phase")
		kc := &v1alpha1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testClusterName,
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, kc)
		Expect(err).NotTo(HaveOccurred())
		updateStatusKC(ctx, kc, v1alpha1.KubernetesClusterPhaseRunning)

		By("reconciling the Pipeline resource")
		res, err := pipelineReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the KubernetesClusterConfiguration resource is created")
		var kcc v1alpha1.KubernetesClusterConfiguration
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      testClusterName,
			Namespace: testNamespace,
		}, &kcc)
		Expect(err).NotTo(HaveOccurred())
		Expect(kcc.Spec.Owner.Name).To(Equal(testClusterName))
	})

	It("should succeed if a KubernetesCluster and a KubernetesClusterConfiguration are both in running phase", func(ctx context.Context) {
		got := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
			Spec: v1alpha1.PipelineSpec{
				Cluster: v1alpha1.PipelineClusterSpec{
					Name: testClusterName,
				},
			},
		}

		By("creating a test Pipeline resource")
		err := k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatusPL(ctx, &got, v1alpha1.PipelinePhaseRunning)

		By("creating a KubernetesCluster resource in running phase")
		kc := &v1alpha1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testClusterName,
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, kc)
		Expect(err).NotTo(HaveOccurred())
		updateStatusKC(ctx, kc, v1alpha1.KubernetesClusterPhaseRunning)

		By("creating a KubernetesClusterConfiguration resource in running phase")
		kcc := &v1alpha1.KubernetesClusterConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testClusterName,
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, kcc)
		Expect(err).NotTo(HaveOccurred())
		updateStatusKCC(ctx, kcc, v1alpha1.KubernetesClusterConfigurationPhaseRunning)

		By("reconciling the Pipeline resource")
		res, err := pipelineReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.IsZero()).To(BeTrue())

		By("verifying the Pipeline resource is in succeeded phase")
		err = k8sClient.Get(ctx, namespacedName, &got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.PipelinePhaseSucceeded))
	})
})
