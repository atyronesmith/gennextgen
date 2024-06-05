package generate

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

type FnMap struct {
	Path     string `json:"path"`
	Function string `json:"function"`
	Type     string `json:"type"`
	Target   string `json:"target"`
}

var varMap map[string]FnMap

func MapVars(varName string, value interface{}) error {

	if fnMap, ok := varMap[varName]; ok {
		switch fnMap.Function {
		case "NeutronPhysnetNUMANodesMapping":
			err := NeutronPhysnetNUMANodesMapping(value, nil)
			if err != nil {
				return err
			}
		case "NeutronTunnelNUMANodes":
			err := NeutronTunnelNUMANodes(value, nil)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("function not found: %s", fnMap.Function)
		}
	}
	return nil
}

// Passthrough function that returns the input string as is.
func Passthrough(value interface{}, config *ini.File) error {

	return nil
}

// Map of phynet name as key and NUMA nodes as value.
// For example: NeutronPhysnetNUMANodesMapping: {'foo': [0, 1], 'bar': [1]}
// where foo and bar are physnet names and corresponding values are list of
// associated numa_nodes
//
// [neutron]
// physnets = foo, bar
// [neutron_physnet_foo]
// numa_nodes = 0
// [neutron_physnet_bar]
// numa_nodes = 0,1
func NeutronPhysnetNUMANodesMapping(s interface{}, config *ini.File) error {

	// TODO: Unmarshal JSON into a map
	var m map[string][]int
	err := json.Unmarshal([]byte(s.([]byte)), &m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !config.HasSection("neutron") {
		_, err := config.NewSection("neutron")
		if err != nil {
			return fmt.Errorf("failed to create section: %w", err)
		}
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	config.Section("neutron").Key("physnets").SetValue(strings.Join(keys, ", "))

	for physnetName, numaNodes := range m {
		sectionName := fmt.Sprintf("neutron_physnet_%s", physnetName)
		if !config.HasSection(sectionName) {
			_, err := config.NewSection(sectionName)
			if err != nil {
				return fmt.Errorf("failed to create section: %w", err)
			}
		}
		config.Section(sectionName).Key("numa_nodes").SetValue(strings.Join(strings.Fields(fmt.Sprint(numaNodes)), ", "))
	}

	return nil
}

func NeutronTunnelNUMANodes(s interface{}, config *ini.File) error {
	var m []int
	err := json.Unmarshal([]byte(s.([]byte)), &m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !config.HasSection("neutron_tunnels") {
		_, err := config.NewSection("neutron_tunnels")
		if err != nil {
			return fmt.Errorf("failed to create section: %w", err)
		}
	}

	config.Section("neutron_tunnels").Key("numa_nodes").SetValue(strings.Join(strings.Fields(fmt.Sprint(m)), ", "))

	return nil
}

func (f *FnMap) NovaPCIPassthrough(value interface{}, config *ini.File) error {
	var spec map[string][]map[string]string
	err := json.Unmarshal([]byte(value.(string)), &spec)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}
