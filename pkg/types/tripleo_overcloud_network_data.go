package types

import "github.com/atyronesmith/gennextgen/pkg/utils"

type TripleoOvercloudNetworkData []TripleoOvercloudNetworkDataEntry

type TripleoOvercloudNetworkDataEntry struct {
	Name      string                `yaml:"name"`
	NameLower string                `yaml:"name_lower"`
	Vip       bool                  `yaml:"vip"`
	DNSDomain string                `yaml:"dns_domain,omitempty"`
	MTU       int                   `yaml:"mtu"`
	Ipv6      bool                  `yaml:"ipv6"`
	Subnets   map[string]TONDSubnet `yaml:"subnets"`
	Enabled   bool                  `yaml:"enabled"`
}

type TONDSubnet struct {
	Enabled             bool                 `yaml:"enabled"`
	VLAN                int64                `yaml:"vlan"`
	AllocationPools     []TONDAllocationPool `yaml:"allocation_pools"`
	GatewayIp           string               `yaml:"gateway_ip"`
	GatewayIpV6         string               `yaml:"gateway_ipv6"`
	Routes              []TONDRoutesIpv6     `yaml:"routes,omitempty"`
	IpSubnet            string               `yaml:"ip_subnet"`
	Ipv6Subnet          string               `yaml:"ipv6_subnet"`
	Ipv6AllocationPools []TONDAllocationPool `yaml:"ipv6_allocation_pools"`
	RoutesIpv6          []TONDRoutesIpv6     `yaml:"routes_ipv6,omitempty"`
}

type TONDAllocationPool struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type TONDRoutesIpv6 struct {
	Default     bool   `yaml:"default"`
	Destination string `yaml:"destination"`
	NextHop     string `yaml:"nexthop"`
}

type TONDInternalAPISubnet1 struct {
	Ipv6Subnet          string        `yaml:"ipv6_subnet"`
	Ipv6AllocationPools []interface{} `yaml:"ipv6_allocation_pools"`
}

var tripleoOvercloudNetworkData TripleoOvercloudNetworkData

func GetTripleoOvercloudNetworkData() (*TripleoOvercloudNetworkData, error) {
	if len(tripleoOvercloudNetworkData) == 0 {
		tripleoOvercloudNetworkData = TripleoOvercloudNetworkData{}
		err := utils.YamlToStruct(utils.GetFullPath(utils.TRIPLEO_OVERCLOUD_NETWORK_DATA), &tripleoOvercloudNetworkData)
		if err != nil {
			return nil, err
		}
	}
	return &tripleoOvercloudNetworkData, nil

}
