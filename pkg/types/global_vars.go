package types

import "github.com/atyronesmith/gennextgen/pkg/utils"

type GlobalVars struct {
	AllNodesExtraMapData    map[string]interface{}       `yaml:"all_nodes_extra_map_data"`
	CloudDomain             string                       `yaml:"cloud_domain"`
	CloudNames              map[string]string            `yaml:"cloud_names"`
	ContainerCLI            string                       `yaml:"container_cli"`
	ControlVirtualIP        string                       `yaml:"control_virtual_ip"`
	DeployArtifactFiles     []interface{}                `yaml:"deploy_artifact_files"`
	DeployArtifactUrls      []interface{}                `yaml:"deploy_artifact_urls"`
	DeployIdentifier        string                       `yaml:"deploy_identifier"`
	DeployStepsMax          int64                        `yaml:"deploy_steps_max"`
	EnableInternalTLS       bool                         `yaml:"enable_internal_tls"`
	EnabledNetworks         []string                     `yaml:"enabled_networks"`
	EnabledServices         []string                     `yaml:"enabled_services"`
	ExtraHostsEntries       []string                     `yaml:"extra_hosts_entries"`
	HideSensitiveLogs       bool                         `yaml:"hide_sensitive_logs"`
	HostsEntry              []string                     `yaml:"hosts_entry"`
	KeystoneResources       map[string]KeystoneResources `yaml:"keystone_resources"`
	NetCIDRMap              map[string][]string          `yaml:"net_cidr_map"`
	NetVipMap               map[string]string            `yaml:"net_vip_map"`
	NetworkSafeDefaults     bool                         `yaml:"network_safe_defaults"`
	NetworkVirtualIPS       map[string]NetworkVirtualIP  `yaml:"network_virtual_ips"`
	Networks                map[string]Networks          `yaml:"networks"`
	NovaAdditionalCell      bool                         `yaml:"nova_additional_cell"`
	PingTestGatewayIPS      map[string][]string          `yaml:"ping_test_gateway_ips"`
	PingTestIPS             map[string]string            `yaml:"ping_test_ips"`
	PrimaryRoleName         string                       `yaml:"primary_role_name"`
	ServiceNetMap           map[string]string            `yaml:"service_net_map"`
	UndercloudHostsEntries  []string                     `yaml:"undercloud_hosts_entries"`
	ValidateControllersICMP bool                         `yaml:"validate_controllers_icmp"`
	ValidateFQDN            bool                         `yaml:"validate_fqdn"`
	ValidateGatewaysICMP    bool                         `yaml:"validate_gateways_icmp"`
	VipHostsEntries         []string                     `yaml:"vip_hosts_entries"`
}

type NetworkVirtualIP struct {
	Index     int64  `yaml:"index"`
	IpAddress string `yaml:"ip_address"`
}

type KeystoneResources struct {
	Domains   []string               `yaml:"domains"`
	Endpoints Endpoints              `yaml:"endpoints"`
	Project   string                 `yaml:"project"`
	Region    string                 `yaml:"region"`
	Roles     []string               `yaml:"roles"`
	Service   string                 `yaml:"service"`
	Users     map[string]interface{} `yaml:"users"`
}

type Endpoints struct {
	Admin    string `yaml:"admin"`
	Internal string `yaml:"internal"`
	Public   string `yaml:"public"`
}

type Networks struct {
	Name      string `yaml:"name"`
	NameLower string `yaml:"name_lower"`
}

var globalVars *GlobalVars

func GetGlobalVars() (*GlobalVars, error) {
	if globalVars == nil {
		globalVars = &GlobalVars{}
		err := utils.YamlToStruct(utils.GetFullPath(utils.GLOBAL_VARS), &globalVars)
		if err != nil {
			return nil, err
		}
	}
	return globalVars, nil

}
