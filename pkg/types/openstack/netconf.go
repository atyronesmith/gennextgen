package openstack

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetNameStr is used for validation of a net name.
type NetNameStr string

// Network definition
type Network struct {
	// Name of the network, e.g. External, InternalApi, ...
	Name NetNameStr `json:"name"`

	// DNSDomain name of the Network
	DNSDomain string `json:"dnsDomain"`

	// MTU of the network
	MTU int `json:"mtu"`

	// Subnets of the tripleo network
	Subnets []Subnet `json:"subnets"`
}

func NewNetwork() *Network {
	return &Network{
		MTU: 1500,
	}
}

// Subnet definition
type Subnet struct {
	// Name of the subnet
	Name NetNameStr `json:"name"`

	// Cidr the cidr to use for this network
	Cidr string `json:"cidr"`

	// DNSDomain name of the subnet, allows to overwrite the DNSDomain of the Network
	DNSDomain *string `json:"dnsDomain,omitempty"`

	// Vlan ID
	Vlan *int `json:"vlan,omitempty"`

	// AllocationRanges a list of AllocationRange for assignment. Allocation will start
	// from first range, first address.
	AllocationRanges []AllocationRange `json:"allocationRanges"`

	// ExcludeAddresses a set of IPs that should be excluded from used as reservation, for both dynamic
	// and static via IPSet FixedIP parameter
	ExcludeAddresses []string `json:"excludeAddresses,omitempty"`

	// Gateway optional gateway for the network
	Gateway string `json:"gateway,omitempty"`

	// Routes, list of networks that should be routed via network gateway.
	Routes []Route `json:"routes,omitempty"`
}

// AllocationRange definition
type AllocationRange struct {
	// Start IP for the AllocationRange
	Start string `json:"start"`

	// End IP for the AllocationRange
	End string `json:"end"`
}

// Route definition
type Route struct {
	// Destination, network CIDR
	Destination string `json:"destination"`

	// Nexthop, gateway for the destination
	Nexthop string `json:"nexthop"`
}

// NetConfigSpec defines the desired state of NetConfig
type NetConfigSpec struct {
	// Networks, list of all tripleo networks of the deployment
	Networks []Network `json:"networks"`
}

// NetConfig is the Schema for the netconfigs API
type NetConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetConfigSpec `json:"spec,omitempty"`
}

// NetConfigList contains a list of NetConfig
type NetConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetConfig `json:"items"`
}

func NewNetConfig() *NetConfig {
	return &NetConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "network.openstack.org/v1beta1",
			Kind:       "NetConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "openstacknetconfig",
		},
	}
}
