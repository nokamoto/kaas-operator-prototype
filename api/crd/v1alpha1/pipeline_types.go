package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=pipelines,scope=Namespaced
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

type PipelineClusterSpec struct {
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

type PipelineSpec struct {
	Cluster PipelineClusterSpec `json:"cluster,omitempty"`
}

type PipelinePhase string

const (
	// PipelinePhasePending indicates that the pipeline is pending and has not yet started processing.
	PipelinePhasePending PipelinePhase = "Pending"
	// PipelinePhaseRunning indicates that the pipeline is currently running.
	// Pipeline should be in progress only one at a time within a namespace.
	PipelinePhaseRunning PipelinePhase = "Running"
	// PipelinePhaseSucceeded indicates that the pipeline has successfully completed.
	PipelinePhaseSucceeded PipelinePhase = "Succeeded"
	// PipelinePhaseFailed indicates that the pipeline has failed.
	PipelinePhaseFailed PipelinePhase = "Failed"
)

type PipelineConditionType string

const (
	// PipelineConditionTypeReady indicates that the pipeline is ready to be processed.
	PipelineConditionTypeReady PipelineConditionType = "Ready"
	// PipelineConditionTypeFailed indicates that the pipeline has failed.
	PipelineConditionTypeFailed PipelineConditionType = "Failed"
)

type PipelineStatus struct {
	Phase          PipelinePhase      `json:"phase,omitempty"`
	Conditions     []metav1.Condition `json:"conditions,omitempty"`
	LastSyncedTime metav1.Time        `json:"lastSyncedTime,omitempty"`
}

// +kubebuilder:object:root=true
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}

var PipelineGVR = GroupVersion.WithResource("pipelines")
