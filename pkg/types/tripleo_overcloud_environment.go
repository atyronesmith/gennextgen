package types

import "github.com/atyronesmith/gennextgen/pkg/utils"

type TripleoOvercloudEnvironment struct {
	ParmaterDefaults TripleoOvercloudEnvironmentParameterDefaults `yaml:"parameter_defaults"`
}

type TripleoOvercloudEnvironmentParameterDefaults struct {
	CtlplaneNetworkAttributes  CtlplaneNetworkAttributes  `yaml:"CtlplaneNetworkAttributes"`
	DeployedNetworkEnvironment DeployedNetworkEnvironment `yaml:"DeployedNetworkEnvironment"`
	ControlPlaneVipData        ControlPlaneVipData        `yaml:"ControlPlaneVipData"`
}

type CtlplaneNetworkAttributes struct {
	Network CtlplaneNetwork           `yaml:"network"`
	Subnets map[string]CtlplaneSubnet `yaml:"subnets"`
}

type CtlplaneNetwork struct {
	DNSDomain string   `yaml:"dns_domain"`
	MTU       int64    `yaml:"mtu"`
	Name      string   `yaml:"name"`
	Tags      []string `yaml:"tags"`
}

type CtlplaneSubnet struct {
	CIDR           string              `yaml:"cidr"`
	DNSNameservers []interface{}       `yaml:"dns_nameservers"`
	GatewayIP      string              `yaml:"gateway_ip"`
	HostRoutes     []CtlplaneHostRoute `yaml:"host_routes"`
	IPVersion      int64               `yaml:"ip_version"`
	Name           string              `yaml:"name"`
}

type CtlplaneHostRoute struct {
	Destination string `yaml:"destination"`
	Nexthop     string `yaml:"nexthop"`
}

type DeployedNetworkEnvironment struct {
	NetAttributesMap map[string]DeployedNetworkEnvironmentNetAttributes `yaml:"net_attributes_map"`
	NetCIDRMap       map[string][]string                                `yaml:"net_cidr_map"`
	NetIPVersionMap  map[string]int64                                   `yaml:"net_ip_version_map"`
}

type DeployedNetworkEnvironmentNetAttributes struct {
	Network DeployedNetworkEnvironmentNetwork           `yaml:"network"`
	Subnets map[string]DeployedNetworkEnvironmentSubnet `yaml:"subnets"`
}

type DeployedNetworkEnvironmentNetwork struct {
	DNSDomain string   `yaml:"dns_domain"`
	MTU       int64    `yaml:"mtu"`
	Name      string   `yaml:"name"`
	Tags      []string `yaml:"tags"`
}

type DeployedNetworkEnvironmentSubnet struct {
	CIDR           string                                `yaml:"cidr"`
	DNSNameservers []interface{}                         `yaml:"dns_nameservers"`
	GatewayIP      interface{}                           `yaml:"gateway_ip"`
	HostRoutes     []DeployedNetworkEnvironmentHostRoute `yaml:"host_routes"`
	IPVersion      int64                                 `yaml:"ip_version"`
	Name           string                                `yaml:"name"`
	Tags           []string                              `yaml:"tags"`
}

type DeployedNetworkEnvironmentHostRoute struct {
	Destination string `yaml:"destination"`
	Nexthop     string `yaml:"nexthop"`
}

type ControlPlaneVipData struct {
	FixedIPS []ControlPlaneVipDataFixedIP `json:"fixed_ips"`
	Name     string                       `json:"name"`
	Network  ControlPlaneVipDataNetwork   `json:"network"`
	Subnets  []ControlPlaneVipDataSubnet  `json:"subnets"`
}

type ControlPlaneVipDataFixedIP struct {
	IPAddress string `json:"ip_address"`
}

type ControlPlaneVipDataNetwork struct {
	Tags []string `json:"tags"`
}

type ControlPlaneVipDataSubnet struct {
	IPVersion int64 `json:"ip_version"`
}

var tripleoOvercloudEnvironment *TripleoOvercloudEnvironment

func GetTripleoOvercloudEnvironment() (*TripleoOvercloudEnvironment, error) {
	if tripleoOvercloudEnvironment == nil {
		tripleoOvercloudEnvironment = &TripleoOvercloudEnvironment{}
		err := utils.YamlToStruct(utils.GetFullPath(utils.TRIPLEO_OVERCLOUD_ENVIRONMENT), &tripleoOvercloudEnvironment)
		if err != nil {
			return nil, err
		}
	}
	return tripleoOvercloudEnvironment, nil
}
