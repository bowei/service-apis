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

// +kubebuilder:object:root=true

// Gateway represents an instantiation of a service-traffic handling infrastructure.
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

// GatewaySpec defines the desired state of Gateway.
//
// The Spec is split into two major pieces: listeners describing
// client-facing properties and routes that describe application-level
// routing.
//
// Not all possible combinations of options specified in the Spec are
// valid. Some invalid configurations can be caught synchronously via a
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
	Listeners []Listener `json:"listeners"`
	// Routes associated with this Gateway. Routes define
	// protocol-specific routing to backends (e.g. Services).
	Routes []core.TypedLocalObjectReference `json:"routes"`
}

// Listener defines a
type Listener struct {
	// Address bound on the listener. This is optional and behavior
	// can depend on GatewayClass. If a value is set in the spec and
	// the request address is invalid, the GatewayClass MUST indicate
	// this in the associated entry in GatewayStatus.Listeners.
	//
	// +optional
	Address *ListenerAddress `json:"address"`
	// Ports is a list of ports associated with the Address.
	Ports []ListenerPort `json:"ports"`
}

const (
	IPAddressType    = "IPAddress"
	NamedAddressType = "NamedAddress"
)

// ListenerAddress describes an address bound by the
// Listener.
type ListenerAddress struct {
	// Type of the Address. This is one of the *AddressType constants.
	//
	// Support: Extended
	Type string `json:"type"`
	// Address value. Examples: "1.2.3.4", "128::1", "my-ip-address".
	Address string `json:"address"`
}

// ListenerPort xxx
type ListenerPort struct {
	Port      *int                            `json:"port"`
	Protocols []string                        `json:"protocols"`
	TLS       *ListenerTLS                    `json:"tls"`
	Extension *core.TypedLocalObjectReference `json:"extension"`
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
	// Certificates is a list of certificates containing resources
	// that are bound to the listener.
	//
	// If apiGroup and kind are empty, will default to Kubernetes Secrets resources.
	//
	// Support: Core (Kubernetes Secrets)
	// Support: Implementation-specific (Other resource types)
	Certificates []core.TypedLocalObjectReference `json:"certificates,omitempty"`
	// MinimumVersion of TLS allowed. It is recommended to use one of
	// the TLSVersion_* constants above. Note: this is not strongly
	// typed to allow newly available version to be used without
	// requiring updates to the API types. String must be of the form
	// "<protocol>_<major>_<minor>".
	//
	// Support: Core
	MinimumVersion *string `json:"minimumVersion"`
	// Options are a list of key/value pairs to give extended options
	// to the provider.
	//
	// There variation among providers as to how ciphersuites are
	// expressed. If there is a common subset for expressing ciphers
	// then it will make sense to loft that as a core API
	// construct.
	//
	// Support: Implementation-specific.
	Options map[string]string `json:"options"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	// XXX overall status

	// Listeners is the status for each listener block in the
	// Spec. The status for a given block will match the order as
	// declared in the Spec, e.g. the status for Spec.Listeners[3]
	// will be in Status.Listeners[3].
	Listeners []ListenerStatus `json:"listeners"`
	// Routes is the status for each attached route to the
	// Gateway. The status for a given route will match the orer as
	// declared in the Spec, e.g. the status for Spec.Routes[3] will
	// be in Status.Routes[3].
	Routes []GatewayRouteStatus `json:"routes"`
}

type ListenerStatus struct {
	// Errors is a list of reasons why a given Listener Spec is
	// not valid. Errors will be empty if the Spec is valid.
	Errors []ListenerError `json:"errors"`

	//
	Address string `json:"address"`
}

type ListenerErrorReason string

const (
	// ErrListenerInvalidSpec is a generic error that is a
	// catch all for unsupported configurations that do not match a
	// more specific error. Implementors should try to use more
	// specific errors instead of this one to give users and
	// automation a more information.
	ErrListenerInvalidSpec ListenerErrorReason = "InvalidSpec"
	// ErrListenerBadAddress indicates the Address
	ErrListenerBadAddress ListenerErrorReason = "InvalidAddress"
)

// ListenerError is an error status for a given ListenerSpec.
type ListenerError struct {
	// Reason is a automation friendly reason code for the error.
	Reason ListenerErrorReason `json:"reason"`
	// Message is a human-understandable error message.
	Message string `json:"message"`
}

type GatewayRouteStatus struct {
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
