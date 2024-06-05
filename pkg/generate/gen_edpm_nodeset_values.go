package generate

import (
	"fmt"

	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
	"github.com/openstack-k8s-operators/dataplane-operator/api/v1beta1"
	infranetworkv1 "github.com/openstack-k8s-operators/infra-operator/apis/network/v1beta1"
)

type EDPMNodesetConfigMap struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string                 `yaml:"name"`
		Annotations map[string]interface{} `yaml:"annotations"`
	} `yaml:"metadata"`
	Data map[string]interface{} `yaml:"data"`
}

func GenEdpmNodesetValues(outDir string, cdl *types.ConfigDownload) error {

	for _, role := range cdl.Roles {
		cfMap, err := genNodeset(role)
		if err != nil {
			return err
		}
		yaml, err := utils.StructToYamlK8s(cfMap)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s-nodeset-values.yaml", role.Name)
		err = utils.WriteByteData(yaml, outDir, fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

// https://github.com/openstack-k8s-operators/dataplane-operator/blob/main/docs/assemblies/ref_example-OpenStackDataPlaneNodeSet-CR-for-preprovisioned-nodes.adoc

func genNodeset(tRole *types.TripleoRole) (EDPMNodesetConfigMap, error) {

	var VANodeDef = EDPMNodesetConfigMap{}

	VANodeDef.APIVersion = "v1"
	VANodeDef.Kind = "ConfigMap"
	VANodeDef.Metadata.Name = fmt.Sprintf("%s-%s", "edpm-nodeset-values", tRole.Name)
	VANodeDef.Metadata.Annotations = map[string]interface{}{
		"config.kubernetes.io/local-config": "true",
	}
	// Use a map to store the values for the nodeset
	// Should this be a struct?
	d := make(map[string]interface{})
	VANodeDef.Data = d

	d["root_password"] = "password"
	d["preProvisioned"] = "false"

	bareMetalSetTemplate := make(map[string]interface{})
	d["baremetalSetTemplate"] = bareMetalSetTemplate

	bareMetalSetTemplate["ctlplaneInterface"] = "CHANGEME - Interface on the provisioned nodes to use for ctlplane network"
	bareMetalSetTemplate["provisioningNetwork"] = "CHANGEME - Network to use for provisioning"
	bareMetalSetTemplate["cloudUserName"] = ""

	nodeset := make(map[string]interface{})
	d["nodeset"] = nodeset

	ansible := make(map[string]interface{})
	nodeset["ansible"] = ansible

	nodeset["networks"] = []infranetworkv1.IPSetNetwork{}

	for _, net := range tRole.Networks {
		rNet := infranetworkv1.IPSetNetwork{}
		for _, sn := range net.Subnets {
			rNet.Name = infranetworkv1.NetNameStr(net.Name)
			rNet.SubnetName = infranetworkv1.NetNameStr(sn.Name)
		}
		nodeset["networks"] = append(nodeset["networks"].([]infranetworkv1.IPSetNetwork), rNet)
	}

	hosts := make(map[string]v1beta1.NodeSection)
	nodeset["nodes"] = hosts

	for _, host := range tRole.Hosts {
		rHost := v1beta1.NodeSection{}

		rHost.HostName = host.Name
		rHost.Ansible.AnsibleHost = host.AnsibleHost

		for _, thn := range host.Networks {
			rNet := infranetworkv1.IPSetNetwork{}
			rNet.Name = infranetworkv1.NetNameStr(thn.Name)
			rNet.SubnetName = infranetworkv1.NetNameStr(thn.BaseSubnet.Name)
			ipAddr := thn.IP.String()
			rNet.FixedIP = &ipAddr

			for _, rn := range host.DefaultRouteNetwork {
				if rn == thn.Name {
					rNet.DefaultRoute = new(bool)
					*rNet.DefaultRoute = true
				}
			}

			rHost.Networks = append(rHost.Networks, rNet)
		}
		hosts[rHost.HostName] = rHost
	}

	return VANodeDef, nil
}
