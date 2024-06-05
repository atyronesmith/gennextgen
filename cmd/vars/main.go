package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	t "github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
	"github.com/wasilibs/go-re2"
	"gopkg.in/yaml.v3"
)

type SearchArg struct {
	VarName   string
	MappedVar string
	Regex     *re2.Regexp
	ArgDef    *t.ArgDef
	Filename  string
	FoundVal  string
}

type Varmap map[string]*SearchArg

func main() {
	isHelp := flag.Bool("help", false, "Print usage information.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] edpm_repo path/config_settings.yaml\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       edpm_repo  -- Path to directory containing the edpm-ansible repo.\n")
		fmt.Fprintf(CommandLine.Output(), "       config_download  -- Path to the config_download directory.\n")
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

	// Sort argDefs alphabetically
	// fileEntires, err := os.ReadDir(filepath.Join(configDownloadDir, "overcloud"))
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	//				sa.Filename = file
	ExtractEPDMVars(edpmRepo)

}

func ExtractEPDMVars(edpmRepo string) {
	// Get a list of all Ansible arg spec files
	fileList, err := utils.SearchFileRegex(edpmRepo, "argument_specs.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	argDefs := t.ArgDefs{}
	// Parse all the arg spec files
	for _, file := range fileList {
		if err := argDefs.ParseArgDefFile(file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	configDownloadDir := flag.Arg(1)

	fileList, err = utils.SearchFileRegex(filepath.Join(configDownloadDir, "overcloud", "config-download", "overcloud"), "*.yaml$")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	varMap := make(map[string]*SearchArg)

	var varList string
	count := 0
	for _, arg := range argDefs {
		varName := strings.Replace(arg.Name, "edpm_", "", 1)
		if count > 0 {
			varList += "|"
		}
		if len(arg.Name) == 0 {
			fmt.Printf("No name found for arg: %+v\n", arg)
			os.Exit(1)
		}
		count++
		varList += varName
		varMap[varName] = &SearchArg{
			VarName: arg.Name,
			ArgDef:  &arg,
		}
		varMap["tripleo_"+varName] = varMap[varName]
	}
	searchVar := fmt.Sprintf("\\s*((tripleo_)*(%s)):\\s*(.*)", varList)
	edpmRegex, err := re2.Compile(searchVar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d edpm definitions...\n", len(argDefs))

	repRegex, err := re2.Compile("config-download/overcloud/.*/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range fileList {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		matches := edpmRegex.FindAllSubmatch(content, -1)
		for _, match := range matches {
			matchStr := string(match[1])
			if sa, ok := varMap[matchStr]; ok {
				foundVal := string(match[4])
				if strings.Contains(foundVal, "{{") {
					continue
				}
				sa.MappedVar = matchStr
				sa.Filename = repRegex.ReplaceAllString(file, "config-download/overcloud/@role/")

				sa.FoundVal = string(match[4])
			}
		}
	}

	foundMap := make(map[string]*SearchArg)
	mappedCount := 0
	for _, vm := range varMap {
		if vm.MappedVar != "" {
			mappedCount++
		}
		foundMap[vm.VarName] = vm
	}
	fmt.Printf("Mapped %d of %d edpm definitions...\n", mappedCount, len(argDefs))

	b, err := yaml.Marshal(foundMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = utils.WriteByteData(b, "/tmp", "varmap.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (sa SearchArg) MarshalYAML() (interface{}, error) {
	node := yaml.Node{}

	if sa.ArgDef == nil {
		return nil, fmt.Errorf("ArgDef is nil")
	}

	err := node.Encode(sa.MappedVar)
	if err != nil {
		return nil, err
	}
	node.FootComment = fmt.Sprintf("Type: %s\nDescription: %s\nRequired: %v\nDefault: %v\nFile: %s\nExample: %s\n",
		sa.ArgDef.Type, strings.TrimSuffix(sa.ArgDef.Description, "\n"), sa.ArgDef.Required,
		sa.ArgDef.Default, sa.Filename, sa.FoundVal)

	return node, nil
}
