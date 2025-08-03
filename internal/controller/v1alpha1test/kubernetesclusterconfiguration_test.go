package v1alpha1test

import (
	"context"

	"github.com/nokamoto/kaas-operator-prototype/api/crd/v1alpha1"
	kccm "github.com/nokamoto/kaas-operator-prototype/internal/controller/kubernetesclusterconfiguration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("KubernetesClusterConfigurationReconciler", func() {
	const testName = "test-kubernetescluster-configuration"
	const testNamespace = "test-kubernetescluster-configuration-reconciler"

	namespacedName := types.NamespacedName{
		Name:      testName,
		Namespace: testNamespace,
	}

	var kccReconciler *kccm.KubernetesClusterConfigurationReconciler

	BeforeEach(func(ctx context.Context) {
		ns := &corev1.Namespace{}

		By("setting up test namespace")
		ns.Name = testNamespace
		err := k8sClient.Create(ctx, ns)
		Expect(client.IgnoreAlreadyExists(err)).NotTo(HaveOccurred())

		By("initializing the KubernetesClusterConfigurationReconciler")
		kccReconciler = kccm.NewKubernetesClusterConfigurationReconciler(k8sClient, kccm.KubernetesClusterConfigurationReconcilerOptions{
			PollingInterval: pollingInterval,
		})
	})

	AfterEach(func(ctx context.Context) {
		By("cleaning up the test namespace")
		err := k8sClient.DeleteAllOf(ctx, &v1alpha1.KubernetesClusterConfiguration{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &v1alpha1.KubernetesClusterConfigurationConfigMap{}, client.InNamespace(testNamespace))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should set creating phase if KubernetesClusterConfiguration is created", func(ctx context.Context) {
		got := &v1alpha1.KubernetesClusterConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		By("creating a test KubernetesClusterConfiguration")
		err := k8sClient.Create(ctx, got)
		Expect(err).NotTo(HaveOccurred())

		By("reconciling the KubernetesClusterConfiguration")
		res, err := kccReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the KubernetesClusterConfiguration is in creating phase")
		err = k8sClient.Get(ctx, namespacedName, got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.KubernetesClusterConfigurationPhaseCreating))
	})

	It("should create a ConfigMap resource and set Running phase", func(ctx context.Context) {
		got := &v1alpha1.KubernetesClusterConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testName,
				Namespace: testNamespace,
			},
		}
		By("creating a test KubernetesClusterConfiguration in Creating phase")
		err := k8sClient.Create(ctx, got)
		Expect(err).NotTo(HaveOccurred())
		updateStatusKCC(ctx, got, v1alpha1.KubernetesClusterConfigurationPhaseCreating)

		By("reconciling the KubernetesClusterConfiguration")
		res, err := kccReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.RequeueAfter).To(Equal(pollingInterval))

		By("verifying the ConfigMap resource is created")
		var kccm v1alpha1.KubernetesClusterConfigurationConfigMap
		err = k8sClient.Get(ctx, namespacedName, &kccm)
		Expect(err).NotTo(HaveOccurred())
		Expect(kccm.Spec.Name).To(Equal("testdata"))

		By("setting ConfigMap phase to Running and reconciling again")
		updateStatusKCCM(ctx, &kccm, v1alpha1.KubernetesClusterConfigurationPhaseRunning)
		res, err = kccReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: namespacedName,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.IsZero()).To(BeTrue())

		By("verifying the KubernetesClusterConfiguration is in Running phase")
		err = k8sClient.Get(ctx, namespacedName, got)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status.Phase).To(Equal(v1alpha1.KubernetesClusterConfigurationPhaseRunning))
	})
})
