package types

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/atyronesmith/gennextgen/pkg/utils"
	v3 "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
)

type ConfigDownload struct {
	Networks        map[string]*OSPNetwork
	EnabledServices map[string]*OSPService
	Hosts           map[string]*TripleoHost
	Roles           []TripleoRole
	PrimaryRoleName string
	Passwords       map[string]string
}

type PasswordMapping struct {
	File    string `yaml:"file"`
	Address string `yaml:"address"`
	Role    string `yaml:"role"`
}

type OSPServicePasswords struct {
	Password   string
	DbPassword string
}

type OSPNetwork struct {
	Name       string              `yaml:"name"`
	IsRoleNet  bool                `yaml:"is_role_net"`
	Domain     string              `yaml:"domain"`
	IP         string              `yaml:"ip"`
	GatewayIp  string              `yaml:"gateway_ip"`
	Cidr       string              `yaml:"cidr"`
	Vip        string              `yaml:"vip"`
	Mtu        int                 `yaml:"mtu"`
	VlanId     int                 `yaml:"vlan_id"`
	HostRoutes []TripleoHostRoutes `yaml:"host_routes"`
	Subnets    []OSPNetworkSubnet  `yaml:"subnets,omitempty"`
}

func (n *OSPNetwork) UnmarshalYaml(node *v3.Node) error {
	fmt.Printf("UnmarshalYaml: %+v\n", node)
	return nil
}

type OSPNetworkSubnet struct {
	Name            string
	Subnet          string
	VlanId          int
	AllocationPools []OSPNetworkSubnetPool
}

type OSPNetworkSubnetPool struct {
	Start string
	End   string
}

type TripleoHost struct {
	Name                string
	CanonicalName       string
	AnsibleHost         string
	DefaultRouteNetwork string
	Networks            map[string]*TripleoHostNetwork
	Vars                map[string]interface{}
}

type TripleoHostNetwork struct {
	Name     string
	IP       net.IP
	Hostname string
}

type TripleoRole struct {
	Name     string
	Hosts    []TripleoHost
	Networks map[string]*OSPNetwork
	Vars     map[string]interface{}
}

type TripleoHostRoutes struct {
	Default     bool   `yaml:"default"`
	Destination string `yaml:"destination"`
	NextHop     string `yaml:"nexthop"`
}
type OSPService struct {
	Name       string
	Network    string
	Password   string
	DbPassword string
}

func NewConfigDownload() *ConfigDownload {
	return &ConfigDownload{
		Networks:        make(map[string]*OSPNetwork),
		Hosts:           make(map[string]*TripleoHost),
		Passwords:       make(map[string]string),
		EnabledServices: make(map[string]*OSPService),
	}
}

type BaremetalDeploy []struct {
	AnsiblePlaybooks []struct {
		ExtraVars struct {
			RoleGrowvolsArgs struct {
				Default string `yaml:"default"`
			} `yaml:"role_growvols_args"`
		} `yaml:"extra_vars"`
		Playbook string `yaml:"playbook"`
	} `yaml:"ansible_playbooks"`
	Count    int `yaml:"count"`
	Defaults struct {
		NetworkConfig struct {
			DefaultRouteNetwork []string `yaml:"default_route_network"`
			Template            string   `yaml:"template"`
		} `yaml:"network_config"`
		ResourceClass string `yaml:"resource_class"`
	} `yaml:"defaults"`
	Instances []struct {
		Hostname string `yaml:"hostname"`
		Name     string `yaml:"name"`
		Networks []struct {
			Network string `yaml:"network"`
			Vif     bool   `yaml:"vif,omitempty"`
			FixedIP string `yaml:"fixed_ip,omitempty"`
			Subnet  string `yaml:"subnet,omitempty"`
		} `yaml:"networks"`
	} `yaml:"instances"`
	Name string `yaml:"name"`
}

func (cdl *ConfigDownload) Process(mappingYaml string) error {
	//	var nncp types.NNCP
	cfg := utils.GetConfig()

	gv, err := utils.YamlToMap(utils.GetFullPath(utils.GLOBAL_VARS))
	if err != nil {
		return err
	}

	err = cdl.ProcessGlobalVars(gv)
	if err != nil {
		return err
	}

	tai, err := utils.YamlToMap(utils.GetFullPath(utils.TRIPLEO_ANSIBLE_INVENTORY_YAML))
	if err != nil {
		return err
	}

	err = cdl.ProcessTripleoAnsibleInventory(tai)
	if err != nil {
		return err
	}

	primaryRole := cdl.PrimaryRoleName
	if cfg.ControllerRole != "" {
		primaryRole = cfg.ControllerRole
	}

	groupVars, err := utils.YamlToMap(utils.GetFullPath("config-download/overcloud/group_vars/" + primaryRole))
	if err != nil {
		fmt.Printf("Error while reading group_vars: %v\n", err)
		return err
	}

	err = cdl.ProcessGroupVars(groupVars)
	if err != nil {
		return err
	}

	err = cdl.ProcessPasswords(mappingYaml, groupVars)
	if err != nil {
		return err
	}

	return nil
}

// Must be run first
func (cdl *ConfigDownload) ProcessGlobalVars(data map[string]interface{}) error {
	// get the global list of infra networks
	if ncm, ok := data["net_cidr_map"]; ok {
		for k, v := range ncm.(map[string]interface{}) {
			for _, cidr := range v.([]interface{}) {
				cdl.Networks[k] = &OSPNetwork{Name: k, Cidr: cidr.(string), IsRoleNet: false}
			}
		}
	} else {
		return fmt.Errorf("Missing net_cidr_map in global_vars.yaml.")
	}

	// Fill in Domain info for the networks
	if cns, ok := data["cloud_names"]; ok {
		for k, v := range cns.(map[string]interface{}) {
			netName := strings.Replace(k, "cloud_name_", "", -1)
			if n, ok := cdl.Networks[netName]; ok {
				n.Domain = v.(string)
			}
		}
	} else {
		return fmt.Errorf("Missing cloud_names in global_vars.yaml.")
	}

	// Fill in vip info for the networks
	if nvm, ok := data["net_vip_map"]; ok {
		for _, v := range cdl.Networks {
			if vip, ok := nvm.(map[string]interface{})[v.Name]; ok {
				v.Vip = vip.(string)
			}
		}
	} else {
		return fmt.Errorf("Missing net_vip_map in global_vars.yaml.")
	}

	if es, ok := data["enabled_services"]; ok {
		for _, v := range es.([]interface{}) {
			cdl.EnabledServices[v.(string)] = &OSPService{Name: v.(string)}
		}
	} else {
		return fmt.Errorf("Missing enabled_services in global_vars.yaml.")
	}

	if snm, ok := data["service_net_map"]; ok {
		for k, v := range snm.(map[string]interface{}) {
			if s, ok := cdl.EnabledServices[k]; ok {
				s.Network = v.(string)
			}
		}
	} else {
		return fmt.Errorf("Missing service_net_map in global_vars.yaml.")
	}

	if es, ok := data["networks"]; ok {
		for _, v := range es.(map[string]interface{}) {
			if n, ok := cdl.Networks[v.(map[string]interface{})["name_lower"].(string)]; ok {
				n.IsRoleNet = true
			}
		}
	} else {
		return fmt.Errorf("Missing enabled_services in global_vars.yaml.")
	}

	if prn, ok := data["primary_role_name"]; ok {
		cdl.PrimaryRoleName = prn.(string)
	} else {
		return fmt.Errorf("Missing primary_role_name in global_vars.yaml.")
	}

	return nil
}

func (cdl *ConfigDownload) ProcessGroupVars(data map[string]interface{}) error {
	for index := range cdl.Networks {
		net := cdl.Networks[index]
		if gw, ok := data[index+"_gateway_ip"]; ok {
			if gw != nil {
				net.GatewayIp = gw.(string)
			}
		}
		// if hr, ok := data[index+"_host_routes"].([]interface{}); ok {
		// 	net.HostRoutes = append(v.HostRoutes, hr.([]string)...)
		// }
		if m, ok := data[index+"_mtu"]; ok {
			net.Mtu = m.(int)
		} else {
			return fmt.Errorf("Missing mtu for %s\n", index)
		}
		if vi, ok := data[index+"_vlan_id"]; ok {
			net.VlanId = vi.(int)
		}
	}

	return nil
}

func mapStringGetter(v interface{}, key string, valPtr interface{}) {
	if value, ok := v.(map[string]interface{})[key]; ok {
		switch p := valPtr.(type) {
		case *int:
			if intVal, ok := value.(int); ok {
				*p = intVal
			} else {
				*p, _ = strconv.Atoi(value.(string))
			}
		case *string:
			if value != nil {
				*p = value.(string)
			}
		}
	}
}

func (cdl *ConfigDownload) ProcessTripleoAnsibleInventory(data map[string]interface{}) error {
	for topKey, topValue := range data {
		if topKey == "Undercloud" {
			continue
		}

		var role TripleoRole

		role.Networks = make(map[string]*OSPNetwork)

		// Get the role networks
		if v, ok := topValue.(map[string]interface{})["vars"]; ok {
			if trn, ok := v.(map[string]interface{})["tripleo_role_networks"]; ok {
				// Get the list of names of networks in this role
				for _, net := range trn.([]interface{}) {
					role.Networks[net.(string)] = &OSPNetwork{
						Name: net.(string),
					}
				}
				// Look for network specific vars
				for netName := range role.Networks {
					osNet := role.Networks[netName]

					mapStringGetter(v, netName+"_cidr", &osNet.Cidr)
					mapStringGetter(v, netName+"_gateway_ip", &osNet.GatewayIp)
					mapStringGetter(v, netName+"_vip", &osNet.Vip)
					mapStringGetter(v, netName+"_mtu", &osNet.Mtu)
					mapStringGetter(v, netName+"_vlan_id", &osNet.VlanId)

					if hostRoutes, ok := v.(map[string]interface{})[netName+"_host_routes"]; ok {
						var hrs []TripleoHostRoutes
						utils.MarshalArray(hostRoutes, &hrs)
						osNet.HostRoutes = append(osNet.HostRoutes, hrs...)
					}
				}
				if _, ok := v.(map[string]interface{})["tripleo_role_name"]; ok {
					mapStringGetter(v, "tripleo_role_name", &role.Name)
				} else {
					return fmt.Errorf("Missing tripleo_role_name in %s\n", topKey)
				}
				role.Vars = v.(map[string]interface{})
			} else {
				return fmt.Errorf("Missing tripleo_role_networks in %s\n", topKey)
			}
		}

		if hosts, ok := topValue.(map[string]interface{})["hosts"]; ok {
			for hostName, hostVarMap := range hosts.(map[string]interface{}) {
				th := TripleoHost{}
				th.Networks = make(map[string]*TripleoHostNetwork)

				mapStringGetter(hostVarMap, "canonical_hostname", &th.CanonicalName)
				mapStringGetter(hostVarMap, "ansible_host", &th.AnsibleHost)
				th.Name = hostName
				th.Vars = hostVarMap.(map[string]interface{})
				// Set host-specific network info
				for netName := range role.Networks {
					th.Networks[netName] = &TripleoHostNetwork{}
					mapStringGetter(hostVarMap, netName+"_ip", &th.Networks[netName].IP)
					mapStringGetter(hostVarMap, netName+"_hostname", &th.Networks[netName].Hostname)
				}
				role.Hosts = append(role.Hosts, th)
			}
		}
		cdl.Roles = append(cdl.Roles, role)
	}
	return nil
}

func (cdl *ConfigDownload) GetNetwork(name string) *OSPNetwork {
	for _, network := range cdl.Networks {
		if network.Name == name {
			return network
		}
	}
	return nil
}

func (cdl *ConfigDownload) ProcessPasswords(mappingYaml string, globalVars map[string]interface{}) error {
	mapping := make(map[string]PasswordMapping)

	err := yaml.Unmarshal([]byte(mappingYaml), &mapping)
	if err != nil {
		return nil
	}

	ocp, err := utils.YamlToMap(utils.GetFullPath(utils.OVERCLOUD_PASSWORDS))
	if err != nil {
		return err
	}

	for k, m := range mapping {
		addr := strings.Split(m.Address, ".")
		var z interface{}
		if m.Role != "" {
			z = globalVars
		} else if m.File != "" {
			z = ocp
		}
		for _, a := range addr {
			if v, ok := z.(map[string]interface{})[a]; ok {
				z = v
			} else {
				fmt.Printf("Error: %s not found in %+v\n", a, z)
				os.Exit(1)
			}
		}
		cdl.Passwords[k] = z.(string)
	}

	return nil
}

func (cdl *ConfigDownload) ProcesNetworks() error {

	networkData, err := utils.YamlToList(utils.GetFullPath(utils.TRIPLEO_OVERCLOUD_NETWORK_DATA))
	if err != nil {
		return err
	}

	for _, netDef := range networkData {
		net, ok := cdl.Networks[netDef["name_lower"].(string)]
		if !ok {
			fmt.Printf("Network %s not found in global_vars.yaml!\n", netDef["name_lower"].(string))
			os.Exit(1)
		}
		if mtu, ok := netDef["mtu"]; ok {
			if mtu != net.Mtu {
				fmt.Printf("MTU mismatch for %s: %d != %d\n", net.Name, mtu, net.Mtu)
				os.Exit(1)
			}
		}
		if vip, ok := netDef["vip"].(string); ok {
			if vip != net.Vip {
				fmt.Printf("VIP mismatch for %s: %s != %s\n", net.Name, vip, net.Vip)
				os.Exit(1)
			}
		}
		subnets, ok := netDef["subnets"].(map[string]interface{})
		if !ok {
			fmt.Printf("No subnets found for %s\n", net.Name)
			os.Exit(1)
		}
		for subnetName, sn := range subnets {
			var ospSn OSPNetworkSubnet
			ospSn.Name = subnetName
			if cidr, ok := sn.(map[string]interface{})["ip_subnet"].(string); ok {
				ospSn.Subnet = cidr
			} else if cidr, ok := sn.(map[string]interface{})["ipv6_subnet"].(string); ok {
				ospSn.Subnet = cidr
			} else {
				fmt.Printf("No subnet definition found for %+v\n", sn)
				os.Exit(1)
			}
			if vlanId, ok := sn.(map[string]interface{})["vlan"]; ok {
				ospSn.VlanId = int(vlanId.(float64))
			}
			if ap, ok := sn.(map[string]interface{})["allocation_pools"].([]interface{}); ok {
				for _, pool := range ap {
					ospSn.AllocationPools = append(ospSn.AllocationPools, OSPNetworkSubnetPool{
						Start: pool.(map[string]interface{})["start"].(string),
						End:   pool.(map[string]interface{})["end"].(string),
					})
				}
			} else if ap, ok := sn.(map[string]interface{})["ipv6_allocation_pools"].([]interface{}); ok {
				for _, pool := range ap {
					ospSn.AllocationPools = append(ospSn.AllocationPools, OSPNetworkSubnetPool{
						Start: pool.(map[string]interface{})["start"].(string),
						End:   pool.(map[string]interface{})["end"].(string),
					})
				}
			} else {
				fmt.Printf("No allocation pools found for %+v\n", sn)
				os.Exit(1)
			}
			net.Subnets = append(net.Subnets, ospSn)
		}
	}

	return nil
}

func (cdl *ConfigDownload) ProcessDeployStepsOne() error {
	steps, err := utils.YamlToList(utils.GetFullPath(utils.DEPLOY_STEPS_ONE))
	if err != nil {
		return err
	}
	importRole := steps[1]

	if vars, ok := importRole["vars"]; ok {
		if tcipc, ok := vars.(map[string]interface{})["tripleo_container_image_prepare_content"]; ok {
			if pd, ok := tcipc.(map[string]interface{})["parameter_defaults"]; ok {
				regexCount := regexp.MustCompile(`([-A-Za-z-0-9.]+)Count`)
				for k, v := range pd.(map[string]interface{}) {
					for _, match := range regexCount.FindAllStringSubmatch(k, -1) {
						fmt.Printf("%s: %d\n", match[1], v)
					}
				}
			} else {
				fmt.Printf("No parameter_defaults found in parameter_defaults\n")
			}
		} else {
			fmt.Printf("No tripleo_container_image_prepare_content found\n")
		}
	}

	return nil
}
