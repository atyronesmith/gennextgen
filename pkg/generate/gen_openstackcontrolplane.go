package generate

import "github.com/atyronesmith/gennextgen/pkg/utils"

func GenOpenStackControlPlane(root string, outDir string) error {

	exportVars, err := utils.YamlToMap(root + "overcloud-export.yaml")
	if err != nil {
		return err
	}

	vals := exportVars["parameter_defaults"].(map[string]interface{})
	services := vals["AllNodesExtraMapData"].(map[string]interface{})

	templateData := struct {
		Services            map[string]interface{}
		GroupVars           map[string]interface{}
		OspSecretName       string
		StorageClass        string
		DnsIps              string
		RabbitMQIps         string
		RabbitMQCell1Ips    string
		DnsServers          []string
		DnsReplicas         int
		GaleraReplicas      int
		GaleraCell1Replicas int
		MemcachedReplicas   int
	}{
		Services:            services,
		StorageClass:        "local-storage",
		OspSecretName:       "osp-secret",
		DnsIps:              "192.168.122.80",
		RabbitMQIps:         "172.17.0.85",
		RabbitMQCell1Ips:    "172.17.0.86",
		DnsServers:          []string{"192.168.122.1"},
		DnsReplicas:         1,
		GaleraReplicas:      1,
		GaleraCell1Replicas: 1,
		MemcachedReplicas:   1,
	}

	secret, err := utils.ProcessTemplate("openstack-control-plane-spec.tmpl", "OpenstackControlPlane", utils.GetFuncMap(), templateData)
	if err != nil {
		return err
	}
	err = utils.WriteByteData(secret, outDir, "openstack-control-plane.yaml")

	return err
}
