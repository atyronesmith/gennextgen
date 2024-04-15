package generate

import (
	"fmt"

	"github.com/atyronesmith/gennextgen/pkg/utils"

	netconfig "github.com/atyronesmith/gennextgen/pkg/types/openstack"
)

func GenNetConfig(root string) {

	ndPath := root + "tripleo-overcloud-network-data.yaml"
	networkData, err := utils.YamlToList(ndPath)
	if err != nil {
		fmt.Printf("Error while reading global vars file %s: %v", ndPath, err)
		return
	}

	nc := netconfig.NewNetConfig()
	nc.ObjectMeta.Namespace = "openstack"

	for _, netDef := range networkData {
		nw := netconfig.NewNetwork()
		nw.Name = netconfig.NetNameStr(netDef["name"].(string))
		nw.DNSDomain = "TODO"

		for key, sn := range netDef["subnets"].(map[string]interface{}) {
			ncsn := netconfig.Subnet{}
			ncsn.Name = netconfig.NetNameStr(key) // Convert key to openstack.NetNameStr type
			ncsn.Cidr = sn.(map[string]interface{})["ip_subnet"].(string)
			val, ok := sn.(map[string]interface{})["gateway_ip"]
			if ok {
				ncsn.Gateway = val.(string)
			}
			for _, pool := range sn.(map[string]interface{})["allocation_pools"].([]interface{}) {
				ncsn.AllocationRanges = append(ncsn.AllocationRanges, netconfig.AllocationRange{
					Start: pool.(map[string]interface{})["start"].(string),
					End:   pool.(map[string]interface{})["end"].(string),
				})
			}
			nw.Subnets = append(nw.Subnets, ncsn)
		}

		nc.Spec.Networks = append(nc.Spec.Networks, *nw)
	}
	ncYaml, err := utils.StructToYaml(nc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", ncYaml)

}
