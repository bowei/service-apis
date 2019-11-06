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

// GatewaySpec defines the desired state of Gateway
type GatewaySpec struct {
	// Class used for this Gateway. This is the name of a GatewayClass.
	Class string `json:"class"`
	// Listeners associated with this Gateway. Listeners define what addresses,
	// ports, protocols, configuration are bound on this Gateway.
	Listeners []GatewayListener
	// Routes associated with this Gateway.
	Routes []core.TypedLocalObjectReference
}

// GatewayListener is xxx
type GatewayListener struct {
	// Address associated with the GatewayListener. This is either an IP
	// address in standard form i.e. IPv4 "1.2.3.4", IPv6
	Address *string
	Ports   []GatewayListenerPort
}

type GatewayListenerPort struct {
	Port      *int
	Protocols []string
	TLS       *GatewayListenerTLS
	Extension *core.TypedLocalObjectReference
}

const (
	TLS_1_0 = "TLS_1_0"
	TLS_1_1 = "TLS_1_1"
	TLS_1_2 = "TLS_1_2"
	TLS_1_3 = "TLS_1_3"
)

// GatewayListenerTLS describes TLS configuration for a given port.
//
// | References TLS support:
// | - nginx: https://nginx.org/en/docs/http/configuring_https_servers.html
// | - envoy: https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/auth/cert.proto
// | - haproxy: https://www.haproxy.com/documentation/aloha/9-5/traffic-management/lb-layer7/tls/
// | - gcp: https://cloud.google.com/load-balancing/docs/use-ssl-policies#creating_an_ssl_policy_with_a_custom_profile
// | - aws: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html#describe-ssl-policies
// | - azure: https://docs.microsoft.com/en-us/azure/app-service/configure-ssl-bindings#enforce-tls-1112
type GatewayListenerTLS struct {
	Certificates []core.TypedLocalObjectReference

	// MinimumVersion of TLS allowed. One of the TLS_* constants above.
	// Note: this is not strongly typed to avoid updating API constants to
	// match available versions. Version constants MUST be of the form
	// TLS_<major>_<minor>.
	//
	// Support: Core.
	MinimumVersion *string
	// Options are a list of key/value pairs to give extended options to the provider.
	//
	// Support: Implementation-specific.
	//
	// | There variation among providers as to how ciphersuites are
	// | expressed. If there is a common subset for expression ciphers then
	// | it will make sense to loft that as a core API construct. We leave
	// | this open ended for now to allow users access to TLS-features which
	// | is an important use case but may not prove to be portable.
	Options map[string]string
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	// Listeners are xxx
	Listeners []GatewayListenerStatus
	// Routes are xxx
	Routes []GatewayRouteStatus
}

type GatewayListenerStatus struct {
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
