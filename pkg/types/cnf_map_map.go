package types

import (
	"embed"
	"fmt"

	"sigs.k8s.io/yaml"
)

type ConfigMapDefs []ConfigMapEntry

type ConfigMapEntry struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
	OOO    string `yaml:"ooo"`
}

func GetConfigMapping(configs embed.FS) error {
	passwordMap := "configs/cfg-map-map.yaml"
	mappingYaml, err := configs.ReadFile(passwordMap)
	if err != nil {
		return fmt.Errorf("unable to read template file: %s, %v", passwordMap, err)
	}

	var mapping ConfigMapDefs
	err = yaml.Unmarshal([]byte(mappingYaml), &mapping)
	if err != nil {
		return nil
	}

	return nil
}
