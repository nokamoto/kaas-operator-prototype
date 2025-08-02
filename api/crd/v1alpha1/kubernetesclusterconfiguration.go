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

type KubernetesClusterConfigurationStatus struct{}

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
