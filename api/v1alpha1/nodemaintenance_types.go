/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ConditionTypeReady is the Ready  Condition.Type for NodeMaintenance
	ConditionTypeReady string = "Ready"
)

const (
	// ConditionReasonPending means that NodeMaintenance is in Pending state
	ConditionReasonPending string = "Pending"
	// ConditionReasonScheduled means that NodeMaintenance is in Scheduled state
	ConditionReasonScheduled string = "Scheduled"
	// ConditionReasonCordon means that NodeMaintenance is in Cordon state
	ConditionReasonCordon string = "Cordon"
	// ConditionReasonWaitForPodCompletion means that NodeMaintenance is in WaitForPodCompletion state
	ConditionReasonWaitForPodCompletion string = "WaitForPodCompletion"
	// ConditionReasonDraining means that NodeMaintenance is in Draining state
	ConditionReasonDraining string = "Draining"
	// ConditionReasonReady means that NodeMaintenance is in Ready state
	ConditionReasonReady string = "Ready"
	// ConditionReasonAborted means that NodeMaintenance is in Aborted state
	ConditionReasonAborted string = "Aborted"
)

// NodeMaintenanceSpec defines the desired state of NodeMaintenance
type NodeMaintenanceSpec struct {
	// RequestorID MUST follow domain name notation format (https://tools.ietf.org/html/rfc1035#section-2.3.1)
	// It MUST be 63 characters or less, beginning and ending with an alphanumeric
	// character ([a-z0-9A-Z]) with dashes (-), dots (.), and alphanumerics between.
	// caller SHOULD NOT create multiple objects with same requestorID and nodeName.
	// This field identifies the requestor of the operation.
	// +kubebuilder:validation:Pattern=`^([a-z0-9A-Z]([-a-z0-9A-Z]*[a-z0-9A-Z])?(\.[a-z0-9A-Z]([-a-z0-9A-Z]*[a-z0-9A-Z])?)*)$`
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:MinLength=2
	RequestorID string `json:"requestorID"`

	// NodeName is The name of the node that maintenance operation will be performed on
	// creation fails if node obj does not exist (webhook)
	NodeName string `json:"nodeName"`

	// Cordon if set, marks node as unschedulable during maintenance operation
	// +kubebuilder:default=true
	Cordon bool `json:"cordon,omitempty"`

	// WaitForPodCompletion specifies pods via selector to wait for completion before performing drain operation
	// if not provided, will not wait for pods to complete
	WaitForPodCompletion *WaitForPodCompletionSpec `json:"waitForPodCompletion,omitempty"`

	// DrainSpec specifies how a node will be drained. if not provided, no draining will be performed.
	DrainSpec *DrainSpec `json:"drainSpec,omitempty"`
}

// WaitForPodCompletionSpec describes the configuration for waiting on pods completion
type WaitForPodCompletionSpec struct {
	// PodSelector specifies a label selector for the pods to wait for completion
	// For more details on label selectors, see:
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	// +kubebuilder:validation:Optional
	// +kubebuilder:example="app=my-workloads"
	PodSelector string `json:"podSelector,omitempty"`

	// TimeoutSecond specifies the length of time in seconds
	// to wait before giving up on pod termination, zero means infinite
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum:=0
	TimeoutSecond uint32 `json:"timeoutSeconds,omitempty"`
}

// DrainSpec describes configuration for node drain during automatic upgrade
type DrainSpec struct {
	// Force indicates if force draining is allowed
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	Force bool `json:"force,omitempty"`

	// PodSelector specifies a label selector to filter pods on the node that need to be drained
	// For more details on label selectors, see:
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	// +kubebuilder:validation:Optional
	PodSelector string `json:"podSelector,omitempty"`

	// TimeoutSecond specifies the length of time in seconds to wait before giving up drain, zero means infinite
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=300
	// +kubebuilder:validation:Minimum:=0
	TimeoutSecond uint32 `json:"timeoutSeconds,omitempty"`

	// DeleteEmptyDir indicates if should continue even if there are pods using emptyDir
	// (local data that will be deleted when the node is drained)
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	DeleteEmptyDir bool `json:"deleteEmptyDir,omitempty"`

	// PodEvictionFilters specifies filters for pods that need to undergo eviction during drain.
	// if specified. only pods that match PodEvictionFilters will be evicted during drain operation.
	// if unspecified. all non-daemonset pods will be evicted.
	// logical OR is performed between filter entires. logical AND is performed within different filters
	// in a filter entry.
	// +kubebuilder:validation:Optional
	PodEvictionFilters []PodEvictionFiterEntry `json:"podEvictionFilters,omitempty"`
}

// PodEvictionFiterEntry defines filters for Pod evictions during drain operation
type PodEvictionFiterEntry struct {
	// ByResourceNameRegex filters pods by the name of the resources they consume using regex.
	ByResourceNameRegex *string `json:"byResourceNameRegex,omitempty"`
}

// NodeMaintenanceStatus defines the observed state of NodeMaintenance
type NodeMaintenanceStatus struct {
	// Conditions represents observations of NodeMaintenance current state
	// +kubebuilder:validation:Optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchMergeKey:"type" patchStrategy:"merge" protobuf:"bytes,1,rep,name=conditions"`

	// Drain represents the drain status of the node
	Drain *DrainStatus `json:"drain,omitempty"`
}

// DrainStatus represents the status of draining for the node
type DrainStatus struct {
	// TotalPods is the number of pods on the node at the time NodeMaintenance was Scheduled
	TotalPods uint32 `json:"totalPods,omitempty"`

	// EvictionPods is the total number of pods that need to be evicted at the time NodeMaintenance was scheduled
	EvictionPods uint32 `json:"evictionPods,omitempty"`

	// DrainProgress represents the draining progress as percentage
	DrainProgress uint32 `json:"drainProgress,omitempty"`

	// WaitForEviction is the list of namespaced named pods that need to be evicted
	WaitForEviction []string `json:"waitForEviction,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Node",type="string",JSONPath=`.spec.nodeName`
// +kubebuilder:printcolumn:name="Requestor",type="string",JSONPath=`.spec.requestorID`
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=`.status.conditions[?(@.type=='Ready')].status`
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=`.status.conditions[?(@.type=='Ready')].reason`

// NodeMaintenance is the Schema for the nodemaintenances API
type NodeMaintenance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeMaintenanceSpec   `json:"spec,omitempty"`
	Status NodeMaintenanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeMaintenanceList contains a list of NodeMaintenance
type NodeMaintenanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeMaintenance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeMaintenance{}, &NodeMaintenanceList{})
}