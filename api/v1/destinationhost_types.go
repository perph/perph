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

// DestinationHostSpec defines the desired state of DestinationHost
type DestinationHostSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	TargetSpec TargetSpec

	// Define the tls requirements for the check
	// +optional
	TLSSpec TLSSpec
}

// DestinationHostStatus defines the observed state of DestinationHost
type DestinationHostStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

// +kubebuilder:object:root=true

// DestinationHost is the Schema for the destinationhosts API
type DestinationHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DestinationHostSpec   `json:"spec,omitempty"`
	Status DestinationHostStatus `json:"status,omitempty"`
}

// TargetSpec defines the target adrress, methods and protocols to use
type TargetSpec struct {
	Address AddressSpec `json:"address"`

	// +optional
	RESTConfig RESTSpec `json:"restConfig,omitempty"`

	// TODO implement gRPC methods later
	//GRPCConfig GRPCSpec `json:"grpcConfig,omitempty"`
}

// RESTSpec defines the configuration options for a REST request
type RESTSpec struct {
	Path   string `json:"path"`
	Method string `json:"method"`

	// +optional
	QueryParams map[string]string `json:"queryParams,omitempty"`

	// +optional
	Headers map[string]string `json:"headers,omitempty"`

	// +optional
	Data []byte `json:"data,omitempty"`
}

// AddressSpec defines the connection details to dial the target
type AddressSpec struct {
	HostName string `json:"hostName"`
	// +optional
	Port int `json:"port,omitempty"`
}

// TLSSpec defines the TLS config required
type TLSSpec struct {
	TLSMode TLSType `json:"tlsMode"`
	// +optional
	CACertRef string `json:"cacert,omitempty"`
	// +optional
	PrivateKeyRef string `json:"privateKey,omitempty"`
	// +optional
	CertRef string `json:"cert,omitempty"`
}

// +kubebuilder:object:root=true

// DestinationHostList contains a list of DestinationHost
type DestinationHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DestinationHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DestinationHost{}, &DestinationHostList{})
}
