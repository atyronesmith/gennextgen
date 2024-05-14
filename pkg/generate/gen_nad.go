package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"

	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	whereabouts "github.com/k8snetworkplumbingwg/whereabouts/pkg/types"

	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/types/nad/netconf"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func GenNetAttachDef(network *types.OSPNetwork) ([]byte, error) {

	cniCfg := netconf.NewNetConfMacvlan()
	cniCfg.Name = network.Name

	cfg := utils.GetConfig()
	if network.VlanId != 0 {
		cniCfg.Master = fmt.Sprintf("%s.%d", cfg.Interface, network.VlanId)
	} else {
		cniCfg.Master = cfg.Interface
	}
	pf, _ := netip.ParsePrefix("192.168.111.0/24")

	start, err := utils.IpOffset(pf, 30)
	if err != nil {
		fmt.Printf("Error getting IP offset: %v", err)
	}
	end, err := utils.IpOffset(pf, 70)
	if err != nil {
		fmt.Printf("Error getting IP offset: %v", err)
	}

	waCfg := whereabouts.IPAMConfig{}
	waCfg.Type = "whereabouts"
	waCfg.RangeStart = net.ParseIP(start.String())
	waCfg.RangeEnd = net.ParseIP(end.String())
	waCfg.Range = pf.String()

	cniCfg.IPAM = waCfg

	c, err := json.MarshalIndent(cniCfg, "", "  ")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func GenNadFile(outDir string, nads []networkv1.NetworkAttachmentDefinition) error {
	var yamlOut bytes.Buffer

	for _, nad := range nads {
		yaml, err := utils.StructToYamlK8s(nad)
		if err != nil {
			return err
		}
		yamlOut.WriteString("---\n")
		yamlOut.Write(yaml)
	}

	err := utils.WriteByteData(yamlOut.Bytes(), outDir, "openstack-nad.yaml")

	return err
}
