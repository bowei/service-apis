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

package v1alpha1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GatewaySpec defines the desired state of Gateway.
//
// The Spec is split into two major pieces: listeners describing
// client-facing properties and routes that describe application-level
// routing.
//
// Not all possible combinations of options specified in the Spec are
// valid. Some invalid configurations can be caught sychronously via a
// webhook, but there are many cases that will require asynchronous
// signaling via the GatewayStatus block.
//
//
//
type GatewaySpec struct {
	// Class used for this Gateway. This is the name of a GatewayClass resource.
	Class string `json:"class"`
	// Listeners associated with this Gateway. Listeners define what addresses,
	// ports, protocols are bound on this Gateway.
	Listeners []Listener
	// Routes associated with this Gateway. Routes define
	// protocol-specific routing to backends (e.g. Services).
	Routes []core.TypedLocalObjectReference
}

// Listener defines a
type Listener struct {
	// Address bound on the listener. This is optional and behavior
	// can depend on GatewayClass. If a value is set in the spec and
	// the request address is invalid, the GatewayClass MUST indicate
	// this in the associated entry in GatewayStatus.Listeners.
	//
	// +optional
	Address *ListenerAddress
	// Ports is a list of ports associated with the Address.
	Ports []ListenerPort
}

const (
	IPAddressType    = "IPAddress"
	NamedAddressType = "NamedAddress"
)

// ListenerAddress describes an address bound by the
// Listener.
type ListenerAddress struct {
	// Type of the Address. This is one of the *AddressType constants.
	Type string
	// Address value. Examples: "1.2.3.4", "128::1", "my-ip-address".
	Address string
}

type ListenerPort struct {
	Port      *int
	Protocols []string
	TLS       *ListenerTLS
	Extension *core.TypedLocalObjectReference
}

const (
	TLSVersion_1_0 = "TLS_1_0"
	TLSVersion_1_1 = "TLS_1_1"
	TLSVersion_1_2 = "TLS_1_2"
	TLSVersion_1_3 = "TLS_1_3"
)

// ListenerTLS describes the TLS configuration for a given port.
//
// References
// - nginx: https://nginx.org/en/docs/http/configuring_https_servers.html
// - envoy: https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/auth/cert.proto
// - haproxy: https://www.haproxy.com/documentation/aloha/9-5/traffic-management/lb-layer7/tls/
// - gcp: https://cloud.google.com/load-balancing/docs/use-ssl-policies#creating_an_ssl_policy_with_a_custom_profile
// - aws: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html#describe-ssl-policies
// - azure: https://docs.microsoft.com/en-us/azure/app-service/configure-ssl-bindings#enforce-tls-1112
type ListenerTLS struct {
	Certificates []core.TypedLocalObjectReference

	// MinimumVersion of TLS allowed. It is recommended to use one of
	// the TLSVersion_* constants above. Note: this is not strongly
	// typed to allow newly available version to be used without
	// requiring updates to the API types. String must be of the form
	// "<protocol>_<major>_<minor>".
	//
	// Support: Core.
	MinimumVersion *string
	// Options are a list of key/value pairs to give extended options
	// to the provider.
	//
	// Support: Implementation-specific.
	//
	// There variation among providers as to how ciphersuites are
	// expressed. If there is a common subset for expressing ciphers
	// then it will make sense to loft that as a core API
	// construct.
	Options map[string]string
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	// Listeners are xxx
	Listeners []ListenerStatus
	// Routes are xxx
	Routes []GatewayRouteStatus
}

type ListenerStatus struct {
	// Errors is a list of reasons why a given Listener Spec is
	// not valid. Errors will be empty if the Spec is valid.
	Errors []ListenerError

	//
	Address string
}

type ListenerErrorCode string

const (
	// ErrListenerInvalidSpec is a generic error that is a
	// catch all for unsupported configurations that do not match a
	// more specific error. Implementors should try to use more
	// specific errors instead of this one to give users and
	// automation a more information.
	ErrListenerInvalidSpec ListenerErrorCode = "InvalidSpec"
	// ErrListenerBadAddress indicates the Address
	ErrListenerBadAddress ListenerErrorCode = "InvalidAddress"
)

type ListenerError struct {
	Code    ListenerErrorCode
	Message string
}

type GatewayRouteStatus struct {
}

// +kubebuilder:object:root=true

// Gateway is the Schema for the gateways API
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
