package types

import (
	"path/filepath"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

type DeployStepsTasksStep0 []struct {
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

func ProcessGetDeployStepsTasksStep0(cdl *ConfigDownload) error {
	for roleIndex, role := range cdl.Roles {
		configPath := filepath.Join("config-download", "overcloud", role.Name, "deploy_steps_tasks_step0.yaml")

		dsts := DeployStepsTasksStep0{}

		err := utils.YamlToStruct(utils.GetFullPath(configPath), &dsts)
		if err != nil {
			return err
		}

		for _, stepValue := range dsts {
			if stepValue.IncludeRole != nil {
				for k, v := range stepValue.Vars {
					cs := TripleoRoleConfigSetting{
						Service: stepValue.IncludeRole.Name,
						Path:    k,
						Value:   v,
					}
					cdl.Roles[roleIndex].ConfigSettings[k] = append(cdl.Roles[roleIndex].ConfigSettings[k], cs)
				}
			}
		}
	}

	return nil
}
