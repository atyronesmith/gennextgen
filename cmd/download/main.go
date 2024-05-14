package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	t "github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
	"gopkg.in/yaml.v3"
)

func main() {
	isHelp := flag.Bool("help", false, "Print usage information.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] edpm_repo\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       edpm_repo  -- Path to directory containing the edpm-ansible repo.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *isHelp {
		flag.Usage()

		os.Exit(0)
	}

	if len(os.Args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	edpmRepo := flag.Arg(0)

	fileList, err := utils.SearchFileRegex(edpmRepo, "*.yml")
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

	// Print argDefs in alphabetical order
	for _, k := range keys {
		v := argDefs[k]
		b, err := yaml.Marshal(v)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Printf("---\n%s\n%s\n", k, string(b))
	}

}
