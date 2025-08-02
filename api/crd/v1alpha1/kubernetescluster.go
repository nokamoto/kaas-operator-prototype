package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubernetesclusters,scope=Namespaced
// +kubebuilder:resource:shortName=kc
type KubernetesCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterSpec   `json:"spec,omitempty"`
	Status KubernetesClusterStatus `json:"status,omitempty"`
}

type KubernetesClusterSpec struct{}

type KubernetesClusterPhase string

const (
	// KubernetesClusterPhaseCreating indicates that the Kubernetes cluster is being created.
	KubernetesClusterPhaseCreating KubernetesClusterPhase = "Creating"
	// KubernetesClusterPhaseRunning indicates that the Kubernetes cluster is currently running.
	KubernetesClusterPhaseRunning KubernetesClusterPhase = "Running"
	// KubernetesClusterPhaseDeleting indicates that the Kubernetes cluster is being deleted.
	KubernetesClusterPhaseDeleting KubernetesClusterPhase = "Deleting"
)

type KubernetesClusterConditionType string

const (
	// KubernetesClusterConditionReady indicates that the Kubernetes cluster is ready to be used.
	KubernetesClusterConditionReady KubernetesClusterConditionType = "Ready"
	// KubernetesClusterConditionFailed indicates that the Kubernetes cluster has failed.
	KubernetesClusterConditionFailed KubernetesClusterConditionType = "Failed"
)

type KubernetesClusterStatus struct {
	Phase          KubernetesClusterPhase `json:"phase,omitempty"`
	Conditions     []metav1.Condition     `json:"conditions,omitempty"`
	LastSyncedTime metav1.Time            `json:"lastSyncedTime,omitempty"`
}

func (obj *KubernetesCluster) SetPhase(s string) {
	obj.Status.Phase = KubernetesClusterPhase(s)
}

func (obj *KubernetesCluster) AddCondition(v metav1.Condition) {
	obj.Status.Conditions = append(obj.Status.Conditions, v)
}

func (obj *KubernetesCluster) SetLastSyncedTime(t metav1.Time) {
	obj.Status.LastSyncedTime = t
}

// +kubebuilder:object:root=true
type KubernetesClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubernetesCluster{}, &KubernetesClusterList{})
}

var KubernetesClusterGVR = GroupVersion.WithResource("kubernetesclusters")

const (
	KubernetesClusterAnnotationDisplayName = "nokamoto.github.com/kubernetescluster.displayName"
	KubernetesClusterAnnotationDescription = "nokamoto.github.com/kubernetescluster.description"
)
