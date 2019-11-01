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
