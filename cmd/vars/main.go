package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	t "github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

func main() {
	isHelp := flag.Bool("help", false, "Print usage information.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] edpm_repo path/config_settings.yaml\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       edpm_repo  -- Path to directory containing the edpm-ansible repo.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *isHelp {
		flag.Usage()

		os.Exit(0)
	}

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	edpmRepo := flag.Arg(0)

	fileList, err := utils.SearchFileRegex(edpmRepo, "argument_specs.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	argDefs := t.ArgDefs{}

	for _, file := range fileList {
		if err := argDefs.ParseArgDefFile(file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Sort argDefs alphabetically
	keys := make([]string, 0, len(argDefs))
	for k := range argDefs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	cs, err := getConfigSettings(flag.Arg(1))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// for k, v := range cs {
	// 	fmt.Printf("%s: %s\n", k, v.Value)
	// }
	// Print argDefs in alphabetical order
	for _, k := range keys {
		v := argDefs[k]
		varName := strings.Replace(v.Name, "edpm_", "", 1)
		if val, ok := cs[varName]; ok {
			fmt.Printf("%s: %+v\n", varName, val)
		}
	}
}

func getConfigSettings(configPath string) (map[string]t.TripleoRoleConfigSetting, error) {
	cfgSet, err := utils.YamlToMap(configPath)
	if err != nil {
		return nil, err
	}

	csm := make(map[string]t.TripleoRoleConfigSetting)

	for k, v := range cfgSet {
		path := strings.Split(k, "::")
		settingKey := path[len(path)-1]
		cs := t.TripleoRoleConfigSetting{
			Service: path[0],
			Path:    k,
			Value:   v,
		}
		if len(path) > 1 {
			cs.Section = path[1]
		}
		csm[settingKey] = cs
	}

	return csm, nil
}
