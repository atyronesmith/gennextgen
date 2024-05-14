package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atyronesmith/gennextgen/pkg/generate"
	"github.com/atyronesmith/gennextgen/pkg/types"
	"github.com/atyronesmith/gennextgen/pkg/utils"
)

//go:embed configs/password-var-map.yaml
var passwordMap string

//go:embed configs/service-var-map.yaml
var serviceMap string

//go:embed templates/*
var templates embed.FS

func main() {
	isVerbose := flag.Bool("verbose", false, "Print extra runtime information.")
	isHelp := flag.Bool("help", false, "Print usage information.")
	outDir := flag.String("outdir", "/tmp/osp", "Where to place output.")

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "Usage: %s [options] config.yaml stack_dir\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(CommandLine.Output(), "       config.yaml -- Path to the tool config file.\n")
		fmt.Fprintf(CommandLine.Output(), "       stack_dir  -- Path to directory containing overcloud-deploy stack.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *isVerbose {
		fmt.Println("Verbose...")
	}

	if *isHelp || flag.NArg() < 2 {
		flag.Usage()

		os.Exit(0)
	}

	dirName := flag.Arg(1)
	configFile := flag.Arg(0)

	utils.SetRootDir(dirName)

	err := utils.ReadRhosoConfig(configFile)
	if err != nil {
		os.Exit(1)
	}

	configDownload := types.NewConfigDownload()

	err = configDownload.Process(passwordMap, serviceMap)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := generate.GenerateNetworkValues(*outDir, configDownload); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := generate.GenSecrets(passwordMap, templates, *outDir, configDownload); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := generate.GenOpenStackControlPlane(templates, *outDir); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := generate.GenNNCP(*outDir, configDownload); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	err = generate.GenGraph(templates, *outDir, configDownload) // Generate the graph
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

}
