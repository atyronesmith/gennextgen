package generate

import (
	"embed"
	"fmt"

	"github.com/atyronesmith/gennextgen/pkg/utils"
)

// TODO
// The keystone-services-secret has more entries than there are values
// in the overcloud-passwords.yaml.  Looking at the Ansible code, it looks like most often,
// the database password is the service password. (service_config_settings.yaml)

func GenOpenStackSecret(templateDir embed.FS, outDir string, passwords map[string]interface{}, groupVars map[string]interface{}) error {

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

	secret, err := utils.ProcessTemplate(templateDir, "keystone-secret.tmpl", "keystone", utils.GetFuncMap(), tPlate)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret.Bytes(), outDir, "openstack-service-secret.yaml")

	return err
}
