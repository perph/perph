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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SyntheticRunSpec defines the desired state of SyntheticRun
type SyntheticRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DestinationHost DestinationHostRef `json:"destinationHost,omitempty"`

	// +optional
	Retries int `json:"retries,omitempty"`

	// +optional
	Check CheckRef `json:"check,omitempty"`

	// TODO complete the metrics export later
	//MetricsTask MetricsTaskRef
}

// DestinationHostRef provides a reference to a specific DestinationHost api object
type DestinationHostRef struct {
	//TODO add the reference struct details
}

// CheckRef provides a reference to a specific Check api object
type CheckRef struct {
	//TODO add the reference struct details
}

// SyntheticRunStatus defines the observed state of SyntheticRun
type SyntheticRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status            RunStatus `json:"status,omitempty"`
	CompletionMessage string    `json:"completionMessage,omitempty"`
}

type RunStatus int

const (
	FAILED = iota
	RETRYING
	IN_PROGRESS
	SUCCESS
)

func (t RunStatus) String() string {
	return []string{"FAILED", "RETRYING", "IN_PROGRESS", "SUCCESS"}[t]
}

// +kubebuilder:object:root=true

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
