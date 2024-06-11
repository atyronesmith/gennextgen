package types

import "github.com/atyronesmith/gennextgen/pkg/utils"

type TripleoOvercloudEnvironment struct {
	ParmaterDefaults TOEParameterDefaults `yaml:"parameter_defaults"`
}

type TOEParameterDefaults struct {
	CtlplaneNetworkAttributes  TOECtlplaneNetworkAttributes  `yaml:"CtlplaneNetworkAttributes"`
	DeployedNetworkEnvironment TOEDeployedNetworkEnvironment `yaml:"DeployedNetworkEnvironment"`
	ControlPlaneVipData        TOEControlPlaneVipData        `yaml:"ControlPlaneVipData"`
}

type TOECtlplaneNetworkAttributes struct {
	Network TOECtlplaneNetwork           `yaml:"network"`
	Subnets map[string]TOECtlplaneSubnet `yaml:"subnets"`
}

type TOECtlplaneNetwork struct {
	DNSDomain string   `yaml:"dns_domain"`
	MTU       int64    `yaml:"mtu"`
	Name      string   `yaml:"name"`
	Tags      []string `yaml:"tags"`
}

type TOECtlplaneSubnet struct {
	CIDR           string                 `yaml:"cidr"`
	DNSNameservers []interface{}          `yaml:"dns_nameservers"`
	GatewayIP      string                 `yaml:"gateway_ip"`
	HostRoutes     []TOECtlplaneHostRoute `yaml:"host_routes"`
	IPVersion      int64                  `yaml:"ip_version"`
	Name           string                 `yaml:"name"`
}

type TOECtlplaneHostRoute struct {
	Destination string `yaml:"destination"`
	Nexthop     string `yaml:"nexthop"`
}

type TOEDeployedNetworkEnvironment struct {
	NetAttributesMap map[string]TOEDeployedNetworkEnvironmentNetAttributes `yaml:"net_attributes_map"`
	NetCIDRMap       map[string][]string                                   `yaml:"net_cidr_map"`
	NetIPVersionMap  map[string]int64                                      `yaml:"net_ip_version_map"`
}

type TOEDeployedNetworkEnvironmentNetAttributes struct {
	Network TOEDeployedNetworkEnvironmentNetwork           `yaml:"network"`
	Subnets map[string]TOEDeployedNetworkEnvironmentSubnet `yaml:"subnets"`
}

type TOEDeployedNetworkEnvironmentNetwork struct {
	DNSDomain string   `yaml:"dns_domain"`
	MTU       int64    `yaml:"mtu"`
	Name      string   `yaml:"name"`
	Tags      []string `yaml:"tags"`
}

type TOEDeployedNetworkEnvironmentSubnet struct {
	CIDR           string                                   `yaml:"cidr"`
	DNSNameservers []interface{}                            `yaml:"dns_nameservers"`
	GatewayIP      interface{}                              `yaml:"gateway_ip"`
	HostRoutes     []TOEDeployedNetworkEnvironmentHostRoute `yaml:"host_routes"`
	IPVersion      int64                                    `yaml:"ip_version"`
	Name           string                                   `yaml:"name"`
	Tags           []string                                 `yaml:"tags"`
}

type TOEDeployedNetworkEnvironmentHostRoute struct {
	Destination string `yaml:"destination"`
	Nexthop     string `yaml:"nexthop"`
}

type TOEControlPlaneVipData struct {
	FixedIPS []TOEControlPlaneVipDataFixedIP `json:"fixed_ips"`
	Name     string                          `json:"name"`
	Network  TOEControlPlaneVipDataNetwork   `json:"network"`
	Subnets  []TOEControlPlaneVipDataSubnet  `json:"subnets"`
}

type TOEControlPlaneVipDataFixedIP struct {
	IPAddress string `json:"ip_address"`
}

type TOEControlPlaneVipDataNetwork struct {
	Tags []string `json:"tags"`
}

type TOEControlPlaneVipDataSubnet struct {
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
