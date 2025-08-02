package pipeline

import (
	"context"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("PipelineQueueReconciler", func() {
	const testName = "test-pipeline"
	const testNamespace = "test-pipeline-queue-reconciler"

	namespacedName := types.NamespacedName{
		Name:      testName,
		Namespace: testNamespace,
	}

	var pipelineQueueReconciler PipelineQueueReconciler

	BeforeEach(func(ctx context.Context) {
		ns := &corev1.Namespace{}

		By("setting up test namespace")
		ns.Name = testNamespace
		err := k8sClient.Create(ctx, ns)
		Expect(client.IgnoreAlreadyExists(err)).NotTo(HaveOccurred())

		By("initializing the PipelineQueueReconciler")
		pipelineQueueReconciler = PipelineQueueReconciler{
			reconciler: reconciler{
				Client: k8sClient,
				now: func() metav1.Time {
					return now
				},
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
		res, err := pipelineQueueReconciler.Reconcile(ctx, reconcile.Request{
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
		otherPipeline := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-pipeline",
				Namespace: testNamespace,
			},
		}
		err := k8sClient.Create(ctx, &otherPipeline)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &otherPipeline, v1alpha1.PipelinePhaseRunning)

		By("creating a test Pipeline resource and setting it to pending phase")
		got = v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("reconciling the test Pipeline resource to check if it remains in pending phase")
		res, err := pipelineQueueReconciler.Reconcile(ctx, reconcile.Request{
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
		otherPipeline := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other-pipeline",
				Namespace: testNamespace,
			},
		}
		err := k8sClient.Create(ctx, &otherPipeline)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &otherPipeline, v1alpha1.PipelinePhasePending)

		By("creating a test Pipeline resource and setting it to pending phase")
		got = v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("reconciling the test Pipeline resource to check if it remains in pending phase")
		res, err := pipelineQueueReconciler.Reconcile(ctx, reconcile.Request{
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
		got = v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		err := k8sClient.Create(ctx, &got)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &got, v1alpha1.PipelinePhasePending)

		By("creating another Pipeline resource in pending phase")
		otherPipeline := v1alpha1.Pipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "z-other-pipeline", // Ensure the name is lexicographically after the test pipeline in case of exact same creation time
				Namespace: testNamespace,
			},
		}
		err = k8sClient.Create(ctx, &otherPipeline)
		Expect(err).NotTo(HaveOccurred())
		updateStatus(ctx, &otherPipeline, v1alpha1.PipelinePhasePending)

		By("reconciling the test Pipeline resource to check if it is set to running phase")
		res, err := pipelineQueueReconciler.Reconcile(ctx, reconcile.Request{
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
