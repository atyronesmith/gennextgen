package generate

import (
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func GenSecrets(root string, outDir string) error {

	parameterDefaults, err := utils.YamlToMap(root + "overcloud-passwords.yaml")
	if err != nil {
		return err
	}
	passwords := parameterDefaults["parameter_defaults"].(map[string]interface{})

	groupVars, err := utils.YamlToMap(root + "config-download/overcloud/group_vars/Controller")
	if err != nil {
		return err
	}
	serviceConfigs := groupVars["service_configs"].(map[string]interface{})

	err = genKeystoneSecret(outDir, passwords)
	if err != nil {
		return err
	}

	err = genOpenStackSecret(outDir, passwords, serviceConfigs)
	if err != nil {
		return err
	}
	return nil
}

func genKeystoneSecret(outDir string, passwords map[string]interface{}) error {

	secret, err := utils.ProcessTemplate("keystone-secret.tmpl", "keystone", utils.GetFuncMap(), passwords)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret, outDir, "keystone-secret.yaml")

	return err
}

func genOpenStackSecret(outDir string, passwords map[string]interface{}, groupVars map[string]interface{}) error {
	templateData := struct {
		Passwords map[string]interface{}
		GroupVars map[string]interface{}
		Namespace string
		Name      string
	}{
		Passwords: passwords,
		GroupVars: groupVars,
		Namespace: "openstack",
		Name:      "osp-secret",
	}

	secret, err := utils.ProcessTemplate("openstack-service-secret.tmpl", "openstack", utils.GetFuncMap(), templateData)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret, outDir, "openstack-service-secret.yaml")

	return err
}
