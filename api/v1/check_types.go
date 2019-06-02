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
	corev1 "k8s.io/apimachinery/pkg/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CheckSpec defines the desired state of Check
type CheckSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CheckDefinition   CheckDefinition
	CheckRequirements CheckRequirements
}

// CheckStatus defines the observed state of Check
type CheckStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// CheckDefinition defines the type of endpoint to target and the require configuration
type CheckDefinition struct {
	DestinationHostBinding DestinationHostBinding
}

// CheckRequirements defines the connection and resource requirements of the check
type CheckRequirements struct {
	// Define the tls requirements for the check
	// +optional
	TLSSpec TLSSpec

	// Define the resources required to complete this check
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// TLSSpec defines the TLS config required
type TLSSpec struct {
	TlsType TLSType
	// +optional
	CACertRef string `json:"cacert,omitempty"`
	// +optional
	PrivateKeyRef string `json:"private_key,omitempty"`
	// +optional
	CertRef string `json:"cert,omitempty"`
}

type TLSType int

const (
	CLIENT = iota
	MUTUAL
)

func (TLSType t) String() string {
	return []string{"CLIENT", "MUTUAL"}[t]
}

// +kubebuilder:object:root=true

// Check is the Schema for the checks API
type Check struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CheckSpec   `json:"spec,omitempty"`
	Status CheckStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CheckList contains a list of Check
type CheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Check `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Check{}, &CheckList{})
}
