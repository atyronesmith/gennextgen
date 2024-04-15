package utils

import (
	"fmt"
)

type RhosoConfig struct {
	WorkerNodes    WorkerNodes `yaml:"worker_nodes"`
	ControllerRole string      `yaml:"controller_role,omitempty"`
}

type WorkerNodes struct {
	Names     []string `yaml:"names"`
	Interface string   `yaml:"interface"`
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
