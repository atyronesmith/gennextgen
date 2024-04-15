package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atyronesmith/gennextgen/pkg/generate"
	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

//go:embed configs/mapping.yaml
var mappingYaml string

func main() {
	isVerbose := flag.Bool("verbose", false, "Print extra runtime information.")
	isHelp := flag.Bool("help", false, "Print usage information.")
	outDir := flag.String("outdir", "/tmp/osp", "Where to place output.")
	configFile := flag.String("config", "./config.yaml", "Path to the config file.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] stack_dir\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       stack_dir  -- Path to directory containing overcloud-deploy stack.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *isVerbose {
		fmt.Println("Verbose...")
	}

	if *isHelp {
		flag.Usage()

		os.Exit(0)
	}

	if dirname := flag.Arg(0); dirname == "" {
		flag.Usage()
		os.Exit(1)
	} else {
		utils.SetRootDir(dirname)

		err := utils.ReadRhosoConfig(*configFile)
		if err != nil {
			os.Exit(1)
		}

		configDownload := types.NewConfigDownload()

		err = configDownload.Process(mappingYaml)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if err := generate.GenSecrets(mappingYaml, *outDir, configDownload); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if err := generate.GenOpenStackControlPlane(*outDir); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		if err := generate.GenNNCP(*outDir, configDownload); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		generate.GenerateNad(*outDir, configDownload)
		//generate.GenNetConfig(dirname)
	}

}
