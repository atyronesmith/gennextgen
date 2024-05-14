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
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
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
		findConfigMaps(flag.Arg(1))
	case "process-kustomize":
		processKustomize(flag.Arg(1))
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

func findConfigMaps(baseDir string) {
	fileList, err := utils.SearchFileRegex(baseDir, "*.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range fileList {
		y, err := utils.YamlToMap(file)
		if err == nil {
			if val, ok := y["kind"]; ok {
				if val == "ConfigMap" {
					fmt.Printf("configmap: %s\n", file)
					findcustomServiceConfig(y)
				}
			}
		}
	}

}

func findcustomServiceConfig(configMap map[string]interface{}) {
	if data, ok := configMap["data"]; ok {
		mapTraverse(data, "data")
	}
}

func mapTraverse(m interface{}, addr string) {
	switch v := m.(type) {
	case map[string]interface{}:
		for k, v := range v {
			switch v := v.(type) {
			case map[string]interface{}:
				mapTraverse(v, fmt.Sprintf("%s.%s", addr, k))
			case string:
				if k == "customServiceConfig" {
					err := parseCustomServiceConfig(v)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				} else {
					fmt.Printf("%s.%s.[%s]\n", addr, k, v)
				}
			}
		}
	case []interface{}:
		for i, v := range v {
			mapTraverse(v, fmt.Sprintf("%s[%d].", addr, i))
		}
	}
}

func parseCustomServiceConfig(cfg string) error {
	ini, err := ini.Load([]byte(cfg))
	if err != nil {
		return err
	}
	for _, section := range ini.Sections() {
		fmt.Printf("\t[%s]\n", section.Name())
		for _, key := range section.Keys() {
			fmt.Printf("\t%s=%s\n", key.Name(), key.Value())
		}
	}
	return nil
}
