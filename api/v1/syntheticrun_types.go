/*

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SyntheticRunSpec defines the desired state of SyntheticRun
type SyntheticRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	InstanceCount int32 `json:"instanceCount,omitempty"`

	// +kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	// +kubebuilder:validation:MinLength=0
	Endpoint string `json:"endpoint"`

	// +kubebuilder:validation:MinLength=0
	Path string `json:"path"`

	//TODO make this a mutating web hook
	JobID string `json:"jobID"`
}

// SyntheticRunStatus defines the observed state of SyntheticRun
type SyntheticRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Agents []corev1.Pod `json:"agents,omitempty"`

	JobID string `json:"jobID"`

	ConfigVersion int32 `json:"configVersion"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SyntheticRun is the Schema for the syntheticruns API
type SyntheticRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyntheticRunSpec   `json:"spec,omitempty"`
	Status SyntheticRunStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SyntheticRunList contains a list of SyntheticRun
type SyntheticRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyntheticRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyntheticRun{}, &SyntheticRunList{})
}
