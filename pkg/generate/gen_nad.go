package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	whereabouts "github.com/k8snetworkplumbingwg/whereabouts/pkg/types"

	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/types/nad/netconf"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func GenerateNad(outDir string, cdl *types.ConfigDownload) {
	cfg := utils.GetConfig()

	nads := []networkv1.NetworkAttachmentDefinition{}

	for _, network := range cdl.Networks {
		nad := networkv1.NetworkAttachmentDefinition{}
		nad.APIVersion = "k8s.cni.cncf.io/v1"
		nad.Kind = "NetworkAttachmentDefinition"

		nad.Name = network.Name
		nad.Namespace = "openstack"

		cniCfg := netconf.NewNetConfMacvlan()
		cniCfg.Name = network.Name
		if network.VlanId != 0 {
			cniCfg.Master = fmt.Sprintf("%s.%d", cfg.WorkerNodes.Interface, network.VlanId)
		} else {
			cniCfg.Master = cfg.WorkerNodes.Interface
		}

		start, err := utils.IpOffset(network.Cidr, 30)
		if err != nil {
			fmt.Printf("Error getting IP offset: %v", err)
		}
		end, err := utils.IpOffset(network.Cidr, 70)
		if err != nil {
			fmt.Printf("Error getting IP offset: %v", err)
		}

		waCfg := whereabouts.IPAMConfig{}
		waCfg.Type = "whereabouts"
		waCfg.RangeStart = net.ParseIP(start.String())
		waCfg.RangeEnd = net.ParseIP(end.String())
		waCfg.Range = network.Cidr.String()

		cniCfg.IPAM = waCfg

		c, err := json.MarshalIndent(cniCfg, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling CNIConfig: %v", err)

			return
		}
		nad.Spec.Config = string(c)

		nads = append(nads, nad)
	}

	err := GenNadFile(outDir, nads)
	if err != nil {
		fmt.Printf("Error generating NAD file: %v", err)
	}
	// Your code for generating NAD goes here
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
