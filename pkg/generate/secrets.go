package generate

import "github.com/atyronesmith/gennextgen/pkg/utils"

func GenSecrets(root string, outDir string) error {

	passwords, err := getCredentials(root)
	if err != nil {
		return err
	}
	err = genKeystoneSecret(outDir, passwords)
	if err != nil {
		return err
	}

	err = genOpenStackSecret(outDir, passwords)
	if err != nil {
		return err
	}
	return nil
}

func getCredentials(root string) (map[string]interface{}, error) {
	parameterDefaults, err := utils.YamlToMap(root + "overcloud-passwords.yaml")
	if err != nil {
		return nil, err
	}
	return parameterDefaults["parameter_defaults"].(map[string]interface{}), nil

}

func genKeystoneSecret(outDir string, passwords map[string]interface{}) error {

	secret, err := utils.ProcessTemplate("keystone-secret.tmpl", "keystone", utils.GetFuncMap(), passwords)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret, outDir, "keystone-secret.yaml")

	return err
}

func genOpenStackSecret(outDir string, passwords map[string]interface{}) error {

	secret, err := utils.ProcessTemplate("openstack-service-secret.tmpl", "openstack", utils.GetFuncMap(), passwords)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret, outDir, "openstack-service-secret.yaml")

	return err
}
