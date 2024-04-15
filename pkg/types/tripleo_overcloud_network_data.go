package types

type OvercloudNetwork struct {
	Name      string                   `json:"name"`
	NameLower string                   `json:"name_lower"`
	Vip       bool                     `json:"vip"`
	Mtu       int                      `json:"mtu"`
	Subnets   []OvercloudNetworkSubnet `json:"subnets"`
}

type OvercloudNetworkSubnet struct {
	IpSubnet        string                                 `json:"ip_subnet"`
	AllocationPools []OvercloudNetworkSubnetAllocationPool `json:"allocation_pools"`
	Vlan            int                                    `json:"vlan"`
}

type OvercloudNetworkSubnetAllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type TripleoOvercloudNetworkData struct {
	Networks []OvercloudNetwork `json:"networks"`
}

func (tond *TripleoOvercloudNetworkData) Process(data []map[string]interface{}) error {
	var taod TripleoOvercloudNetworkData

	for _, net := range data {
		var on OvercloudNetwork
		on.Name = net["name"].(string)
		on.NameLower = net["name_lower"].(string)
		on.Vip = net["vip"].(bool)
		on.Mtu = net["mtu"].(int)
		for _, subnet := range net["subnets"].(map[string]interface{}) {
			var os OvercloudNetworkSubnet
			os.IpSubnet = subnet.(map[string]interface{})["ip_subnet"].(string)
			os.Vlan = subnet.(map[string]interface{})["vlan"].(int)
			for _, pool := range subnet.(map[string]interface{})["allocation_pools"].([]interface{}) {
				var oap OvercloudNetworkSubnetAllocationPool
				oap.Start = pool.(map[string]interface{})["start"].(string)
				oap.End = pool.(map[string]interface{})["end"].(string)
				os.AllocationPools = append(os.AllocationPools, oap)
			}
			on.Subnets = append(on.Subnets, os)
		}
		taod.Networks = append(taod.Networks, on)
	}

	return nil
}

func (tond *TripleoOvercloudNetworkData) GetNetwork(name string) *OvercloudNetwork {
	for _, network := range tond.Networks {
		if network.Name == name {
			return &network
		}
	}
	return nil
}
