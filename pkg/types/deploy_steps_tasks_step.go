package types

import (
	"fmt"
	"path/filepath"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

type DeployStepsTasksStep []struct {
	IncludeRole *IncludeRole           `yaml:"include_role,omitempty"`
	Name        string                 `yaml:"name"`
	Vars        map[string]interface{} `yaml:"vars,omitempty"`
	When        string                 `yaml:"when"`
	Block       []Block                `yaml:"block,omitempty"`
}
type IncludeRole struct {
	Name      string `yaml:"name"`
	TasksFrom string `yaml:"tasks_from"`
}
type Block struct {
	Name     string                 `yaml:"name"`
	Register string                 `yaml:"register,omitempty"`
	Shell    string                 `yaml:"shell"`
	Vars     map[string]interface{} `yaml:"vars,omitempty"`
	When     string                 `yaml:"when,omitempty"`
}

func ProcessGetDeploySteps(cdl *ConfigDownload) error {
	err := processGetDeployStepsTasksStep(cdl, "deploy_steps_tasks_step0.yaml")
	if err != nil {
		return err
	}

	err = processGetDeployStepsTasksStep(cdl, "pre_deploy_step_tasks.yaml")
	if err != nil {
		return err
	}

	return nil
}

func processGetDeployStepsTasksStep(cdl *ConfigDownload, fileName string) error {
	for roleIndex, role := range cdl.Roles {
		configPath := filepath.Join("config-download", "overcloud", role.Name, fileName)

		dsts := DeployStepsTasksStep{}

		err := utils.YamlToStruct(utils.GetFullPath(configPath), &dsts)
		if err != nil {
			return err
		}

		for _, stepValue := range dsts {
			if stepValue.IncludeRole != nil {
				for k, v := range stepValue.Vars {
					if cdl.Roles[roleIndex].ConfigSettings[k] == nil {
						cdl.Roles[roleIndex].ConfigSettings[k] = v
					} else {
						return fmt.Errorf("processGetDeployStepsTasksStep: %s already exists for role %s", k, role.Name)
					}
				}
			}
		}
	}

	return nil
}
