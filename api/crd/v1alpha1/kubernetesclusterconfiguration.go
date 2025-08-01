package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubernetesclusterconfigurations,scope=Namespaced
// +kubebuilder:resource:shortName=kcc
type KubernetesClusterConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterConfigurationSpec   `json:"spec,omitempty"`
	Status KubernetesClusterConfigurationStatus `json:"status,omitempty"`
}

type KubernetesClusterConfigurationSpec struct {
	// Owner specifies the owner KubernetesCluster resource for this configuration.
	Owner KubernetesClusterConfigurationSpecOwner `json:"owner,omitempty"`
}

type KubernetesClusterConfigurationSpecOwner struct {
	Name string `json:"name,omitempty"`
}

type KubernetesClusterConfigurationPhase string

const (
	// KubernetesClusterConfigurationPhaseCreating indicates that the configuration is being created.
	KubernetesClusterConfigurationPhaseCreating KubernetesClusterConfigurationPhase = "Creating"
	// KubernetesClusterConfigurationPhaseRunning indicates that the configuration is running.
	KubernetesClusterConfigurationPhaseRunning KubernetesClusterConfigurationPhase = "Running"
)

type KubernetesClusterConfigurationConditionType string

const (
	// KubernetesClusterConfigurationConditionReady indicates that the configuration is ready.
	KubernetesClusterConfigurationConditionReady KubernetesClusterConfigurationConditionType = "Ready"
	// KubernetesClusterConfigurationConditionFailed indicates that the configuration has failed.
	KubernetesClusterConfigurationConditionFailed KubernetesClusterConfigurationConditionType = "Failed"
)

type KubernetesClusterConfigurationStatus struct {
	Phase          KubernetesClusterConfigurationPhase `json:"phase,omitempty"`
	Conditions     []metav1.Condition                  `json:"conditions,omitempty"`
	LastSyncedTime metav1.Time                         `json:"lastSyncedTime,omitempty"`
}

func (obj *KubernetesClusterConfiguration) SetPhase(s string) {
	obj.Status.Phase = KubernetesClusterConfigurationPhase(s)
}

func (obj *KubernetesClusterConfiguration) AddCondition(condition metav1.Condition) {
	obj.Status.Conditions = append(obj.Status.Conditions, condition)
}

func (obj *KubernetesClusterConfiguration) SetLastSyncedTime(t metav1.Time) {
	obj.Status.LastSyncedTime = t
}

// +kubebuilder:object:root=true
type KubernetesClusterConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesClusterConfiguration `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubernetesclusterconfigurationconfigmaps,scope=Namespaced
// +kubebuilder:resource:shortName=kccm
type KubernetesClusterConfigurationConfigMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterConfigurationConfigMapSpec `json:"spec,omitempty"`
	Status KubernetesClusterConfigurationStatus        `json:"status,omitempty"`
}

type KubernetesClusterConfigurationConfigMapSpec struct {
	// Name is the ConfigMap name that holds the configuration.
	Name string `json:"name,omitempty"`
}

// +kubebuilder:object:root=true
type KubernetesClusterConfigurationConfigMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesClusterConfigurationConfigMap `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&KubernetesClusterConfiguration{}, &KubernetesClusterConfigurationList{},
		&KubernetesClusterConfigurationConfigMap{}, &KubernetesClusterConfigurationConfigMapList{},
	)
}
