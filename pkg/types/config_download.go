package types

import (
	"embed"
	"fmt"
	"net/netip"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/atyronesmith/gennextgen/pkg/utils"
	"sigs.k8s.io/yaml"
)

type ConfigDownload struct {
	Networks        map[string]*OSPNetwork
	EnabledServices map[string]*OSPService
	Hosts           map[string]*TripleoHost
	Roles           map[string]*TripleoRole
	PrimaryRoleName string
	Passwords       map[string]string
	ConfigSettings  ConfigSettings
}

type ConfigSettings struct {
	MechanismDrivers  []string
	NetworkVlanRanges []string
	GlobalPhysnetMtu  int
}

type PasswordMapping struct {
	File    string `yaml:"file"`
	Address string `yaml:"address"`
	Role    string `yaml:"role"`
}

type ServiceMapping struct {
	Setting       string `yaml:"setting"`
	ConfigVarName string `yaml:"config_var_name"`
	Service       string `yaml:"service"`
	Section       string `yaml:"section"`
}

type OSPServicePasswords struct {
	Password   string
	DbPassword string
}

type OSPNetwork struct {
	Name        string             `json:"name"`
	NameLower   string             `json:"name_lower"`
	DnsDomain   string             `json:"domain,omitempty"`
	CloudName   string             `json:"cloud_name,omitempty"`
	Mtu         int                `json:"mtu"`
	IpV6        bool               `json:"ipv6"`
	Vip         netip.Addr         `json:"vip"`
	Subnets     []OSPNetworkSubnet `json:"subnets,omitempty"`
	PrefixLen   int                `json:"prefix_len"`
	VlanId      int                `json:"vlan_id"`
	HostRoutes  []TripleoRoutes    `json:"host_routes"`
	IsCtrlPlane bool               `json:"is_ctrl_plane"`
}

type OSPNetworkSubnet struct {
	Name                string                 `json:"name"`
	IpSubnet            netip.Prefix           `json:"ip_subnet"`
	Ipv6Subnet          netip.Prefix           `json:"ipv6_subnet,omitempty"`
	GatewayIp           netip.Addr             `json:"gateway_ip,omitempty"`
	GatewayIpV6         netip.Addr             `json:"gateway_ip_v6,omitempty"`
	AllocationPools     []OSPNetworkSubnetPool `json:"allocationRanges,omitempty"`
	Ipv6AllocationPools []OSPNetworkSubnetPool `json:"ipv6_allocationRanges,omitempty"`
	Routes              []TripleoRoutes        `json:"routes,omitempty"`
	RoutesIpv6          []TripleoRoutes        `json:"routes_ipv6,omitempty"`
	Vlan                int                    `json:"vlan,omitempty"`
}

type OSPNetworkSubnetPool struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type TripleoRoutes struct {
	Default     bool   `yaml:"default"`
	Destination string `yaml:"destination"`
	NextHop     string `yaml:"nexthop"`
}

type TripleoHost struct {
	Name                string
	CanonicalName       string
	AnsibleHost         string
	DefaultRouteNetwork []string
	Networks            map[string]*TripleoHostNetwork
	TripleoRole         *TripleoRole
	Vars                map[string]interface{}
}

type TripleoHostNetwork struct {
	Name        string
	IP          netip.Addr
	Hostname    string
	SubnetName  string
	NetworkName string
	BaseSubnet  *OSPNetworkSubnet
	RoleNetwork *TripleoRoleNetwork
}

//   - tags: (list) list of tags used by other parts of the deployment process to
//     find the role for a specific type of functionality. Currently a role
//     with both 'primary' and 'controller' is used as the primary role for the
//     deployment process. If no roles have 'primary' and 'controller', the
//     first role in this file is used as the primary role.
//     The third tag that can be defined here is external_bridge, which is used
//     to define which node must have a bridge created in a multiple-nic network
//     config.
//   - ovsdpdk: (boolean) whether the role is using OVS-DPDK or not.
//   - storage: (boolean) whether the role is using storage or not.
//   - ceph: (boolean) whether the role is using Ceph or not.
//   - compute: (boolean) whether the role is a compute role or not.

type TripleoRole struct {
	Name string

	ConfigSettings map[string]interface{}

	RoleTags     []string
	Hosts        []*TripleoHost
	Networks     map[string]*TripleoRoleNetwork
	GrowvolsArgs map[string]interface{}

	Vars map[string]interface{}
}

type TripleoRoleNetwork struct {
	Name       string          `json:"name"`
	NameLower  string          `json:"name_lower"`
	DnsDomain  string          `json:"domain,omitempty"`
	CloudName  string          `json:"cloud_name,omitempty"`
	Mtu        int             `json:"mtu"`
	IpV6       bool            `json:"ipv6"`
	Vip        netip.Addr      `json:"vip"`
	GatewayIp  netip.Addr      `json:"gateway_ip"`
	Subnets    []string        `json:"subnets,omitempty"`
	PrefixLen  int             `json:"prefix_len"`
	VlanId     int             `json:"vlan_id"`
	HostRoutes []TripleoRoutes `json:"host_routes"`
	IsRoleNet  bool            `json:"is_role_net"`
}

type RoleType string

const (
	RoleTypeController    RoleType = "controller"
	RoleTypeCompute       RoleType = "compute"
	RoleTypeCephStorage   RoleType = "ceph-storage"
	RoleTypeBlockStorage  RoleType = "block-storage"
	RoleTypeObjectStorage RoleType = "object-storage"
	RoleTypeNetworker     RoleType = "networker"
	Unknown               RoleType = "unknown"
)

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
		Roles:           make(map[string]*TripleoRole, 0),
	}
}

func NewTripleoRole() *TripleoRole {
	return &TripleoRole{
		Networks:       make(map[string]*TripleoRoleNetwork),
		ConfigSettings: make(map[string]interface{}),
	}
}

func (cdl *ConfigDownload) Process(outDir string, configs embed.FS, serviceMap string) error {
	//	var nncp types.NNCP
	cfg := utils.GetConfig()

	gvs, err := GetGlobalVars()
	if err != nil {
		return err
	}
	err = cdl.ProcessGlobalVars(gvs)
	if err != nil {
		return err
	}

	tond, err := GetTripleoOvercloudNetworkData()
	if err != nil {
		return err
	}
	err = cdl.ProcessTripleoOvercloudNetworkData(tond)
	if err != nil {
		return err
	}

	// b, err := yaml.Marshal(toe)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%s\n", string(b))

	err = cdl.ProcessTripleoAnsibleInventory()
	if err != nil {
		return err
	}

	err = cdl.ProcessTripleoOvercloudBaremetalDeployment()
	if err != nil {
		return err
	}

	tor, err := GetTripleoOvercloudRolesData("")
	if err != nil {
		return err
	}

	err = cdl.ProcessTripleoOvercloudRolesData(tor)
	if err != nil {
		return err
	}

	primaryRole := cdl.PrimaryRoleName
	if cfg.ControllerRole != "" {
		primaryRole = cfg.ControllerRole
	}

	groupVars, err := cdl.ProcessGroupVars(primaryRole)
	if err != nil {
		return err
	}

	err = cdl.ProcessPasswords(configs, groupVars)
	if err != nil {
		return err
	}
	err = ProcessConfigSettings(cdl)
	if err != nil {
		return err
	}
	err = cdl.ProcessTripleoOvercloudEnvironment()
	if err != nil {
		return err
	}

	err = ProcessGetDeploySteps(cdl)
	if err != nil {
		return err
	}

	// err = cdl.SaveConfigSettings(outDir)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Must be run first
func (cdl *ConfigDownload) ProcessGlobalVars(globalVars *GlobalVars) error {
	for netName, cidrs := range globalVars.NetCIDRMap {
		_, ok := cdl.Networks[netName]
		if !ok {
			cdl.Networks[netName] = &OSPNetwork{
				NameLower: netName,
			}
			for _, cidr := range cidrs {
				if addr, err := netip.ParsePrefix(cidr); err == nil {
					cdl.Networks[netName].Subnets = append(cdl.Networks[netName].Subnets, OSPNetworkSubnet{
						IpSubnet: addr,
					})
				} else {
					return err
				}
			}
		}
	}

	for cloudName, cn := range globalVars.CloudNames {
		netName := strings.Replace(cloudName, "cloud_name_", "", -1)
		if _, ok := cdl.Networks[netName]; ok {
			cdl.Networks[netName].CloudName = cn
		}
	}

	for netName, vip := range globalVars.NetworkVirtualIPS {
		if _, ok := cdl.Networks[netName]; ok {
			if addr, err := netip.ParseAddr(vip.IpAddress); err == nil {
				cdl.Networks[netName].Vip = addr
			} else {
				return err
			}
		}
	}

	cdl.PrimaryRoleName = globalVars.PrimaryRoleName

	return nil
}

func (cdl *ConfigDownload) ProcessGroupVars(primaryRole string) (map[string]interface{}, error) {
	groupVars, err := utils.YamlToMap(utils.GetFullPath("config-download/overcloud/group_vars/" + primaryRole))
	if err != nil {
		fmt.Printf("Error while reading group_vars: %v\n", err)
		return nil, err
	}

	for index := range cdl.Networks {
		if m, ok := groupVars[index+"_mtu"]; ok {
			cdl.Networks[index].Mtu = m.(int)
		} else {
			return nil, fmt.Errorf("Missing mtu for %s, in group_vars.yaml\n", index)
		}
		if vi, ok := groupVars[index+"_vlan_id"]; ok {
			cdl.Networks[index].VlanId = vi.(int)
		}
	}

	return groupVars, nil
}

func mapStringGetter(v interface{}, key string, valPtr interface{}) error {
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
		case *netip.Prefix:
			if value != nil {
				pf, err := netip.ParsePrefix(value.(string))
				if err != nil {
					fmt.Printf("Error parsing prefix: %v\n", err)
					os.Exit(1)
				}
				*p = pf
			}
		case *netip.Addr:
			if value != nil {
				pf, err := netip.ParseAddr(value.(string))
				if err != nil {
					fmt.Printf("Error parsing address: %v\n", err)
					os.Exit(1)
				}
				*p = pf
			}
		default:
			fmt.Printf("Unknown type: %T\n", p)
			os.Exit(1)
		}
	} else {
		return fmt.Errorf("Missing %s\n", key)
	}

	return nil
}

func (cdl *ConfigDownload) ProcessTripleoAnsibleInventory() error {
	data, err := utils.YamlToMap(utils.GetFullPath(utils.TRIPLEO_ANSIBLE_INVENTORY_YAML))
	if err != nil {
		return err
	}

	for roleNameKey, roleData := range data {
		if roleNameKey == "Undercloud" {
			continue
		}
		if _, ok := roleData.(map[string]interface{})["children"]; ok {
			continue
		}

		role := NewTripleoRole()

		// Get the role networks
		// role networks are common to all hosts in the role
		if v, ok := roleData.(map[string]interface{})["vars"]; ok {

			if trn, ok := v.(map[string]interface{})["tripleo_role_name"]; ok {
				role.Name = trn.(string)
			} else {
				return fmt.Errorf("Missing tripleo_role_name in %s\n", roleNameKey)
			}
			role.Networks = make(map[string]*TripleoRoleNetwork)
			// tripleo_role_networks is a list of network names
			// present in the this role
			if trn, ok := v.(map[string]interface{})["tripleo_role_networks"]; ok {
				networks := trn.([]interface{})

				// Get the list of names of networks in this role
				for _, nn := range networks {
					netName := nn.(string)
					osNet := &TripleoRoleNetwork{
						NameLower: netName,
					}
					role.Networks[netName] = osNet

					// There are two possible names for the subnet cidr,
					// try both
					// A PrefixLen is required, so if neither is found, return an error
					if err := mapStringGetter(v, netName+"_subnet_cidr", &osNet.PrefixLen); err != nil {
						if err = mapStringGetter(v, netName+"_cidr", &osNet.PrefixLen); err != nil {
							return err
						}
					}
					// The gateway_ip may be null
					if err := mapStringGetter(v, netName+"_gateway_ip", &osNet.GatewayIp); err != nil {
						return err
					}
					// There is no Vip
					// The MTU must exist
					if err := mapStringGetter(v, netName+"_mtu", &osNet.Mtu); err != nil {
						return err
					}
					// The vlan_id must exist, but may be 0
					if err := mapStringGetter(v, netName+"_vlan_id", &osNet.VlanId); err != nil {
						return err
					}
					// HostRoutes are optional.  The control plane network usually has them.
					if hostRoutes, ok := v.(map[string]interface{})[netName+"_host_routes"]; ok {
						var hrs []TripleoRoutes

						utils.MarshalArray(hostRoutes, &hrs)
						osNet.HostRoutes = append(osNet.HostRoutes, hrs...)
					}
				}
			} else {
				return fmt.Errorf("Missing networks_lower in %s\n", roleNameKey)
			}

			for roleVarName, roleVarValue := range v.(map[string]interface{}) {
				if strings.HasPrefix(roleVarName, "tripleo_") {
					err = SetConfigSetting(role, roleVarName, roleVarValue)
					if err != nil {
						return err
					}
				}
			}
		} else {
			return fmt.Errorf("%s: Missing vars in %s\n", "TripleoAnsibleInventory", roleNameKey)
		}

		if hosts, ok := roleData.(map[string]interface{})["hosts"]; ok {
			for hostName, hostVarMap := range hosts.(map[string]interface{}) {
				// Create the host object
				th := TripleoHost{}

				// Set a pointer to the role to speed searches
				th.TripleoRole = role

				th.Networks = make(map[string]*TripleoHostNetwork)

				if err := mapStringGetter(hostVarMap, "canonical_hostname", &th.CanonicalName); err != nil {
					return err
				}

				if err := mapStringGetter(hostVarMap, "ansible_host", &th.AnsibleHost); err != nil {
					return err
				}

				if drn, ok := hostVarMap.(map[string]interface{})["default_route_network"]; ok {
					var drns []string
					utils.MarshalArray(drn, &drns)
				}

				// TODO - this is a hack to work with GraphViz
				// Rmove the hyphens from the name in the template
				th.Name = strings.ReplaceAll(hostName, "-", "_")
				// Save all the vars for later
				th.Vars = hostVarMap.(map[string]interface{})
				// Set host-specific network info
				for netName := range role.Networks {
					thn := &TripleoHostNetwork{}

					thn.RoleNetwork = role.Networks[netName]

					th.Networks[netName] = thn

					thn.Name = netName

					if err := mapStringGetter(hostVarMap, netName+"_ip", &thn.IP); err != nil {
						return err
					}
					if err := mapStringGetter(hostVarMap, netName+"_hostname", &thn.Hostname); err != nil {
						return err
					}
					if net, ok := cdl.Networks[netName]; ok {
						found := false
						// Set a pointer to the subnet the network is attached to
						for snIndex, sn := range net.Subnets {
							if sn.IpSubnet.Contains(thn.IP) {
								thn.BaseSubnet = &net.Subnets[snIndex]
								found = true
							}
							if sn.Ipv6Subnet.Contains(thn.IP) {
								thn.BaseSubnet = &net.Subnets[snIndex]
								found = true
							}
						}
						if !found {
							return fmt.Errorf("ProcessTripleoAnsibleInventory: IP %s not found in subnet for network %s\n", thn.IP, netName)
						}
					} else {
						return fmt.Errorf("ProcessTripleoAnsibleInventory: Network %s not found in global_vars.yaml\n", netName)
					}
				}

				if drn, ok := hostVarMap.(map[string]interface{})["default_route_network"]; ok {
					var drns []string
					utils.MarshalArray(drn, &drns)

					th.DefaultRouteNetwork = append(th.DefaultRouteNetwork, drns...)
				}

				role.Hosts = append(role.Hosts, &th)

				cdl.Hosts[hostName] = &th

			}

			cdl.Roles[role.Name] = role
		}
	}
	return nil
}

func (cdl *ConfigDownload) ProcessPasswords(configs embed.FS, globalVars map[string]interface{}) error {
	mapping := make(map[string]PasswordMapping)

	passwordMap := "configs/password-var-map.yaml"
	mappingYaml, err := configs.ReadFile(passwordMap)
	if err != nil {
		return fmt.Errorf("unable to read template file: %s, %v", passwordMap, err)
	}

	err = yaml.Unmarshal([]byte(mappingYaml), &mapping)
	if err != nil {
		return nil
	}

	ocp, err := utils.YamlToMap(utils.GetFullPath(utils.OVERCLOUD_PASSWORDS))
	if err != nil {
		return err
	}

	for k, m := range mapping {
		fmt.Printf("Processing %s, %+v\n", k, m)
		addr := strings.Split(m.Address, ".")
		fmt.Printf("Address: %v\n", addr)
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
				z = "CHANGE_ME"
			}
		}
		cdl.Passwords[k] = z.(string)
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

func (cdl *ConfigDownload) ProcessTripleoOvercloudRolesData(tord *TripleoOvercloudRolesData) error {

	for _, role := range *tord {
		for trIndex, tripleoRole := range cdl.Roles {
			if tripleoRole.Name == role.Name {
				cdl.Roles[trIndex].RoleTags = append(cdl.Roles[trIndex].RoleTags, role.Tags...)
			}
		}
	}

	return nil
}

func (cdl *ConfigDownload) ProcessTripleoOvercloudNetworkData(networkData *TripleoOvercloudNetworkData) error {

	for _, netDef := range *networkData {
		ospNetwork := &OSPNetwork{}
		cdl.Networks[netDef.NameLower] = ospNetwork

		ospNetwork.Name = netDef.Name
		ospNetwork.NameLower = netDef.NameLower
		ospNetwork.DnsDomain = netDef.DNSDomain
		ospNetwork.IpV6 = netDef.Ipv6
		ospNetwork.Mtu = netDef.MTU
		//		ospNetwork.Vip = netDef.Vip

		_, ok := cdl.Networks[netDef.NameLower]
		if !ok {
			fmt.Printf("Network %s not found in global_vars.yaml!\n", netDef.Name)
			os.Exit(1)
		}
		for subnetName, sn := range netDef.Subnets {
			var ospSubnet OSPNetworkSubnet

			ospSubnet.Name = subnetName
			ospSubnet.Vlan = int(sn.VLAN)

			for _, pool := range sn.AllocationPools {
				ospSubnet.AllocationPools = append(ospSubnet.AllocationPools, OSPNetworkSubnetPool(pool))
			}

			if sn.GatewayIpV6 != "" {
				gatewayIpV6, err := netip.ParseAddr(sn.GatewayIpV6)
				if err != nil {
					return err
				}
				ospSubnet.GatewayIpV6 = gatewayIpV6
			}

			for _, route := range sn.Routes {
				ospSubnet.Routes = append(ospSubnet.Routes, TripleoRoutes(route))
			}

			if sn.IpSubnet != "" {
				subnetPrefix, err := netip.ParsePrefix(sn.IpSubnet)
				if err != nil {
					return err
				}
				ospSubnet.IpSubnet = subnetPrefix
			}

			if sn.Ipv6Subnet != "" {
				subnetPrefix, err := netip.ParsePrefix(sn.Ipv6Subnet)
				if err != nil {
					return err
				}
				ospSubnet.Ipv6Subnet = subnetPrefix
			}

			for _, route := range sn.RoutesIpv6 {
				ospSubnet.RoutesIpv6 = append(ospSubnet.RoutesIpv6, TripleoRoutes(route))
			}

			ospNetwork.Subnets = append(ospNetwork.Subnets, ospSubnet)
		}
	}

	return nil
}

// Extract Control Plane information from the TripleoOvercloudEnvironment
// The Control Plane network is defined as part of the undercloud
func (cdl *ConfigDownload) ProcessTripleoOvercloudEnvironment() error {

	environment, err := GetTripleoOvercloudEnvironment("")
	if err != nil {
		return err
	}

	// Get the name of the control plane network
	net := OSPNetwork{}

	net.IsCtrlPlane = true

	net.NameLower = environment.ParmaterDefaults.CtlplaneNetworkAttributes.Network.Name
	net.DnsDomain = environment.ParmaterDefaults.CtlplaneNetworkAttributes.Network.DNSDomain
	net.Mtu = int(environment.ParmaterDefaults.CtlplaneNetworkAttributes.Network.MTU)

	for _, subnet := range environment.ParmaterDefaults.CtlplaneNetworkAttributes.Subnets {
		ospSubnet := OSPNetworkSubnet{}
		ospSubnet.Name = subnet.Name

		prefix, err := netip.ParsePrefix(subnet.CIDR)
		if err != nil {
			return err
		}
		ospSubnet.IpSubnet = prefix
		//TODO Host Routes
		net.Subnets = append(net.Subnets, ospSubnet)
	}

	if len(environment.ParmaterDefaults.ControlPlaneVipData.FixedIPS) > 0 {
		vipAddr, err := netip.ParseAddr(environment.ParmaterDefaults.ControlPlaneVipData.FixedIPS[0].IPAddress)
		if err != nil {
			return err
		}
		// TODO -- needs to be an array?
		net.Vip = vipAddr
	}

	cdl.Networks[net.NameLower] = &net

	env, err := utils.YamlToMap(utils.GetFullPath(utils.TRIPLEO_OVERCLOUD_ENVIRONMENT))
	if err != nil {
		return err
	}

	defaults := env["parameter_defaults"].(map[string]interface{})

	mapTypes := []string{"ExtraConfig", "ExtraGroupVars", "Parameters"}

	for _, role := range cdl.Roles {
		for _, mapType := range mapTypes {
			key := fmt.Sprintf("%s%s", role.Name, mapType)
			if roleParams, ok := defaults[key]; ok {
				fmt.Printf("Processing %s\n", key)
				for k, v := range roleParams.(map[string]interface{}) {
					varName := strings.ReplaceAll(k, "::", ".")
					err = SetConfigSetting(role, varName, v)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (cdl *ConfigDownload) ProcessConfigSettings(cfgSet map[string]interface{}, serviceMap string) error {
	mapping := make(map[string]ServiceMapping)

	err := yaml.Unmarshal([]byte(serviceMap), &mapping)
	if err != nil {
		return nil
	}

	drivers, ok := cfgSet["neutron::plugins::ml2::mechanism_drivers"].([]string)
	if ok {
		cdl.ConfigSettings.MechanismDrivers = append(cdl.ConfigSettings.MechanismDrivers, drivers...)
	}

	vlanRanges, ok := cfgSet["neutron::plugins::ml2::network_vlan_ranges"].([]string)
	if ok {
		cdl.ConfigSettings.NetworkVlanRanges = append(cdl.ConfigSettings.NetworkVlanRanges, vlanRanges...)
	}

	physnetMtu, ok := cfgSet["neutron::global_physnet_mtu"].(int)
	if ok {
		cdl.ConfigSettings.GlobalPhysnetMtu = physnetMtu
	}

	return nil
}

func (cdl *ConfigDownload) ProcessTripleoOvercloudBaremetalDeployment() error {
	bmd, err := GetTripleoOvercloudBaremetalDeployment("")
	if err != nil {
		return err
	}

	for _, bm := range *bmd {
		if role, ok := cdl.Roles[bm.Name]; !ok {
			return fmt.Errorf("Role %s not found in roles\n", bm.Name)
		} else {
			for _, bmAnsible := range bm.AnsiblePlaybooks {
				if strings.Contains(bmAnsible.Playbook, "cli-overcloud-node-growvols.yaml") {
					role.GrowvolsArgs = bmAnsible.ExtraVars["role_growvols_args"].(map[string]interface{})
				}
			}
			for _, inst := range bm.Instances {
				if host, ok := cdl.Hosts[inst.Hostname]; !ok {
					return fmt.Errorf("Host %s not found in hosts\n", inst.Name)
				} else {
					for _, net := range inst.Networks {
						if _, ok := role.Networks[net.Network]; !ok {
							return fmt.Errorf("Network %s not found in role %s\n", net.Network, role.Name)
						} else {
							host.Networks[net.Network].SubnetName = net.Subnet
							host.Networks[net.Network].NetworkName = net.Network
							// TODO check that the IP address matches what we already have
							role.Networks[net.Network].Subnets = append(role.Networks[net.Network].Subnets, net.Subnet)
						}
					}
				}
			}
		}
	}

	return nil
}

func (cdl *ConfigDownload) SaveConfigSettings(outDir string) error {
	// Sort the keys of ConfigSettings alphabetically
	var keys []string
	for _, role := range cdl.Roles {
		for key := range role.ConfigSettings {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// Print ConfigSettings in alphabetical order
	for _, role := range cdl.Roles {
		for _, key := range keys {
			fmt.Printf("%s: %+v\n", key, role.ConfigSettings[key])
		}
	}

	return nil
}

var edpmMap map[string]string

func InitEDPMVarMap(configs embed.FS) error {
	if edpmMap == nil {
		edpmMap = make(map[string]string)

		varMap := "configs/edpm-var-map.yaml"
		mappingYaml, err := configs.ReadFile(varMap)
		if err != nil {
			return fmt.Errorf("unable to read template file: %s, %v", varMap, err)
		}

		err = yaml.Unmarshal([]byte(mappingYaml), &edpmMap)
		if err != nil {
			return nil
		}
	}

	return nil
}

func EDPMVarMap(varName string, tr *TripleoRole) interface{} {

	if mapped, ok := edpmMap[varName]; ok {
		if val, ok := tr.ConfigSettings[mapped]; ok {
			return val
		}
		return "CHANGE_ME"
	}

	fmt.Printf("unable to find mapping for %s", varName)

	return "Error"
}

func GetServiceVars(role *TripleoRole, service string) map[string]interface{} {
	serviceVars := make(map[string]interface{})

	for k, v := range role.ConfigSettings {
		path := strings.Split(k, ".")
		if len(path) > 1 {
			if strings.Contains(path[0], service) {
				serviceVars[k] = v
			}
		}
	}

	return serviceVars
}

func SetConfigSetting(role *TripleoRole, key string, value interface{}) error {
	strVal := fmt.Sprintf("%v", value)
	if len(strVal) == 0 {
		fmt.Printf("\tSetConfigSetting: <%s> value is empty, skipping...\n", key)

		return nil
	}
	if strings.ContainsAny(strVal, "{}&:*#?|<>=!%@\\") {
		fmt.Printf("\tSetConfigSetting: <%s> value contains a special char <%s>, skipping...\n", key, strVal)

		return nil
	}
	if role.ConfigSettings[key] != nil {
		if role.ConfigSettings[key] != value {
			fmt.Printf("\tSetConfigSetting: <%s> already exists for role %s\n", key, role.Name)
			fmt.Printf("\t<%+v> != <%+v>\n", role.ConfigSettings[key], value)
		}
	}

	role.ConfigSettings[key] = value

	return nil
}
