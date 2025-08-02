package v1alpha1test

import (
	"context"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	"github.com/nokamoto/kaas-operator-prototype/internal/controller/kubernetescluster"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("KubernetesClusterReconciler", func() {
	const testName = "test-cluster"
	const testNamespace = "test-kubernetescluster-reconciler"

	namespacedName := types.NamespacedName{
		Name:      testName,
		Namespace: testNamespace,
	}

	var kubernetesClusterReconciler *kubernetescluster.KubernetesClusterReconciler

	BeforeEach(func(ctx context.Context) {
		ns := &corev1.Namespace{}

		By("setting up test namespace")
		ns.Name = testNamespace
		err := k8sClient.Create(ctx, ns)
		Expect(client.IgnoreAlreadyExists(err)).NotTo(HaveOccurred())

		By("initializing the PipelineReconciler")
		kubernetesClusterReconciler = kubernetescluster.NewKubernetesClusterReconciler(k8sClient, kubernetescluster.KubernetesClusterReconcilerOptions{
			PollingInterval: pollingInterval,
		})
	})

	AfterEach(func(ctx context.Context) {
		By("cleaning up the test namespace")
		err := k8sClient.DeleteAllOf(ctx, &v1alpha1.KubernetesCluster{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should set creating phase if KubernetesCluster is created", func(ctx context.Context) {
		got := &v1alpha1.KubernetesCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		By("creating a test KubernetesCluster")
		err := k8sClient.Create(ctx, got)
		Expect(err).NotTo(HaveOccurred())

		By("reconciling the KubernetesCluster")
		res, err := kubernetesClusterReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))
	})
})
