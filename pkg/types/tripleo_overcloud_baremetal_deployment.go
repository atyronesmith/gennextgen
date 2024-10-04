package types

import "github.com/atyronesmith/gennextgen/pkg/utils"

type TripleoOvercloudBaremetalDeployment []struct {
	AnsiblePlaybooks []TOBDAnsiblePlaybooks `yaml:"ansible_playbooks"`
	Count            int                    `yaml:"count"`
	Defaults         TOBDDefaults           `yaml:"defaults"`
	Instances        []TOBDInstances        `yaml:"instances"`
	Name             string                 `yaml:"name"`
}
type TOBDRoleGrowvolsArgs struct {
	Default string `yaml:"default"`
}

type TOBDExtraVars struct {
	RoleGrowvolsArgs TOBDRoleGrowvolsArgs `yaml:"role_growvols_args"`
}
type TOBDAnsiblePlaybooks struct {
	ExtraVars map[string]interface{} `yaml:"extra_vars"`
	Playbook  string                 `yaml:"playbook"`
}
type TOBDNetworkConfig struct {
	DefaultRouteNetwork []string `yaml:"default_route_network"`
	Template            string   `yaml:"template"`
}
type TOBDDefaults struct {
	NetworkConfig TOBDNetworkConfig `yaml:"network_config"`
	ResourceClass string            `yaml:"resource_class"`
}
type TOBDNetworks struct {
	Network string `yaml:"network"`
	Vif     bool   `yaml:"vif,omitempty"`
	FixedIP string `yaml:"fixed_ip,omitempty"`
	Subnet  string `yaml:"subnet,omitempty"`
}
type TOBDInstances struct {
	Hostname string         `yaml:"hostname"`
	Name     string         `yaml:"name"`
	Networks []TOBDNetworks `yaml:"networks"`
}

var tripleoOvercloudBaremetalDeployment TripleoOvercloudBaremetalDeployment

func GetTripleoOvercloudBaremetalDeployment(path string) (*TripleoOvercloudBaremetalDeployment, error) {
	if path == "" {
		path = utils.GetFullPath(utils.BAREMETAL_DEPLOY)
	}

	if len(tripleoOvercloudBaremetalDeployment) == 0 {
		tripleoOvercloudBaremetalDeployment = TripleoOvercloudBaremetalDeployment{}
		err := utils.YamlToStruct(path, &tripleoOvercloudBaremetalDeployment)
		if err != nil {
			return nil, err
		}
	}
	return &tripleoOvercloudBaremetalDeployment, nil

}
