package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	t "github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
	v1 "k8s.io/api/core/v1"

	v1beta1 "github.com/openstack-k8s-operators/dataplane-operator/api/v1beta1"
	infranetworkv1 "github.com/openstack-k8s-operators/infra-operator/apis/network/v1beta1"
	baremetalv1 "github.com/openstack-k8s-operators/openstack-baremetal-operator/api/v1beta1"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

func main() {
	isHelp := flag.Bool("help", false, "Print usage information.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] architecture_repo\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       architecture_repo  -- Path to directory containing the architecture_repo repo.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *isHelp {
		flag.Usage()

		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	cmd := flag.Arg(0)
	switch cmd {
	case "find-configmaps":
		err := findConfigMaps("/tmp", "archvars.yaml", flag.Arg(1))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "process-kustomize":
		processKustomize(flag.Arg(1))
	case "gen-nodeset":
		genOpenStackDataPlaneNodeSet("/tmp/osp")
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func processKustomize(archRepo string) {

	absArchRepo, err := filepath.Abs(archRepo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resources := make(map[string]interface{})

	traverseKustomizeFiles(absArchRepo, resources)
}

func traverseKustomizeFiles(baseDir string, resources map[string]interface{}) {
	kFile := t.NewKustomizeFile()
	kFile.Path = path.Join(baseDir, "kustomization.yaml")
	fmt.Printf("Processing file: %s\n", kFile.Path)

	k := t.Kustomize{}

	err := utils.YamlToStruct(kFile.Path, &k)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for index, p := range k.Components {
		k.Components[index] = path.Join(path.Dir(kFile.Path), p)
	}
	for index, p := range k.Resources {
		k.Resources[index] = path.Join(path.Dir(kFile.Path), p)
		m, err := utils.YamlToMap(k.Resources[index])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		kind, ok := m["kind"]
		if !ok {
			fmt.Printf("Resource %s does not have a kind\n", k.Resources[index])
			os.Exit(1)
		}
		name, ok := m["metadata"].(map[string]interface{})["name"]
		if !ok {
			fmt.Printf("Resource %s does not have a name\n", k.Resources[index])
			os.Exit(1)
		}
		resKey := ""
		switch kind {
		case "L2Advertisement":
			fallthrough
		case "ConfigMap":
			fallthrough
		case "IPAddressPool":
			resKey = fmt.Sprintf("%s:%s", kind, name.(string))
		default:
			fmt.Printf("%s/%s %s\n", kind, name, k.Resources[index])
			os.Exit(1)
		}
		fmt.Printf("Adding resource %s\n", resKey)
		kFile.Resources[resKey] = m
		resources[resKey] = m
	}
	for _, p := range k.Components {
		traverseKustomizeFiles(p, resources)
	}
	for _, r := range k.Replacements {
		resKey := fmt.Sprintf("%s:%s", r.Source.Kind, r.Source.Name)
		res, ok := resources[resKey]
		if !ok {
			fmt.Printf("Resource %s not found\n", resKey)
			os.Exit(1)
		}
		fieldPath := strings.Split(r.Source.FieldPath, ".")
		resPtr := res
		for _, fp := range fieldPath {
			resPtr, ok = resPtr.(map[string]interface{})[fp]
			if !ok {
				fmt.Printf("Field %s/%s not found in resource %s\n", fp, r.Source.FieldPath, resKey)
				os.Exit(1)
			}
		}
		fmt.Printf("%s: <%s>\n", r.Source.FieldPath, resPtr)
	}

	out, err := yaml.Marshal(kFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	kFile.Data = k
	fmt.Println(string(out))
}

func findConfigMaps(outDir string, fileName string, baseDir string) error {
	fileList, err := utils.SearchFileRegex(baseDir, "*.yaml")
	if err != nil {
		return err
	}

	_, err = os.Stat(outDir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("unable to mkdir: %s", outDir)
		}
	}
	file := filepath.Join(outDir, fileName)
	dFile, err := os.Create(file)
	if err != nil {
		fmt.Printf("unable to create/open file: %s", file)
	}
	defer dFile.Close()

	fmt.Printf("Writing %s...\n", file)

	for _, file := range fileList {
		y, err := utils.YamlToMap(file)
		if err == nil {
			if val, ok := y["kind"]; ok {
				if val == "ConfigMap" {
					fmt.Printf("configmap: %s\n", file)
					findcustomServiceConfig(y, dFile)
				}
			}
		}
	}
	return nil
}

func findcustomServiceConfig(configMap map[string]interface{}, file *os.File) {
	if data, ok := configMap["data"]; ok {
		mapTraverse(data, "data", file)
	}
}

func mapTraverse(m interface{}, addr string, file *os.File) {
	switch v := m.(type) {
	case map[string]interface{}:
		for k, v := range v {
			switch v := v.(type) {
			case map[string]interface{}:
				mapTraverse(v, fmt.Sprintf("%s.%s", addr, k), file)
			case string:
				if k == "customServiceConfig" || k == "conf" {
					err := parseCustomServiceConfig(v, addr, file)
					if err != nil {
						fmt.Printf("Cannot parse INI file...%s\n", err)
					}
				} else {
					_, err := file.WriteString(fmt.Sprintf("%s.%s.[%s]\n", addr, k, v))
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}
		}
	case []interface{}:
		for i, v := range v {
			mapTraverse(v, fmt.Sprintf("%s[%d].", addr, i), file)
		}
	default:
		fmt.Printf("unhandled map type %s=%v\n", addr, v)
	}
}

func parseCustomServiceConfig(cfg string, addr string, file *os.File) error {
	ini, err := ini.Load([]byte(cfg))
	if err != nil {
		return err
	}
	for _, section := range ini.Sections() {
		//		fmt.Printf("\t[%s]\n", section.Name())
		for _, key := range section.Keys() {
			_, err := file.WriteString(fmt.Sprintf("%s.%s.%s=%s\n", addr, section.Name(), key.Name(), key.Value()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genOpenStackDataPlaneNodeSet(outDir string) {
	var vnn = v1beta1.OpenStackDataPlaneNodeSetSpec{}

	vnn.BaremetalSetTemplate = baremetalv1.OpenStackBaremetalSetSpec{
		CtlplaneInterface:     "eth0",
		CloudUserName:         "cloud-admin",
		ProvisioningInterface: "eth1",
		BmhLabelSelector: map[string]string{
			"app": "openstack",
		},
		PasswordSecret: &v1.SecretReference{
			Name:      "baremetalset-passowrd-secret",
			Namespace: "openstack",
		},
	}

	vnn.NodeTemplate.Ansible = v1beta1.AnsibleOpts{
		AnsibleUser: "cloud-admin",
		AnsiblePort: 22,
	}

	dRoute := true
	vnn.NodeTemplate.Networks = []infranetworkv1.IPSetNetwork{
		{
			Name:         "ctlplane",
			DefaultRoute: &dRoute,
			SubnetName:   "ctlplane",
		},
		{
			Name:       "internalapi",
			SubnetName: "subnet1",
		},
		{
			Name:       "storage",
			SubnetName: "subnet1",
		},
		{
			Name:       "tenant",
			SubnetName: "subnet1",
		},
	}
	vnn.Nodes = make(map[string]v1beta1.NodeSection)
	vnn.Nodes["edpm-compute-0"] = v1beta1.NodeSection{
		HostName: "compute-0",
		Networks: []infranetworkv1.IPSetNetwork{
			{
				Name:       "ctlplane",
				SubnetName: "ctlplane",
			},
			{
				Name:       "internalapi",
				SubnetName: "subnet1",
			},
		},
	}

	// serviceMap := map[string]string{
	// 	"ceph_client":        "ceph_client",
	// 	"ovn_metadata":       "neutron-metadata",
	// 	"neutron_ovn_dpdk":   "configure-ovs-dpdk",
	// 	"ca_certs":           "install-certs",
	// 	"ovn_controller":     "ovn",
	// 	"nova_libvirt":       "libvirt",
	// 	"nova_libvirt_guest": "", // No mapping
	// 	"nova_compute":       "nova",
	// 	"iscsid":             "nova",
	// }

	// edpm_nova contains
	// - NovaComputeImage
	// - EdpmIscsiImage

	// install_os contains
	// - epdm_podman
	// - edpm_sshd
	// - dataplane_chrony / timesync
	// - edpm_logrotate

	vnn.Services = []string{
		"bootstrap",
		"download-cache",
		"reboot-os",
		"configure-ovs-dpdk",
		"configure-network",
		"validate-network",
		"install-os",
		"configure-os",
		"ssh-known-hosts",
		"run-os",
		"install-certs",
		"ovn",
		"neutron-ovn",
		"neutron-metadata",
		"neutron-sriov",
		"libvirt",
		"nova-custom-ovsdpdksriov",
		"telemetry",
	}

	yaml, err := utils.StructToYamlK8s(vnn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = utils.WriteByteData(yaml, outDir, "openstack-dataplane-nodeset-spec.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
