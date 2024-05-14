package utils

import (
	"fmt"
)

type RhosoConfig struct {
	WorkerNodes    []WorkerNode `yaml:"worker_nodes"`
	Interface      string       `yaml:"interface"`
	ControllerRole string       `yaml:"controller_role,omitempty"`
}

type WorkerNode struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

var rhosoConfig RhosoConfig

func ReadRhosoConfig(path string) error {
	err := YamlToStruct(path, &rhosoConfig)
	if err != nil {
		fmt.Printf("Error while reading config file %s: %v", path, err)
		return err
	}
	fmt.Printf("rhosoConfig: %v\n", rhosoConfig)

	return nil
}

func GetConfig() RhosoConfig {
	return rhosoConfig
}
