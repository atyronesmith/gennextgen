package generate

import (
	"fmt"

	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func GenerateNetworkValues(outDir string, cdl *types.ConfigDownload) error {

	configMap := make(map[string]interface{})

	configMap["APIVersion"] = "v1"
	configMap["kind"] = "ConfigMap"
	configMap["metadata"] = map[string]interface{}{
		"name": "network-values",
		"annotations": map[string]interface{}{
			"config.kubernetes.io/local-config": "true",
		},
	}
	data := make(map[string]interface{})
	configMap["data"] = data
	networks := make(map[string]interface{})
	nodes := make(map[string]interface{})
	data["networks"] = networks
	data["nodes"] = nodes

	cfg := utils.GetConfig()

	for _, network := range cdl.Networks {
		currentNet := make(map[string]interface{})
		currentNet["dnsDomain"] = network.DnsDomain
		if network.VlanId != 0 {
			currentNet["vlan"] = network.VlanId
			currentNet["base_iface"] = cfg.Interface
			currentNet["iface"] = fmt.Sprintf("%s.%d", cfg.Interface, network.VlanId)
		}
		currentNet["vip"] = network.Vip
		currentNet["mtu"] = network.Mtu
		currentNet["prefix-length"] = network.PrefixLen

		currentNet["subnets"] = network.Subnets
		// Generate NetworkAttachmentDefinition
		nad, err := GenNetAttachDef(network)
		if err != nil {
			return err
		}
		currentNet["network-attachment-definition"] = string(nad)

		networks[network.NameLower] = currentNet
	}

	controllerHosts := make([]*types.TripleoHost, 0)
	for _, role := range cdl.Roles {
		for _, tag := range role.RoleTags {
			if tag == "controller" {
				controllerHosts = append(controllerHosts, role.Hosts...)
			}
		}
	}
	// Need to map the controller role to the worker node
	// And then allocate acutal hosts from each controller role
	// take the network information from each host
	for index, workerNode := range cfg.WorkerNodes {
		currentNode := make(map[string]interface{})
		nodes[fmt.Sprintf("node_%d", index)] = currentNode
		currentNode["name"] = workerNode.Name
		host := controllerHosts[index]
		for _, network := range host.Networks {
			currentNode[fmt.Sprintf("%s_ip", network.Name)] = network.IP
		}
	}

	yaml, err := utils.StructToYamlK8s(configMap)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(yaml, outDir, "network-values.yaml")
	if err != nil {
		return err
	}

	return nil
}
