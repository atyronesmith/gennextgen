package generate

import (
	"fmt"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

// TODO
// The keystone-services-secret has more entries than there are values
// in the overcloud-passwords.yaml.  Looking at the Ansible code, it looks like most often,
// the database password is the service password. (service_config_settings.yaml)

func GenOpenStackSecret(outDir string, passwords map[string]interface{}, groupVars map[string]interface{}) error {

	type tStruct struct {
		Passwords map[string]interface{}
		GroupVars map[string]interface{}
	}

	tPlate := tStruct{
		Passwords: passwords,
		GroupVars: groupVars,
	}

	for k := range passwords {
		fmt.Printf("XXXX: %s\n", k)
	}

	secret, err := utils.ProcessTemplate("keystone-secret.tmpl", "keystone", utils.GetFuncMap(), tPlate)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret, outDir, "openstack-service-secret.yaml")

	return err
}

func GetOpenStackPasswords(root string, outDir string) error {
	parameterDefaults, err := utils.YamlToMap(root + "overcloud-passwords.yaml")
	if err != nil {
		return err
	}
	passwords := parameterDefaults["parameter_defaults"].(map[string]interface{})

	groupVars, err := utils.YamlToMap(root + "config-download/overcloud/group_vars/Controller")
	if err != nil {
		return err
	}

	err = GenOpenStackSecret(outDir, passwords, groupVars)
	if err != nil {
		return err
	}

	return nil
}
